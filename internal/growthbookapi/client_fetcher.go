package growthbookapi

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type handler[T any] struct {
	client *Client
	path   string
	method string
}

func fetcher[T any](client *Client, method, path string) *handler[T] {
	return &handler[T]{
		client: client,
		path:   path,
		method: method,
	}
}

func (h *handler[T]) One(ctx context.Context, body any, resultKey string) (T, error) {
	var zero T
	tflog.Debug(ctx, "Fetching single resource", map[string]any{
		"method": h.method,
		"path":   h.path,
	})
	resp, err := h.client.do(ctx, h.method, h.path, body)
	if err != nil {
		tflog.Error(ctx, "Error in doRequest for fetchSingle", map[string]any{
			"error":  err.Error(),
			"method": h.method,
			"path":   h.path,
		})
		return zero, err
	}
	defer func() { _ = resp.Body.Close() }()

	if err := checkStatuses(h.method, resp); err != nil {
		tflog.Error(ctx, "Status check failed in fetchSingle", map[string]any{
			"error":  err.Error(),
			"status": resp.StatusCode,
			"method": h.method,
			"path":   h.path,
		})
		return zero, err
	}
	var respMap map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&respMap); err != nil {
		tflog.Error(ctx, "JSON decode failed in fetchSingle", map[string]any{
			"error":  err.Error(),
			"method": h.method,
			"path":   h.path,
		})
		return zero, err
	}
	return decodeResultKey[T](respMap, resultKey)
}

func (h *handler[T]) page(
	ctx context.Context,
	urlStr string,
	body any,
	resultKey string,
) ([]T, bool, int, error) {
	resp, err := h.client.do(ctx, h.method, urlStr, body)
	if err != nil {
		return nil, false, 0, err
	}
	defer func() { _ = resp.Body.Close() }()
	if err := checkStatuses(h.method, resp); err != nil {
		return nil, false, 0, err
	}
	var page map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
		return nil, false, 0, err
	}
	items, err := decodeResultKey[[]T](page, resultKey)
	if err != nil {
		return nil, false, 0, err
	}
	hasMore, _ := page["hasMore"].(bool)
	nextOffset, _ := page["nextOffset"].(float64)
	return items, hasMore, int(nextOffset), nil
}

func (h *handler[T]) All(
	ctx context.Context,
	body any,
	resultKey string,
) ([]T, error) {
	var allItems []T
	offset := 0
	for {
		tflog.Debug(ctx, "Fetching page for paginated API", map[string]any{
			"method": h.method,
			"path":   h.path,
			"offset": offset,
			"limit":  h.client.Limit,
		})
		parsedURL, err := url.Parse(h.path)
		if err != nil {
			return nil, err
		}
		q := parsedURL.Query()
		q.Set("limit", strconv.Itoa(h.client.Limit))
		if offset > 0 {
			q.Set("offset", strconv.Itoa(offset))
		}
		parsedURL.RawQuery = q.Encode()
		urlStr := parsedURL.String()
		items, hasMore, nextOffset, err := h.page(ctx, urlStr, body, resultKey)
		if err != nil {
			tflog.Error(ctx, "Error fetching page", map[string]any{
				"error":  err.Error(),
				"offset": offset,
				"limit":  h.client.Limit,
			})
			return nil, err
		}
		allItems = append(allItems, items...)
		tflog.Debug(ctx, "Fetched page", map[string]any{
			"items_fetched": len(items),
			"total_items":   len(allItems),
			"has_more":      hasMore,
			"next_offset":   nextOffset,
			"limit":         h.client.Limit,
		})
		if !hasMore {
			break
		}
		offset = nextOffset
	}
	return allItems, nil
}

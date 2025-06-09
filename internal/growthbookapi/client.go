// Package growthbookapi provides a client for interacting with the GrowthBook API.
package growthbookapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ErrNotFound is returned when a resource is not found (HTTP 404).
var ErrNotFound = errors.New("growthbookapi: resource not found")

// Client is a GrowthBook API client that can be used to interact with the GrowthBook API.
// It supports making HTTP requests to the API and includes options for customization.
// The Client is initialized with a base URL, API key, and optional HTTP client.
// It provides methods to perform API requests and handles logging of requests and responses.
// The API key is redacted in logs for security purposes.
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// Option is a function that configures a Client.
type Option func(*Client)

// WithHTTPClient sets a custom http.Client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		if httpClient != nil {
			c.HTTPClient = httpClient
		}
	}
}

// NewClient creates a Client with optional configuration options.
func NewClient(baseURL, apiKey string, opts ...Option) *Client {
	client := &Client{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		HTTPClient: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(client)
	}

	return client
}

func redactAPIKey(apiKey string) string {
	if len(apiKey) <= 6 {
		return "***REDACTED***"
	}

	return apiKey[:3] + "***REDACTED***" + apiKey[len(apiKey)-3:]
}

// doRequest executes an HTTP request to the GrowthBook API with the provided context.
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var buf io.Reader
	var bodyLog string
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(b)
		bodyLog = string(b)
	}
	url := c.BaseURL + path
	redactedAPIKey := redactAPIKey(c.APIKey)
	tflog.Debug(ctx,
		"HTTP Request",
		map[string]interface{}{
			"method":        method,
			"url":           url,
			"authorization": "Bearer " + redactedAPIKey,
			"body":          bodyLog,
		},
	)

	req, err := http.NewRequestWithContext(ctx, method, url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		tflog.Debug(ctx,
			"HTTP Response Error",
			map[string]interface{}{
				"error":  err.Error(),
				"method": method,
				"url":    url,
			},
		)
		return nil, err
	}
	defer func() {
		_ = req.Body
	}()
	respBody, _ := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewBuffer(respBody))
	tflog.Debug(ctx,
		"HTTP Response",
		map[string]interface{}{
			"method": method,
			"url":    url,
			"status": resp.StatusCode,
			"body":   strings.TrimSpace(string(respBody)),
		},
	)

	return resp, nil
}

// fetchOne performs an HTTP request and decodes the JSON response into out using generics.
func fetchOne[T any](
	ctx context.Context,
	c *Client,
	method string,
	path string,
	body interface{},
	expectedStatus ...int,
) (T, error) {
	var zero T
	resp, err := c.doRequest(ctx, method, path, body)
	if err != nil {
		return zero, err
	}
	defer func() { _ = resp.Body.Close() }()

	ok := false
	for _, code := range expectedStatus {
		if resp.StatusCode == code {
			ok = true
			break
		}
	}
	if !ok {
		b, _ := io.ReadAll(resp.Body)
		return zero, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(b))
	}
	var out T
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return zero, err
	}
	return out, nil
}

// fetchPage performs a single paginated API request and decodes the response.
func fetchPage[T any](
	ctx context.Context,
	c *Client, method string,
	urlStr string,
	body interface{},
	resultKey string,
	expectedStatus ...int,
) ([]T, bool, int, error) {
	resp, err := c.doRequest(ctx, method, urlStr, body)
	if err != nil {
		return nil, false, 0, err
	}
	defer func() { _ = resp.Body.Close() }()
	ok := false
	for _, code := range expectedStatus {
		if resp.StatusCode == code {
			ok = true
			break
		}
	}
	if !ok {
		b, _ := io.ReadAll(resp.Body)
		return nil, false, 0, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(b))
	}
	var page map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
		return nil, false, 0, err
	}
	itemsRaw, ok := page[resultKey]
	if !ok {
		return nil, false, 0, fmt.Errorf("response missing '%s' key", resultKey)
	}
	itemsBytes, err := json.Marshal(itemsRaw)
	if err != nil {
		return nil, false, 0, err
	}
	var items []T
	if err := json.Unmarshal(itemsBytes, &items); err != nil {
		return nil, false, 0, err
	}
	hasMore, _ := page["hasMore"].(bool)
	nextOffset, _ := page["nextOffset"].(float64)
	return items, hasMore, int(nextOffset), nil
}

// fetchAllPages performs repeated HTTP requests to fetch all pages of a paginated list
// endpoint and decodes the results into a single slice.
func fetchAllPages[T any](
	ctx context.Context,
	c *Client,
	method string,
	path string,
	body interface{},
	resultKey string,
	expectedStatus ...int,
) ([]T, error) {
	var allItems []T
	offset := 0
	for {
		parsedURL, err := url.Parse(path)
		if err != nil {
			return nil, err
		}
		q := parsedURL.Query()
		q.Set("limit", strconv.Itoa(100))
		if offset > 0 {
			q.Set("offset", strconv.Itoa(offset))
		}
		parsedURL.RawQuery = q.Encode()
		urlStr := parsedURL.String()
		items, hasMore, nextOffset, err := fetchPage[T](ctx, c, method, urlStr, body, resultKey, expectedStatus...)
		if err != nil {
			return nil, err
		}
		allItems = append(allItems, items...)
		if !hasMore {
			break
		}
		offset = nextOffset
	}
	return allItems, nil
}

// remove performs a DELETE request and checks for expected status codes.
func (c *Client) remove(ctx context.Context, path string, expectedStatus ...int) error {
	resp, err := c.doRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	for _, code := range expectedStatus {
		if resp.StatusCode == code {
			return nil
		}
	}
	b, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(b))
}

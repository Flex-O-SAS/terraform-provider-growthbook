package growthbookapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// APIError represents an error response from the GrowthBook API.
type APIError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	HTTPStatus string `json:"-"`
}

var (
	//nolint:gochecknoglobals
	createStatuses = []int{http.StatusOK, http.StatusCreated}
	//nolint:gochecknoglobals
	updateStatuses = []int{http.StatusOK}
	//nolint:gochecknoglobals
	readStatuses = []int{http.StatusOK}
	//nolint:gochecknoglobals
	deleteStatuses = []int{http.StatusOK, http.StatusNoContent}
	//nolint:gochecknoglobals
	methodStatuses = map[string][]int{
		"GET":    readStatuses,
		"POST":   createStatuses,
		"PUT":    updateStatuses,
		"DELETE": deleteStatuses,
		"PATCH":  updateStatuses,
	}
)

func (e APIError) Error() string {
	msg := e.Message
	if msg == "" {
		msg = e.HTTPStatus
	}
	return fmt.Sprintf("unexpected status: %d %s", e.StatusCode, msg)
}

func checkStatuses(method string, statusCode int, body []byte) error {
	expected, found := methodStatuses[method]
	if !found {
		return fmt.Errorf("unsupported method %s", method)
	}
	if statusCode == 0 {
		return errors.New("response is nil")
	}

	if statusCode == http.StatusNotFound {
		return ErrNotFound
	}
	for _, code := range expected {
		if statusCode == code {
			return nil
		}
	}

	// For error responses, try to extract the API error message
	apiErr := APIError{
		StatusCode: statusCode,
		HTTPStatus: http.StatusText(statusCode),
	}

	if len(body) > 0 {
		var errorResp map[string]interface{}
		if err := json.Unmarshal(body, &errorResp); err == nil {
			if msg, ok := errorResp["message"].(string); ok {
				apiErr.Message = msg
			}
		}
	}

	return apiErr
}

func decodeResultKey[T any](m map[string]any, key string) (T, error) {
	var zero T
	val, ok := m[key]
	if !ok {
		return zero, fmt.Errorf("response missing '%s' key", key)
	}
	valBytes, err := json.Marshal(val)
	if err != nil {
		return zero, err
	}
	var out T
	if err := json.Unmarshal(valBytes, &out); err != nil {
		return zero, err
	}
	return out, nil
}

// 1. read response bod for logging purposes, replace it with NopCloser buffer.
func (c *Client) do(ctx context.Context, method, path string, body any) (*http.Response, error) {
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
		map[string]any{
			"method":        method,
			"url":           url,
			"authorization": "Bearer " + redactedAPIKey,
			"body":          bodyLog,
		},
	)

	resp, err := c.withRetry(ctx, method, url, buf)
	if err != nil {
		tflog.Debug(ctx,
			"HTTP Response Error",
			map[string]any{
				"error":  err.Error(),
				"method": method,
				"url":    url,
			},
		)
		return nil, err
	}

	// 1.
	respBody, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewBuffer(respBody))

	tflog.Debug(ctx,
		"HTTP Response",
		map[string]any{
			"method": method,
			"url":    url,
			"status": resp.StatusCode,
			"body":   strings.TrimSpace(string(respBody)),
		},
	)

	return resp, nil
}

func (c *Client) retryAfter(
	ctx context.Context,
	resp *http.Response,
	attempt int,
	interval,
	maxInterval time.Duration,
) time.Duration {
	retryAfter := resp.Header.Get("Retry-After")
	var wait time.Duration
	if retryAfter != "" {
		if seconds, parseErr := strconv.Atoi(retryAfter); parseErr == nil {
			wait = time.Duration(seconds) * time.Second
		} else if t, parseErr := http.ParseTime(retryAfter); parseErr == nil {
			wait = time.Until(t)
		}
	}
	if wait <= 0 {
		wait = interval
	}
	tflog.Warn(ctx, "Received 429 Too Many Requests, backing off", map[string]any{
		"attempt":     attempt + 1,
		"retry_after": retryAfter,
		"wait_ms":     wait.Milliseconds(),
	})
	time.Sleep(wait)
	interval = time.Duration(math.Min(float64(maxInterval), float64(interval)*2.0))
	return interval
}

// 1. success or non-retryable error, exit reties and return.
// 2. close response body as it won't be closed by caller.
func (c *Client) withRetry(ctx context.Context, method, url string, buf io.Reader) (*http.Response, error) {
	var resp *http.Response
	var err error

	attempt := 0
	interval := c.Backoff.InitialInterval
	for {
		req, reqErr := http.NewRequestWithContext(ctx, method, url, buf)
		if reqErr != nil {
			return nil, reqErr
		}

		req.Header.Set("Authorization", "Bearer "+c.APIKey)
		req.Header.Set("Content-Type", "application/json")
		resp, err = c.HTTPClient.Do(req)

		if err == nil {
			if resp.StatusCode == http.StatusTooManyRequests {
				interval = c.retryAfter(ctx, resp, attempt, interval, c.Backoff.MaxInterval)
			}
			if resp.StatusCode < 500 {
				// 1.
				break
			}
		}

		if attempt >= c.Backoff.MaxRetries {
			break
		}

		attempt++
		tflog.Warn(ctx, "Transient error, retrying request", map[string]any{
			"attempt":     attempt,
			"interval_ms": interval.Milliseconds(),
			"status": func() int {
				if resp != nil {
					return resp.StatusCode
				}
				return 0
			}(),
			"error": func() string {
				if err != nil {
					return err.Error()
				}
				return ""
			}(),
		})
		time.Sleep(interval)
		interval = time.Duration(math.Min(float64(c.Backoff.MaxInterval), float64(interval)*c.Backoff.Multiplier))
		// 2.
		if err != nil {
			_ = resp.Body.Close()
		}
	}

	return resp, err
}

func (c *Client) delete(ctx context.Context, path string) error {
	resp, err := c.do(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	for _, code := range deleteStatuses {
		if resp.StatusCode == code {
			return nil
		}
	}

	b, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(b))
}

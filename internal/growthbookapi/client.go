package growthbookapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ErrNotFound is returned when a resource is not found (HTTP 404).
var ErrNotFound = errors.New("growthbookapi: resource not found")

type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
	context    context.Context
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
func NewClient(ctx context.Context, baseURL, apiKey string, opts ...Option) *Client {
	client := &Client{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		HTTPClient: http.DefaultClient,
		context:    ctx,
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

// doRequest is a shared helper for making HTTP requests.
func (c *Client) doRequest(method, path string, body interface{}) (*http.Response, error) {
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
	tflog.Debug(c.context,
		"HTTP Request",
		"method", method,
		"url", url,
		"authorization", "Bearer "+redactedAPIKey,
		"body", bodyLog,
	)

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		tflog.Debug(c.context,
			"HTTP Response Error",
			"error", err.Error(),
			"method", method,
			"url", url,
		)
		return nil, err
	}
	defer func() {
		_ = req.Body
	}()
	respBody, _ := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewBuffer(respBody))
	tflog.Debug(c.context,
		"HTTP Response",
		"method", method,
		"url", url,
		"status", resp.StatusCode,
		"body", strings.TrimSpace(string(respBody)),
	)
	return resp, nil
}

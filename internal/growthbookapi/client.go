// Package growthbookapi provides a client for interacting with the GrowthBook API.
package growthbookapi

import (
	"context"
	"errors"
	"net/http"
	"time"
)

var (
	// ErrNotFound is returned when a resource is not found (HTTP 404).
	ErrNotFound = errors.New("growthbookapi: resource not found")
)

// ClientAPI defines the interface for the GrowthBook API client.
type ClientAPI interface {
	// FindProjectByName retrieves a project by its name.
	FindProjectByName(ctx context.Context, name string) (*Project, error)
	// CreateProject creates a new project.
	CreateProject(ctx context.Context, p *Project) (*Project, error)
	// GetProject retrieves a project by its ID.
	GetProject(ctx context.Context, id string) (*Project, error)
	// UpdateProject updates an existing project by its ID.
	UpdateProject(ctx context.Context, id string, p *Project) (*Project, error)
	// DeleteProject deletes a project by its ID.
	DeleteProject(ctx context.Context, id string) error
	// FindEnvironmentByID retrieves an environment by its ID.
	FindEnvironmentByID(ctx context.Context, id string) (*Environment, error)
	// CreateEnvironment creates a new environment.
	CreateEnvironment(ctx context.Context, e *Environment) (*Environment, error)
	// UpdateEnvironment updates an existing environment by its ID.
	UpdateEnvironment(ctx context.Context, id string, e *Environment) (*Environment, error)
	// DeleteEnvironment deletes an environment by its ID.
	DeleteEnvironment(ctx context.Context, id string) error
	// CreateFeature creates a new feature.
	CreateFeature(ctx context.Context, f *Feature) (*Feature, error)
	// GetFeature retrieves a feature by its ID.
	GetFeature(ctx context.Context, id string) (*Feature, error)
	// UpdateFeature updates an existing feature by its ID.
	UpdateFeature(ctx context.Context, id string, f *Feature) (*Feature, error)
	// DeleteFeature deletes a feature by its ID.
	DeleteFeature(ctx context.Context, id string) error
	// FindFeatureByName retrieves a feature by its ID.
	FindFeatureByName(ctx context.Context, id string) (*Feature, error)
	// CreateSDKConnection creates a new SDK connection.
	CreateSDKConnection(ctx context.Context, c *SDKConnection) (*SDKConnection, error)
	// GetSDKConnection retrieves an SDK connection by its ID.
	GetSDKConnection(ctx context.Context, id string) (*SDKConnection, error)
	// UpdateSDKConnection updates an existing SDK connection by its ID.
	UpdateSDKConnection(ctx context.Context, id string, c *SDKConnection) (*SDKConnection, error)
	// DeleteSDKConnection deletes an SDK connection by its ID.
	DeleteSDKConnection(ctx context.Context, id string) error
	// FindSDKConnectionByName retrieves an SDK connection by its name.
	FindSDKConnectionByName(ctx context.Context, name string) (*SDKConnection, error)
}

// BackoffConfig defines the configuration for retrying transient errors.
type BackoffConfig struct {
	MaxRetries      int
	InitialInterval time.Duration
	Multiplier      float64
	MaxInterval     time.Duration
}

// Option is a function that configures a Client.
type Option func(*Client)

// Client is a GrowthBook API client that can be used to interact with the GrowthBook API.
// It supports making HTTP requests to the API and includes options for customization.
// The Client is initialized with a base URL, API key, and optional HTTP client.
// It provides methods to perform API requests and handles logging of requests and responses.
// The API key is redacted in logs for security purposes.
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
	Backoff    BackoffConfig
	Limit      int
}

// NewClient creates a Client with optional configuration options.
func NewClient(baseURL, apiKey string, opts ...Option) ClientAPI {
	client := &Client{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		HTTPClient: http.DefaultClient,
		Backoff: BackoffConfig{
			MaxRetries:      3,
			InitialInterval: 500 * time.Millisecond,
			Multiplier:      2.0,
			MaxInterval:     5 * time.Second,
		},
		Limit: 100,
	}
	for _, opt := range opts {
		opt(client)
	}

	return client
}

// WithHTTPClient sets a custom http.Client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		if httpClient != nil {
			c.HTTPClient = httpClient
		}
	}
}

// WithPageLimit sets the maximum number of items to return per page in paginated API responses.
func WithPageLimit(limit int) Option {
	return func(c *Client) {
		c.Limit = limit
	}
}

// WithBackoff sets a custom backoff configuration for transient error retries.
func WithBackoff(cfg BackoffConfig) Option {
	return func(c *Client) {
		c.Backoff = cfg
	}
}

func redactAPIKey(apiKey string) string {
	if len(apiKey) <= 6 {
		return "***REDACTED***"
	}

	return apiKey[:3] + "***REDACTED***" + apiKey[len(apiKey)-3:]
}

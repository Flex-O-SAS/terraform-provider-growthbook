// Package internal provides the Terraform provider implementation for GrowthBook.
package internal

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-provider-growthbook/internal/growthbookapi"
)

// Provider returns the Terraform provider schema.ResourceProvider for GrowthBook.
// It defines the provider's configuration schema, sets up the API client, and
// manages resources and data sources for GrowthBook projects, features, and
// environments.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GROWTHBOOK_API_KEY", nil),
				Description: "The GrowthBook API key.",
			},
			"api_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GROWTHBOOK_API_URL", "https://api.growthbook.io/api/v1"),
				Description: "The GrowthBook API base URL.",
			},
			"http_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60,
				Description: "Timeout in seconds for HTTP requests to the GrowthBook API.",
			},
			"insecure_skip_verify": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Description: "If true, disables SSL certificate verification for GrowthBook API requests " +
					"(not recommended for production).",
			},
			"retry_max_attempts": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "Maximum number of retry attempts for transient API errors.",
			},
			"retry_min_backoff_ms": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     500,
				Description: "Minimum backoff (in milliseconds) between retries.",
			},
			"retry_max_backoff_ms": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5000,
				Description: "Maximum backoff (in milliseconds) between retries.",
			},
			"query_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: "Maximum number of items to fetch per page for paginated API requests.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"growthbook_project":        resourceProject(),
			"growthbook_feature":        resourceFeature(),
			"growthbook_environment":    resourceEnvironment(),
			"growthbook_sdk_connection": resourceSDKConnection(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"growthbook_project":        dataSourceProject(),
			"growthbook_environment":    dataSourceEnvironment(),
			"growthbook_feature":        dataSourceFeature(),
			"growthbook_sdk_connection": dataSourceSDKConnection(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := d.Get("api_key").(string)
	if apiKey == "" {
		return nil, diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Missing GrowthBook API key",
				Detail: "The 'api_key' property must be set in the provider configuration " +
					"or via the GROWTHBOOK_API_KEY environment variable.",
			},
		}
	}
	baseURL := d.Get("api_url").(string)
	timeout := d.Get("http_timeout").(int)
	insecure := d.Get("insecure_skip_verify").(bool)

	retryMaxAttempts := d.Get("retry_max_attempts").(int)
	retryMinBackoff := d.Get("retry_min_backoff_ms").(int)
	retryMaxBackoff := d.Get("retry_max_backoff_ms").(int)
	queryLimit := d.Get("query_limit").(int)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}
	httpClient := &http.Client{
		Timeout:   time.Duration(timeout) * time.Second,
		Transport: transport,
	}

	client := growthbookapi.NewClient(
		baseURL,
		apiKey,
		growthbookapi.WithHTTPClient(httpClient),
		growthbookapi.WithBackoff(growthbookapi.BackoffConfig{
			MaxRetries:      retryMaxAttempts,
			InitialInterval: time.Duration(retryMinBackoff) * time.Millisecond,
			Multiplier:      2.0,
			MaxInterval:     time.Duration(retryMaxBackoff) * time.Millisecond,
		}),
		growthbookapi.WithPageLimit(queryLimit),
	)

	return client, nil
}

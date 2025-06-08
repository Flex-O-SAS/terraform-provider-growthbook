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
	)

	return client, nil
}

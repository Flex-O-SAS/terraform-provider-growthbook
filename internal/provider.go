// Package internal provides the Terraform provider implementation for GrowthBook.
package internal

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-growthbook/internal/growthbookapi"
)

var _ provider.Provider = &growthbookProvider{}

// New returns a new GrowthBook provider.
func New() provider.Provider {
	return &growthbookProvider{}
}

type growthbookProvider struct{}

type growthbookProviderModel struct {
	APIKey             types.String `tfsdk:"api_key"`
	APIURL             types.String `tfsdk:"api_url"`
	HTTPTimeout        types.Int64  `tfsdk:"http_timeout"`
	InsecureSkipVerify types.Bool   `tfsdk:"insecure_skip_verify"`
	RetryMaxAttempts   types.Int64  `tfsdk:"retry_max_attempts"`
	RetryMinBackoffMs  types.Int64  `tfsdk:"retry_min_backoff_ms"`
	RetryMaxBackoffMs  types.Int64  `tfsdk:"retry_max_backoff_ms"`
	QueryLimit         types.Int64  `tfsdk:"query_limit"`
}

func (p *growthbookProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "growthbook"
}

func (p *growthbookProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The GrowthBook API key. Can also be set via GROWTHBOOK_API_KEY env var.",
			},
			"api_url": schema.StringAttribute{
				Optional:    true,
				Description: "The GrowthBook API base URL. Can also be set via GROWTHBOOK_API_URL env var.",
			},
			"http_timeout": schema.Int64Attribute{
				Optional:    true,
				Description: "Timeout in seconds for HTTP requests to the GrowthBook API.",
			},
			"insecure_skip_verify": schema.BoolAttribute{
				Optional: true,
				Description: "If true, disables SSL certificate verification for GrowthBook API requests " +
					"(not recommended for production).",
			},
			"retry_max_attempts": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum number of retry attempts for transient API errors.",
			},
			"retry_min_backoff_ms": schema.Int64Attribute{
				Optional:    true,
				Description: "Minimum backoff (in milliseconds) between retries.",
			},
			"retry_max_backoff_ms": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum backoff (in milliseconds) between retries.",
			},
			"query_limit": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum number of items to fetch per page for paginated API requests.",
			},
		},
	}
}

func (p *growthbookProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config growthbookProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiKey := config.APIKey.ValueString()
	if apiKey == "" {
		apiKey = os.Getenv("GROWTHBOOK_API_KEY")
	}
	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing GrowthBook API key",
			"The 'api_key' property must be set in the provider configuration or via the GROWTHBOOK_API_KEY environment variable.",
		)
		return
	}

	apiURL := config.APIURL.ValueString()
	if apiURL == "" {
		apiURL = os.Getenv("GROWTHBOOK_API_URL")
	}
	if apiURL == "" {
		apiURL = "https://api.growthbook.io/api/v1"
	}

	timeout := int64(60)
	if !config.HTTPTimeout.IsNull() && !config.HTTPTimeout.IsUnknown() {
		timeout = config.HTTPTimeout.ValueInt64()
	}

	insecure := false
	if !config.InsecureSkipVerify.IsNull() && !config.InsecureSkipVerify.IsUnknown() {
		insecure = config.InsecureSkipVerify.ValueBool()
	}

	retryMaxAttempts := int64(5)
	if !config.RetryMaxAttempts.IsNull() && !config.RetryMaxAttempts.IsUnknown() {
		retryMaxAttempts = config.RetryMaxAttempts.ValueInt64()
	}

	retryMinBackoff := int64(500)
	if !config.RetryMinBackoffMs.IsNull() && !config.RetryMinBackoffMs.IsUnknown() {
		retryMinBackoff = config.RetryMinBackoffMs.ValueInt64()
	}

	retryMaxBackoff := int64(5000)
	if !config.RetryMaxBackoffMs.IsNull() && !config.RetryMaxBackoffMs.IsUnknown() {
		retryMaxBackoff = config.RetryMaxBackoffMs.ValueInt64()
	}

	queryLimit := int64(100)
	if !config.QueryLimit.IsNull() && !config.QueryLimit.IsUnknown() {
		queryLimit = config.QueryLimit.ValueInt64()
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure}, //nolint:gosec
	}
	httpClient := &http.Client{
		Timeout:   time.Duration(timeout) * time.Second,
		Transport: transport,
	}

	client := growthbookapi.NewClient(
		apiURL,
		apiKey,
		growthbookapi.WithHTTPClient(httpClient),
		growthbookapi.WithBackoff(growthbookapi.BackoffConfig{
			MaxRetries:      int(retryMaxAttempts),
			InitialInterval: time.Duration(retryMinBackoff) * time.Millisecond,
			Multiplier:      2.0,
			MaxInterval:     time.Duration(retryMaxBackoff) * time.Millisecond,
		}),
		growthbookapi.WithPageLimit(int(queryLimit)),
	)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *growthbookProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		newProjectResource,
		newFeatureResource,
		newEnvironmentResource,
		newSDKConnectionResource,
		newAttributeResource,
		newSavedGroupResource,
		newSegmentResource,
		newMetricResource,
		newDimensionResource,
		newFactTableResource,
		newFactMetricResource,
		newNamespaceResource,
		newExperimentResource,
	}
}

func (p *growthbookProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newProjectDataSource,
		newEnvironmentDataSource,
		newFeatureDataSource,
		newSDKConnectionDataSource,
		newAttributeDataSource,
		newSavedGroupDataSource,
		newSegmentDataSource,
		newMetricDataSource,
		newDimensionDataSource,
		newFactTableDataSource,
		newFactMetricDataSource,
		newNamespaceDataSource,
		newExperimentDataSource,
		newDataSourceDataSource,
	}
}

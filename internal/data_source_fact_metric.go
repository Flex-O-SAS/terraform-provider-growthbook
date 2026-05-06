package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-growthbook/internal/growthbookapi"
)

var _ datasource.DataSource = &factMetricDataSource{}

func newFactMetricDataSource() datasource.DataSource {
	return &factMetricDataSource{}
}

type factMetricDataSource struct {
	client *growthbookapi.Client
}

type factMetricDataSourceModel = factMetricModel

func (d *factMetricDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fact_metric"
}

func (d *factMetricDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name":                           schema.StringAttribute{Required: true},
			"id":                             schema.StringAttribute{Computed: true},
			"metric_type":                    schema.StringAttribute{Computed: true},
			"numerator":                      schema.StringAttribute{Computed: true},
			"datasource":                     schema.StringAttribute{Computed: true},
			"description":                    schema.StringAttribute{Computed: true},
			"owner":                          schema.StringAttribute{Computed: true},
			"projects":                       schema.ListAttribute{ElementType: types.StringType, Computed: true},
			"tags":                           schema.ListAttribute{ElementType: types.StringType, Computed: true},
			"denominator":                    schema.StringAttribute{Computed: true},
			"inverse":                        schema.BoolAttribute{Computed: true},
			"capping_settings":               schema.StringAttribute{Computed: true},
			"window_settings":                schema.StringAttribute{Computed: true},
			"prior_settings":                 schema.StringAttribute{Computed: true},
			"regression_adjustment_settings": schema.StringAttribute{Computed: true},
			"risk_threshold_success":         schema.Float64Attribute{Computed: true},
			"risk_threshold_danger":          schema.Float64Attribute{Computed: true},
			"min_percent_change":             schema.Float64Attribute{Computed: true},
			"max_percent_change":             schema.Float64Attribute{Computed: true},
			"min_sample_size":                schema.Float64Attribute{Computed: true},
			"target_mde":                     schema.Float64Attribute{Computed: true},
			"managed_by":                     schema.StringAttribute{Computed: true},
			"archived":                       schema.BoolAttribute{Computed: true},
			"date_created":                   schema.StringAttribute{Computed: true},
			"date_updated":                   schema.StringAttribute{Computed: true},
		},
	}
}

func (d *factMetricDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*growthbookapi.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data type", "Expected *growthbookapi.Client")
		return
	}
	d.client = client
}

func (d *factMetricDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data factMetricDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := d.client.FindFactMetricByName(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to find GrowthBook fact metric by name", err.Error())
		return
	}

	factMetricModelFromAPI((*factMetricModel)(&data), item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

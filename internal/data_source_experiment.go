package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-growthbook/internal/growthbookapi"
)

var _ datasource.DataSource = &experimentDataSource{}

func newExperimentDataSource() datasource.DataSource {
	return &experimentDataSource{}
}

type experimentDataSource struct {
	client *growthbookapi.Client
}

type experimentDataSourceModel = experimentModel

func (d *experimentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_experiment"
}

func (d *experimentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name":                schema.StringAttribute{Required: true},
			"id":                  schema.StringAttribute{Computed: true},
			"tracking_key":        schema.StringAttribute{Computed: true},
			"variations":          schema.StringAttribute{Computed: true},
			"type":                schema.StringAttribute{Computed: true},
			"project":             schema.StringAttribute{Computed: true},
			"hypothesis":          schema.StringAttribute{Computed: true},
			"description":         schema.StringAttribute{Computed: true},
			"tags":                schema.ListAttribute{ElementType: types.StringType, Computed: true},
			"owner":               schema.StringAttribute{Computed: true},
			"archived":            schema.BoolAttribute{Computed: true},
			"status":              schema.StringAttribute{Computed: true},
			"auto_refresh":        schema.BoolAttribute{Computed: true},
			"hash_attribute":      schema.StringAttribute{Computed: true},
			"fallback_attribute":  schema.StringAttribute{Computed: true},
			"datasource_id":       schema.StringAttribute{Computed: true},
			"assignment_query_id": schema.StringAttribute{Computed: true},
			"segment_id":          schema.StringAttribute{Computed: true},
			"metrics":             schema.ListAttribute{ElementType: types.StringType, Computed: true},
			"secondary_metrics":   schema.ListAttribute{ElementType: types.StringType, Computed: true},
			"guardrail_metrics":   schema.ListAttribute{ElementType: types.StringType, Computed: true},
			"stats_engine":        schema.StringAttribute{Computed: true},
			"phases":              schema.StringAttribute{Computed: true},
			"date_created":        schema.StringAttribute{Computed: true},
			"date_updated":        schema.StringAttribute{Computed: true},
		},
	}
}

func (d *experimentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *experimentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data experimentDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := d.client.FindExperimentByName(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to find GrowthBook experiment by name", err.Error())
		return
	}

	experimentModelFromAPI((*experimentModel)(&data), item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

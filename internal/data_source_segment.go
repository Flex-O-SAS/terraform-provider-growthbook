package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-growthbook/internal/growthbookapi"
)

var _ datasource.DataSource = &segmentDataSource{}

func newSegmentDataSource() datasource.DataSource {
	return &segmentDataSource{}
}

type segmentDataSource struct {
	client *growthbookapi.Client
}

type segmentDataSourceModel = segmentModel

func (d *segmentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_segment"
}

func (d *segmentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name":            schema.StringAttribute{Required: true},
			"id":              schema.StringAttribute{Computed: true},
			"type":            schema.StringAttribute{Computed: true},
			"datasource_id":   schema.StringAttribute{Computed: true},
			"identifier_type": schema.StringAttribute{Computed: true},
			"owner":           schema.StringAttribute{Computed: true},
			"description":     schema.StringAttribute{Computed: true},
			"query":           schema.StringAttribute{Computed: true},
			"fact_table_id":   schema.StringAttribute{Computed: true},
			"projects":        schema.ListAttribute{ElementType: types.StringType, Computed: true},
			"managed_by":      schema.StringAttribute{Computed: true},
			"date_created":    schema.StringAttribute{Computed: true},
			"date_updated":    schema.StringAttribute{Computed: true},
		},
	}
}

func (d *segmentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *segmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data segmentDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := d.client.FindSegmentByName(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to find GrowthBook segment by name", err.Error())
		return
	}

	segmentModelFromAPI((*segmentModel)(&data), item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-growthbook/internal/growthbookapi"
)

var _ datasource.DataSource = &factTableDataSource{}

func newFactTableDataSource() datasource.DataSource {
	return &factTableDataSource{}
}

type factTableDataSource struct {
	client *growthbookapi.Client
}

type factTableDataSourceModel = factTableModel

func (d *factTableDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fact_table"
}

func (d *factTableDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name":          schema.StringAttribute{Required: true},
			"id":            schema.StringAttribute{Computed: true},
			"datasource":    schema.StringAttribute{Computed: true},
			"user_id_types": schema.ListAttribute{ElementType: types.StringType, Computed: true},
			"sql":           schema.StringAttribute{Computed: true},
			"description":   schema.StringAttribute{Computed: true},
			"owner":         schema.StringAttribute{Computed: true},
			"projects":      schema.ListAttribute{ElementType: types.StringType, Computed: true},
			"tags":          schema.ListAttribute{ElementType: types.StringType, Computed: true},
			"event_name":    schema.StringAttribute{Computed: true},
			"managed_by":    schema.StringAttribute{Computed: true},
			"archived":      schema.BoolAttribute{Computed: true},
			"date_created":  schema.StringAttribute{Computed: true},
			"date_updated":  schema.StringAttribute{Computed: true},
		},
	}
}

func (d *factTableDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *factTableDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data factTableDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := d.client.FindFactTableByName(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to find GrowthBook fact table by name", err.Error())
		return
	}

	factTableModelFromAPI((*factTableModel)(&data), item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

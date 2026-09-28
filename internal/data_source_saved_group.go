package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-growthbook/internal/growthbookapi"
)

var _ datasource.DataSource = &savedGroupDataSource{}

func newSavedGroupDataSource() datasource.DataSource {
	return &savedGroupDataSource{}
}

type savedGroupDataSource struct {
	client *growthbookapi.Client
}

type savedGroupDataSourceModel = savedGroupModel

func (d *savedGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_saved_group"
}

func (d *savedGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name":          schema.StringAttribute{Required: true},
			"id":            schema.StringAttribute{Computed: true},
			"type":          schema.StringAttribute{Computed: true},
			"condition":     schema.StringAttribute{Computed: true},
			"attribute_key": schema.StringAttribute{Computed: true},
			"values":        schema.ListAttribute{ElementType: types.StringType, Computed: true},
			"owner":         schema.StringAttribute{Computed: true},
			"projects":      schema.ListAttribute{ElementType: types.StringType, Computed: true},
			"description":   schema.StringAttribute{Computed: true},
			"date_created":  schema.StringAttribute{Computed: true},
			"date_updated":  schema.StringAttribute{Computed: true},
		},
	}
}

func (d *savedGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *savedGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data savedGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := d.client.FindSavedGroupByName(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to find GrowthBook saved group by name", err.Error())
		return
	}

	savedGroupModelFromAPI((*savedGroupModel)(&data), item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

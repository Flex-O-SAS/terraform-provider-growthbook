package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"terraform-provider-growthbook/internal/growthbookapi"
)

var _ datasource.DataSource = &namespaceDataSource{}

func newNamespaceDataSource() datasource.DataSource {
	return &namespaceDataSource{}
}

type namespaceDataSource struct {
	client *growthbookapi.Client
}

type namespaceDataSourceModel = namespaceModel

func (d *namespaceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace"
}

func (d *namespaceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"display_name":   schema.StringAttribute{Required: true},
			"id":             schema.StringAttribute{Computed: true},
			"description":    schema.StringAttribute{Computed: true},
			"status":         schema.StringAttribute{Computed: true},
			"format":         schema.StringAttribute{Computed: true},
			"hash_attribute": schema.StringAttribute{Computed: true},
			"seed":           schema.StringAttribute{Computed: true},
		},
	}
}

func (d *namespaceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *namespaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data namespaceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := d.client.FindNamespaceByDisplayName(ctx, data.DisplayName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to find GrowthBook namespace by display name", err.Error())
		return
	}

	namespaceModelFromAPI((*namespaceModel)(&data), item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

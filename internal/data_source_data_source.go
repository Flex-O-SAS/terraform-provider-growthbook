package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-growthbook/internal/growthbookapi"
)

var _ datasource.DataSource = &dataSourceDataSource{}

func newDataSourceDataSource() datasource.DataSource {
	return &dataSourceDataSource{}
}

type dataSourceDataSource struct {
	client *growthbookapi.Client
}

type dataSourceDataModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Type         types.String `tfsdk:"type"`
	Description  types.String `tfsdk:"description"`
	ProjectIDs   types.List   `tfsdk:"project_ids"`
	EventTracker types.String `tfsdk:"event_tracker"`
	DateCreated  types.String `tfsdk:"date_created"`
	DateUpdated  types.String `tfsdk:"date_updated"`
}

func (d *dataSourceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_source"
}

func (d *dataSourceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":            schema.StringAttribute{Required: true},
			"name":          schema.StringAttribute{Computed: true},
			"type":          schema.StringAttribute{Computed: true},
			"description":   schema.StringAttribute{Computed: true},
			"project_ids":   schema.ListAttribute{ElementType: types.StringType, Computed: true},
			"event_tracker": schema.StringAttribute{Computed: true},
			"date_created":  schema.StringAttribute{Computed: true},
			"date_updated":  schema.StringAttribute{Computed: true},
		},
	}
}

func (d *dataSourceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dataSourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dataSourceDataModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := d.client.GetDataSource(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to find GrowthBook data source by ID", err.Error())
		return
	}

	data.ID = types.StringValue(item.ID)
	data.Name = types.StringValue(item.Name)
	data.Type = types.StringValue(item.Type)
	data.Description = types.StringValue(item.Description)
	data.ProjectIDs = stringsToList(item.ProjectIDs)
	data.EventTracker = types.StringValue(item.EventTracker)
	data.DateCreated = types.StringValue(item.DateCreated)
	data.DateUpdated = types.StringValue(item.DateUpdated)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

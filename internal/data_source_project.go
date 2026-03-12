package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-growthbook/internal/growthbookapi"
)

var _ datasource.DataSource = &projectDataSource{}

func newProjectDataSource() datasource.DataSource {
	return &projectDataSource{}
}

type projectDataSource struct {
	client *growthbookapi.Client
}

type projectDataModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	StatsEngine types.String `tfsdk:"stats_engine"`
	DateCreated types.String `tfsdk:"date_created"`
	DateUpdated types.String `tfsdk:"date_updated"`
}

func (d *projectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (d *projectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the GrowthBook project.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the GrowthBook project.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "The description of the GrowthBook project.",
			},
			"stats_engine": schema.StringAttribute{
				Computed:    true,
				Description: "The stats engine used by the project.",
			},
			"date_created": schema.StringAttribute{
				Computed:    true,
				Description: "The creation date of the project.",
			},
			"date_updated": schema.StringAttribute{
				Computed:    true,
				Description: "The last update date of the project.",
			},
		},
	}
}

func (d *projectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *projectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data projectDataModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, err := d.client.FindProjectByName(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to find GrowthBook project by name", err.Error())
		return
	}

	data.ID = types.StringValue(project.ID)
	data.Name = types.StringValue(project.Name)
	data.Description = types.StringValue(project.Description)
	data.StatsEngine = types.StringValue(project.Settings.StatsEngine)
	data.DateCreated = types.StringValue(project.DateCreated)
	data.DateUpdated = types.StringValue(project.DateUpdated)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

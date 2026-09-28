package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-growthbook/internal/growthbookapi"
)

var _ datasource.DataSource = &environmentDataSource{}

func newEnvironmentDataSource() datasource.DataSource {
	return &environmentDataSource{}
}

type environmentDataSource struct {
	client *growthbookapi.Client
}

type environmentDataModel struct {
	ID           types.String `tfsdk:"id"`
	Description  types.String `tfsdk:"description"`
	ToggleOnList types.Bool   `tfsdk:"toggle_on_list"`
	DefaultState types.Bool   `tfsdk:"default_state"`
	Projects     types.List   `tfsdk:"projects"`
}

func (d *environmentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

func (d *environmentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the GrowthBook environment.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "The description of the GrowthBook environment.",
			},
			"toggle_on_list": schema.BoolAttribute{
				Computed:    true,
				Description: "Show toggle on feature list.",
			},
			"default_state": schema.BoolAttribute{
				Computed:    true,
				Description: "Default state for new features.",
			},
			"projects": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Projects associated with the environment.",
			},
		},
	}
}

func (d *environmentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *environmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data environmentDataModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	env, err := d.client.FindEnvironmentByID(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to find GrowthBook environment by ID", err.Error())
		return
	}

	data.ID = types.StringValue(env.ID)
	data.Description = types.StringValue(env.Description)
	data.ToggleOnList = types.BoolValue(env.ToggleOnList)
	data.DefaultState = types.BoolValue(env.DefaultState)
	data.Projects = stringsToList(env.Projects)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-growthbook/internal/growthbookapi"
)

var _ datasource.DataSource = &environmentsDataSource{}

func newEnvironmentsDataSource() datasource.DataSource {
	return &environmentsDataSource{}
}

type environmentsDataSource struct {
	client *growthbookapi.Client
}

type environmentsDataModel struct {
	Environments []environmentDataModel `tfsdk:"environments"`
}

func (d *environmentsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environments"
}

func (d *environmentsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"environments": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
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
				},
			},
		},
	}
}

func (d *environmentsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *environmentsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	envs, err := d.client.ListEnvironments(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to list GrowthBook environments", err.Error())
		return
	}

	items := make([]environmentDataModel, len(envs))
	for i, env := range envs {
		items[i] = environmentDataModel{
			ID:           types.StringValue(env.ID),
			Description:  types.StringValue(env.Description),
			ToggleOnList: types.BoolValue(env.ToggleOnList),
			DefaultState: types.BoolValue(env.DefaultState),
			Projects:     stringsToList(env.Projects),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &environmentsDataModel{Environments: items})...)
}

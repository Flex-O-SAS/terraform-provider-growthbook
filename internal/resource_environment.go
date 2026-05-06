package internal

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-growthbook/internal/growthbookapi"
)

var _ resource.Resource = &environmentResource{}
var _ resource.ResourceWithImportState = &environmentResource{}

func newEnvironmentResource() resource.Resource {
	return &environmentResource{}
}

type environmentResource struct {
	client *growthbookapi.Client
}

type environmentModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	ToggleOnList types.Bool   `tfsdk:"toggle_on_list"`
	DefaultState types.Bool   `tfsdk:"default_state"`
	Projects     types.List   `tfsdk:"projects"`
}

func (r *environmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

func (r *environmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"toggle_on_list": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"default_state": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"projects": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *environmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*growthbookapi.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data type", "Expected *growthbookapi.Client")
		return
	}
	r.client = client
}

func (r *environmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data environmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	env := &growthbookapi.Environment{
		ID:           data.Name.ValueString(),
		Description:  data.Description.ValueString(),
		ToggleOnList: data.ToggleOnList.ValueBool(),
		DefaultState: data.DefaultState.ValueBool(),
	}
	projects := []string{}
	if !data.Projects.IsNull() && !data.Projects.IsUnknown() {
		resp.Diagnostics.Append(data.Projects.ElementsAs(ctx, &projects, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	env.Projects = projects

	created, err := r.client.CreateEnvironment(ctx, env)
	if err != nil {
		resp.Diagnostics.AddError("Error creating environment", err.Error())
		return
	}

	data.ID = types.StringValue(created.ID)
	data.Name = types.StringValue(created.ID)
	data.Description = types.StringValue(created.Description)
	data.ToggleOnList = types.BoolValue(created.ToggleOnList)
	data.DefaultState = types.BoolValue(created.DefaultState)
	data.Projects = stringsToList(created.Projects)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *environmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data environmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	env, err := r.client.FindEnvironmentByID(ctx, data.ID.ValueString())
	if err != nil {
		if errors.Is(err, growthbookapi.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading environment", err.Error())
		return
	}

	data.Name = types.StringValue(env.ID)
	data.Description = types.StringValue(env.Description)
	data.ToggleOnList = types.BoolValue(env.ToggleOnList)
	data.DefaultState = types.BoolValue(env.DefaultState)
	data.Projects = stringsToList(env.Projects)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *environmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data environmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state environmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	env := &growthbookapi.Environment{
		Description:  data.Description.ValueString(),
		ToggleOnList: data.ToggleOnList.ValueBool(),
		DefaultState: data.DefaultState.ValueBool(),
	}
	projects := []string{}
	if !data.Projects.IsNull() && !data.Projects.IsUnknown() {
		resp.Diagnostics.Append(data.Projects.ElementsAs(ctx, &projects, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	env.Projects = projects

	_, err := r.client.UpdateEnvironment(ctx, state.ID.ValueString(), env)
	if err != nil {
		resp.Diagnostics.AddError("Error updating environment", err.Error())
		return
	}

	data.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *environmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data environmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteEnvironment(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting environment", err.Error())
	}
}

func (r *environmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

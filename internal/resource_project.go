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

var _ resource.Resource = &projectResource{}
var _ resource.ResourceWithImportState = &projectResource{}

func newProjectResource() resource.Resource {
	return &projectResource{}
}

type projectResource struct {
	client *growthbookapi.Client
}

type projectModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	StatsEngine types.String `tfsdk:"stats_engine"`
	DateCreated types.String `tfsdk:"date_created"`
	DateUpdated types.String `tfsdk:"date_updated"`
}

func (r *projectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *projectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"stats_engine": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"date_created": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"date_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *projectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *projectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data projectModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	project := &growthbookapi.Project{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
	}
	if !data.StatsEngine.IsNull() && !data.StatsEngine.IsUnknown() {
		project.Settings.StatsEngine = data.StatsEngine.ValueString()
	}

	created, err := r.client.CreateProject(ctx, project)
	if err != nil {
		resp.Diagnostics.AddError("Error creating project", err.Error())
		return
	}

	data.ID = types.StringValue(created.ID)
	data.Name = types.StringValue(created.Name)
	data.Description = types.StringValue(created.Description)
	data.StatsEngine = types.StringValue(created.Settings.StatsEngine)
	data.DateCreated = types.StringValue(created.DateCreated)
	data.DateUpdated = types.StringValue(created.DateUpdated)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *projectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data projectModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, err := r.client.GetProject(ctx, data.ID.ValueString())
	if err != nil {
		if errors.Is(err, growthbookapi.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading project", err.Error())
		return
	}

	data.Name = types.StringValue(project.Name)
	data.Description = types.StringValue(project.Description)
	data.StatsEngine = types.StringValue(project.Settings.StatsEngine)
	data.DateCreated = types.StringValue(project.DateCreated)
	data.DateUpdated = types.StringValue(project.DateUpdated)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *projectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data projectModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state projectModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	project := &growthbookapi.Project{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
	}
	if !data.StatsEngine.IsNull() && !data.StatsEngine.IsUnknown() {
		project.Settings.StatsEngine = data.StatsEngine.ValueString()
	}

	updated, err := r.client.UpdateProject(ctx, state.ID.ValueString(), project)
	if err != nil {
		resp.Diagnostics.AddError("Error updating project", err.Error())
		return
	}

	data.ID = state.ID
	data.Name = types.StringValue(updated.Name)
	data.Description = types.StringValue(updated.Description)
	data.StatsEngine = types.StringValue(updated.Settings.StatsEngine)
	data.DateCreated = state.DateCreated
	data.DateUpdated = types.StringValue(updated.DateUpdated)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *projectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data projectModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteProject(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting project", err.Error())
	}
}

func (r *projectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

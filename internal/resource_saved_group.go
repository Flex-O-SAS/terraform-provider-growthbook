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

var _ resource.Resource = &savedGroupResource{}
var _ resource.ResourceWithImportState = &savedGroupResource{}

func newSavedGroupResource() resource.Resource {
	return &savedGroupResource{}
}

type savedGroupResource struct {
	client *growthbookapi.Client
}

type savedGroupModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Type         types.String `tfsdk:"type"`
	Condition    types.String `tfsdk:"condition"`
	AttributeKey types.String `tfsdk:"attribute_key"`
	Values       types.List   `tfsdk:"values"`
	Owner        types.String `tfsdk:"owner"`
	Projects     types.List   `tfsdk:"projects"`
	Description  types.String `tfsdk:"description"`
	DateCreated  types.String `tfsdk:"date_created"`
	DateUpdated  types.String `tfsdk:"date_updated"`
}

func (r *savedGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_saved_group"
}

func (r *savedGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name":          schema.StringAttribute{Required: true},
			"type":          schema.StringAttribute{Optional: true, Computed: true},
			"condition":     schema.StringAttribute{Optional: true, Computed: true},
			"attribute_key": schema.StringAttribute{Optional: true, Computed: true},
			"values":        schema.ListAttribute{ElementType: types.StringType, Optional: true, Computed: true},
			"owner":         schema.StringAttribute{Optional: true, Computed: true},
			"projects":      schema.ListAttribute{ElementType: types.StringType, Optional: true, Computed: true},
			"description":   schema.StringAttribute{Optional: true, Computed: true},
			"date_created":  schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"date_updated":  schema.StringAttribute{Computed: true},
		},
	}
}

func (r *savedGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func savedGroupFromModel(data savedGroupModel) *growthbookapi.SavedGroup {
	return &growthbookapi.SavedGroup{
		Name:         data.Name.ValueString(),
		Type:         data.Type.ValueString(),
		Condition:    data.Condition.ValueString(),
		AttributeKey: data.AttributeKey.ValueString(),
		Values:       listToStrings(data.Values),
		Owner:        data.Owner.ValueString(),
		Projects:     listToStrings(data.Projects),
		Description:  data.Description.ValueString(),
	}
}

func savedGroupModelFromAPI(data *savedGroupModel, item *growthbookapi.SavedGroup) {
	data.ID = types.StringValue(item.ID)
	data.Name = types.StringValue(item.Name)
	data.Type = types.StringValue(item.Type)
	data.Condition = types.StringValue(item.Condition)
	data.AttributeKey = types.StringValue(item.AttributeKey)
	data.Values = stringsToList(item.Values)
	data.Owner = types.StringValue(item.Owner)
	data.Projects = stringsToList(item.Projects)
	data.Description = types.StringValue(item.Description)
	data.DateCreated = types.StringValue(item.DateCreated)
	data.DateUpdated = types.StringValue(item.DateUpdated)
}

func (r *savedGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data savedGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.CreateSavedGroup(ctx, savedGroupFromModel(data))
	if err != nil {
		resp.Diagnostics.AddError("Error creating saved group", err.Error())
		return
	}

	savedGroupModelFromAPI(&data, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *savedGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data savedGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.GetSavedGroup(ctx, data.ID.ValueString())
	if err != nil {
		if errors.Is(err, growthbookapi.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading saved group", err.Error())
		return
	}

	savedGroupModelFromAPI(&data, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *savedGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data savedGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state savedGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.UpdateSavedGroup(ctx, state.ID.ValueString(), savedGroupFromModel(data))
	if err != nil {
		resp.Diagnostics.AddError("Error updating saved group", err.Error())
		return
	}

	savedGroupModelFromAPI(&data, item)
	data.ID = state.ID
	if data.DateCreated.ValueString() == "" {
		data.DateCreated = state.DateCreated
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *savedGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data savedGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSavedGroup(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting saved group", err.Error())
	}
}

func (r *savedGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

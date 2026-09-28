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

var _ resource.Resource = &factTableResource{}
var _ resource.ResourceWithImportState = &factTableResource{}

func newFactTableResource() resource.Resource {
	return &factTableResource{}
}

type factTableResource struct {
	client *growthbookapi.Client
}

type factTableModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Datasource  types.String `tfsdk:"datasource"`
	UserIDTypes types.List   `tfsdk:"user_id_types"`
	SQL         types.String `tfsdk:"sql"`
	Description types.String `tfsdk:"description"`
	Owner       types.String `tfsdk:"owner"`
	Projects    types.List   `tfsdk:"projects"`
	Tags        types.List   `tfsdk:"tags"`
	EventName   types.String `tfsdk:"event_name"`
	ManagedBy   types.String `tfsdk:"managed_by"`
	Archived    types.Bool   `tfsdk:"archived"`
	DateCreated types.String `tfsdk:"date_created"`
	DateUpdated types.String `tfsdk:"date_updated"`
}

func (r *factTableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fact_table"
}

func (r *factTableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":            schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"name":          schema.StringAttribute{Required: true},
			"datasource":    schema.StringAttribute{Required: true},
			"user_id_types": schema.ListAttribute{ElementType: types.StringType, Required: true},
			"sql":           schema.StringAttribute{Required: true},
			"description":   schema.StringAttribute{Optional: true, Computed: true},
			"owner":         schema.StringAttribute{Optional: true, Computed: true},
			"projects":      schema.ListAttribute{ElementType: types.StringType, Optional: true, Computed: true},
			"tags":          schema.ListAttribute{ElementType: types.StringType, Optional: true, Computed: true},
			"event_name":    schema.StringAttribute{Optional: true, Computed: true},
			"managed_by":    schema.StringAttribute{Optional: true, Computed: true},
			"archived":      schema.BoolAttribute{Optional: true, Computed: true},
			"date_created":  schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"date_updated":  schema.StringAttribute{Computed: true},
		},
	}
}

func (r *factTableResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func factTableFromModel(data factTableModel) *growthbookapi.FactTable {
	return &growthbookapi.FactTable{
		Name:        data.Name.ValueString(),
		Datasource:  data.Datasource.ValueString(),
		UserIDTypes: listToStrings(data.UserIDTypes),
		SQL:         data.SQL.ValueString(),
		Description: data.Description.ValueString(),
		Owner:       data.Owner.ValueString(),
		Projects:    listToStrings(data.Projects),
		Tags:        listToStrings(data.Tags),
		EventName:   data.EventName.ValueString(),
		ManagedBy:   data.ManagedBy.ValueString(),
		Archived:    data.Archived.ValueBool(),
	}
}

func factTableModelFromAPI(data *factTableModel, item *growthbookapi.FactTable) {
	data.ID = types.StringValue(item.ID)
	data.Name = types.StringValue(item.Name)
	data.Datasource = types.StringValue(item.Datasource)
	data.UserIDTypes = stringsToList(item.UserIDTypes)
	data.SQL = types.StringValue(item.SQL)
	data.Description = types.StringValue(item.Description)
	data.Owner = types.StringValue(item.Owner)
	data.Projects = stringsToList(item.Projects)
	data.Tags = stringsToList(item.Tags)
	data.EventName = types.StringValue(item.EventName)
	data.ManagedBy = types.StringValue(item.ManagedBy)
	data.Archived = types.BoolValue(item.Archived)
	data.DateCreated = types.StringValue(item.DateCreated)
	data.DateUpdated = types.StringValue(item.DateUpdated)
}

func (r *factTableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data factTableModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.CreateFactTable(ctx, factTableFromModel(data))
	if err != nil {
		resp.Diagnostics.AddError("Error creating fact table", err.Error())
		return
	}

	factTableModelFromAPI(&data, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *factTableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data factTableModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.GetFactTable(ctx, data.ID.ValueString())
	if err != nil {
		if errors.Is(err, growthbookapi.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading fact table", err.Error())
		return
	}

	factTableModelFromAPI(&data, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *factTableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data factTableModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state factTableModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.UpdateFactTable(ctx, state.ID.ValueString(), factTableFromModel(data))
	if err != nil {
		resp.Diagnostics.AddError("Error updating fact table", err.Error())
		return
	}

	factTableModelFromAPI(&data, item)
	data.ID = state.ID
	if data.DateCreated.ValueString() == "" {
		data.DateCreated = state.DateCreated
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *factTableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data factTableModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteFactTable(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting fact table", err.Error())
	}
}

func (r *factTableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

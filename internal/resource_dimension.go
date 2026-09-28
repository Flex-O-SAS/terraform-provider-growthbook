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

var _ resource.Resource = &dimensionResource{}
var _ resource.ResourceWithImportState = &dimensionResource{}

func newDimensionResource() resource.Resource {
	return &dimensionResource{}
}

type dimensionResource struct {
	client *growthbookapi.Client
}

type dimensionModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	DatasourceID   types.String `tfsdk:"datasource_id"`
	IdentifierType types.String `tfsdk:"identifier_type"`
	Query          types.String `tfsdk:"query"`
	Description    types.String `tfsdk:"description"`
	Owner          types.String `tfsdk:"owner"`
	ManagedBy      types.String `tfsdk:"managed_by"`
	DateCreated    types.String `tfsdk:"date_created"`
	DateUpdated    types.String `tfsdk:"date_updated"`
}

func (r *dimensionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dimension"
}

func (r *dimensionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":              schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"name":            schema.StringAttribute{Required: true},
			"datasource_id":   schema.StringAttribute{Required: true},
			"identifier_type": schema.StringAttribute{Required: true},
			"query":           schema.StringAttribute{Required: true},
			"description":     schema.StringAttribute{Optional: true, Computed: true},
			"owner":           schema.StringAttribute{Optional: true, Computed: true},
			"managed_by":      schema.StringAttribute{Optional: true, Computed: true},
			"date_created":    schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"date_updated":    schema.StringAttribute{Computed: true},
		},
	}
}

func (r *dimensionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func dimensionFromModel(data dimensionModel) *growthbookapi.Dimension {
	return &growthbookapi.Dimension{
		Name:           data.Name.ValueString(),
		DatasourceID:   data.DatasourceID.ValueString(),
		IdentifierType: data.IdentifierType.ValueString(),
		Query:          data.Query.ValueString(),
		Description:    data.Description.ValueString(),
		Owner:          data.Owner.ValueString(),
		ManagedBy:      data.ManagedBy.ValueString(),
	}
}

func dimensionModelFromAPI(data *dimensionModel, item *growthbookapi.Dimension) {
	data.ID = types.StringValue(item.ID)
	data.Name = types.StringValue(item.Name)
	data.DatasourceID = types.StringValue(item.DatasourceID)
	data.IdentifierType = types.StringValue(item.IdentifierType)
	data.Query = types.StringValue(item.Query)
	data.Description = types.StringValue(item.Description)
	data.Owner = types.StringValue(item.Owner)
	data.ManagedBy = types.StringValue(item.ManagedBy)
	data.DateCreated = types.StringValue(item.DateCreated)
	data.DateUpdated = types.StringValue(item.DateUpdated)
}

func (r *dimensionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data dimensionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.CreateDimension(ctx, dimensionFromModel(data))
	if err != nil {
		resp.Diagnostics.AddError("Error creating dimension", err.Error())
		return
	}

	dimensionModelFromAPI(&data, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dimensionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data dimensionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.GetDimension(ctx, data.ID.ValueString())
	if err != nil {
		if errors.Is(err, growthbookapi.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading dimension", err.Error())
		return
	}

	dimensionModelFromAPI(&data, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dimensionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data dimensionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state dimensionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.UpdateDimension(ctx, state.ID.ValueString(), dimensionFromModel(data))
	if err != nil {
		resp.Diagnostics.AddError("Error updating dimension", err.Error())
		return
	}

	dimensionModelFromAPI(&data, item)
	data.ID = state.ID
	if data.DateCreated.ValueString() == "" {
		data.DateCreated = state.DateCreated
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *dimensionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data dimensionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteDimension(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting dimension", err.Error())
	}
}

func (r *dimensionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

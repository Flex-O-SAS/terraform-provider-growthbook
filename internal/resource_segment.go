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

var _ resource.Resource = &segmentResource{}
var _ resource.ResourceWithImportState = &segmentResource{}

func newSegmentResource() resource.Resource {
	return &segmentResource{}
}

type segmentResource struct {
	client *growthbookapi.Client
}

type segmentModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Type           types.String `tfsdk:"type"`
	DatasourceID   types.String `tfsdk:"datasource_id"`
	IdentifierType types.String `tfsdk:"identifier_type"`
	Owner          types.String `tfsdk:"owner"`
	Description    types.String `tfsdk:"description"`
	Query          types.String `tfsdk:"query"`
	FactTableID    types.String `tfsdk:"fact_table_id"`
	Projects       types.List   `tfsdk:"projects"`
	ManagedBy      types.String `tfsdk:"managed_by"`
	DateCreated    types.String `tfsdk:"date_created"`
	DateUpdated    types.String `tfsdk:"date_updated"`
}

func (r *segmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_segment"
}

func (r *segmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":              schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"name":            schema.StringAttribute{Required: true},
			"type":            schema.StringAttribute{Required: true},
			"datasource_id":   schema.StringAttribute{Required: true},
			"identifier_type": schema.StringAttribute{Required: true},
			"owner":           schema.StringAttribute{Optional: true, Computed: true},
			"description":     schema.StringAttribute{Optional: true, Computed: true},
			"query":           schema.StringAttribute{Optional: true, Computed: true},
			"fact_table_id":   schema.StringAttribute{Optional: true, Computed: true},
			"projects":        schema.ListAttribute{ElementType: types.StringType, Optional: true, Computed: true},
			"managed_by":      schema.StringAttribute{Optional: true, Computed: true},
			"date_created":    schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"date_updated":    schema.StringAttribute{Computed: true},
		},
	}
}

func (r *segmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func segmentFromModel(data segmentModel) *growthbookapi.Segment {
	return &growthbookapi.Segment{
		Name:           data.Name.ValueString(),
		Type:           data.Type.ValueString(),
		DatasourceID:   data.DatasourceID.ValueString(),
		IdentifierType: data.IdentifierType.ValueString(),
		Owner:          data.Owner.ValueString(),
		Description:    data.Description.ValueString(),
		Query:          data.Query.ValueString(),
		FactTableID:    data.FactTableID.ValueString(),
		Projects:       listToStrings(data.Projects),
		ManagedBy:      data.ManagedBy.ValueString(),
	}
}

func segmentModelFromAPI(data *segmentModel, item *growthbookapi.Segment) {
	data.ID = types.StringValue(item.ID)
	data.Name = types.StringValue(item.Name)
	data.Type = types.StringValue(item.Type)
	data.DatasourceID = types.StringValue(item.DatasourceID)
	data.IdentifierType = types.StringValue(item.IdentifierType)
	data.Owner = types.StringValue(item.Owner)
	data.Description = types.StringValue(item.Description)
	data.Query = types.StringValue(item.Query)
	data.FactTableID = types.StringValue(item.FactTableID)
	data.Projects = stringsToList(item.Projects)
	data.ManagedBy = types.StringValue(item.ManagedBy)
	data.DateCreated = types.StringValue(item.DateCreated)
	data.DateUpdated = types.StringValue(item.DateUpdated)
}

func (r *segmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data segmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.CreateSegment(ctx, segmentFromModel(data))
	if err != nil {
		resp.Diagnostics.AddError("Error creating segment", err.Error())
		return
	}

	segmentModelFromAPI(&data, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *segmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data segmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.GetSegment(ctx, data.ID.ValueString())
	if err != nil {
		if errors.Is(err, growthbookapi.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading segment", err.Error())
		return
	}

	segmentModelFromAPI(&data, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *segmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data segmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state segmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.UpdateSegment(ctx, state.ID.ValueString(), segmentFromModel(data))
	if err != nil {
		resp.Diagnostics.AddError("Error updating segment", err.Error())
		return
	}

	segmentModelFromAPI(&data, item)
	data.ID = state.ID
	if data.DateCreated.ValueString() == "" {
		data.DateCreated = state.DateCreated
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *segmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data segmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSegment(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting segment", err.Error())
	}
}

func (r *segmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

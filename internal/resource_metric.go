package internal

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-growthbook/internal/growthbookapi"
)

var _ resource.Resource = &metricResource{}
var _ resource.ResourceWithImportState = &metricResource{}

func newMetricResource() resource.Resource {
	return &metricResource{}
}

type metricResource struct {
	client *growthbookapi.Client
}

type metricModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	DatasourceID types.String `tfsdk:"datasource_id"`
	Type         types.String `tfsdk:"type"`
	Description  types.String `tfsdk:"description"`
	Owner        types.String `tfsdk:"owner"`
	Tags         types.List   `tfsdk:"tags"`
	Projects     types.List   `tfsdk:"projects"`
	Archived     types.Bool   `tfsdk:"archived"`
	Behavior     types.String `tfsdk:"behavior"`
	SQL          types.String `tfsdk:"sql"`
	ManagedBy    types.String `tfsdk:"managed_by"`
	DateCreated  types.String `tfsdk:"date_created"`
	DateUpdated  types.String `tfsdk:"date_updated"`
}

func (r *metricResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metric"
}

func (r *metricResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":            schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"name":          schema.StringAttribute{Required: true},
			"datasource_id": schema.StringAttribute{Required: true},
			"type":          schema.StringAttribute{Required: true},
			"description":   schema.StringAttribute{Optional: true, Computed: true},
			"owner":         schema.StringAttribute{Optional: true, Computed: true},
			"tags":          schema.ListAttribute{ElementType: types.StringType, Optional: true, Computed: true},
			"projects":      schema.ListAttribute{ElementType: types.StringType, Optional: true, Computed: true},
			"archived":      schema.BoolAttribute{Optional: true, Computed: true},
			"behavior":      schema.StringAttribute{Optional: true, Computed: true, Validators: []validator.String{JSONString()}},
			"sql":           schema.StringAttribute{Optional: true, Computed: true, Validators: []validator.String{JSONString()}},
			"managed_by":    schema.StringAttribute{Optional: true, Computed: true},
			"date_created":  schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"date_updated":  schema.StringAttribute{Computed: true},
		},
	}
}

func (r *metricResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func metricFromModel(data metricModel) *growthbookapi.Metric {
	return &growthbookapi.Metric{
		Name:         data.Name.ValueString(),
		DatasourceID: data.DatasourceID.ValueString(),
		Type:         data.Type.ValueString(),
		Description:  data.Description.ValueString(),
		Owner:        data.Owner.ValueString(),
		Tags:         listToStrings(data.Tags),
		Projects:     listToStrings(data.Projects),
		Archived:     data.Archived.ValueBool(),
		Behavior:     stringToRawJSON(data.Behavior.ValueString()),
		SQL:          stringToRawJSON(data.SQL.ValueString()),
		ManagedBy:    data.ManagedBy.ValueString(),
	}
}

func metricModelFromAPI(data *metricModel, item *growthbookapi.Metric) {
	data.ID = types.StringValue(item.ID)
	data.Name = types.StringValue(item.Name)
	data.DatasourceID = types.StringValue(item.DatasourceID)
	data.Type = types.StringValue(item.Type)
	data.Description = types.StringValue(item.Description)
	data.Owner = types.StringValue(item.Owner)
	data.Tags = stringsToList(item.Tags)
	data.Projects = stringsToList(item.Projects)
	data.Archived = types.BoolValue(item.Archived)
	data.Behavior = types.StringValue(rawJSONToString(item.Behavior))
	data.SQL = types.StringValue(rawJSONToString(item.SQL))
	data.ManagedBy = types.StringValue(item.ManagedBy)
	data.DateCreated = types.StringValue(item.DateCreated)
	data.DateUpdated = types.StringValue(item.DateUpdated)
}

func (r *metricResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data metricModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.CreateMetric(ctx, metricFromModel(data))
	if err != nil {
		resp.Diagnostics.AddError("Error creating metric", err.Error())
		return
	}

	metricModelFromAPI(&data, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *metricResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data metricModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.GetMetric(ctx, data.ID.ValueString())
	if err != nil {
		if errors.Is(err, growthbookapi.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading metric", err.Error())
		return
	}

	metricModelFromAPI(&data, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *metricResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data metricModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state metricModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.UpdateMetric(ctx, state.ID.ValueString(), metricFromModel(data))
	if err != nil {
		resp.Diagnostics.AddError("Error updating metric", err.Error())
		return
	}

	metricModelFromAPI(&data, item)
	data.ID = state.ID
	if data.DateCreated.ValueString() == "" {
		data.DateCreated = state.DateCreated
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *metricResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data metricModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteMetric(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting metric", err.Error())
	}
}

func (r *metricResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

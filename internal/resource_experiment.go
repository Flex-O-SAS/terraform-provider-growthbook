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

var _ resource.Resource = &experimentResource{}
var _ resource.ResourceWithImportState = &experimentResource{}

func newExperimentResource() resource.Resource {
	return &experimentResource{}
}

type experimentResource struct {
	client *growthbookapi.Client
}

type experimentModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	TrackingKey       types.String `tfsdk:"tracking_key"`
	Variations        types.String `tfsdk:"variations"`
	Type              types.String `tfsdk:"type"`
	Project           types.String `tfsdk:"project"`
	Hypothesis        types.String `tfsdk:"hypothesis"`
	Description       types.String `tfsdk:"description"`
	Tags              types.List   `tfsdk:"tags"`
	Owner             types.String `tfsdk:"owner"`
	Archived          types.Bool   `tfsdk:"archived"`
	Status            types.String `tfsdk:"status"`
	AutoRefresh       types.Bool   `tfsdk:"auto_refresh"`
	HashAttribute     types.String `tfsdk:"hash_attribute"`
	FallbackAttribute types.String `tfsdk:"fallback_attribute"`
	DatasourceID      types.String `tfsdk:"datasource_id"`
	AssignmentQueryID types.String `tfsdk:"assignment_query_id"`
	SegmentID         types.String `tfsdk:"segment_id"`
	Metrics           types.List   `tfsdk:"metrics"`
	SecondaryMetrics  types.List   `tfsdk:"secondary_metrics"`
	GuardrailMetrics  types.List   `tfsdk:"guardrail_metrics"`
	StatsEngine       types.String `tfsdk:"stats_engine"`
	Phases            types.String `tfsdk:"phases"`
	DateCreated       types.String `tfsdk:"date_created"`
	DateUpdated       types.String `tfsdk:"date_updated"`
}

func (r *experimentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_experiment"
}

func (r *experimentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                  schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"name":                schema.StringAttribute{Required: true},
			"tracking_key":        schema.StringAttribute{Required: true},
			"variations":          schema.StringAttribute{Required: true, Validators: []validator.String{JSONString()}},
			"type":                schema.StringAttribute{Optional: true, Computed: true},
			"project":             schema.StringAttribute{Optional: true, Computed: true},
			"hypothesis":          schema.StringAttribute{Optional: true, Computed: true},
			"description":         schema.StringAttribute{Optional: true, Computed: true},
			"tags":                schema.ListAttribute{ElementType: types.StringType, Optional: true, Computed: true},
			"owner":               schema.StringAttribute{Optional: true, Computed: true},
			"archived":            schema.BoolAttribute{Optional: true, Computed: true},
			"status":              schema.StringAttribute{Optional: true, Computed: true},
			"auto_refresh":        schema.BoolAttribute{Optional: true, Computed: true},
			"hash_attribute":      schema.StringAttribute{Optional: true, Computed: true},
			"fallback_attribute":  schema.StringAttribute{Optional: true, Computed: true},
			"datasource_id":       schema.StringAttribute{Optional: true, Computed: true},
			"assignment_query_id": schema.StringAttribute{Optional: true, Computed: true},
			"segment_id":          schema.StringAttribute{Optional: true, Computed: true},
			"metrics":             schema.ListAttribute{ElementType: types.StringType, Optional: true, Computed: true},
			"secondary_metrics":   schema.ListAttribute{ElementType: types.StringType, Optional: true, Computed: true},
			"guardrail_metrics":   schema.ListAttribute{ElementType: types.StringType, Optional: true, Computed: true},
			"stats_engine":        schema.StringAttribute{Optional: true, Computed: true},
			"phases":              schema.StringAttribute{Optional: true, Computed: true, Validators: []validator.String{JSONString()}},
			"date_created":        schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"date_updated":        schema.StringAttribute{Computed: true},
		},
	}
}

func (r *experimentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func experimentFromModel(data experimentModel) *growthbookapi.Experiment {
	return &growthbookapi.Experiment{
		Name:              data.Name.ValueString(),
		TrackingKey:       data.TrackingKey.ValueString(),
		Variations:        stringToRawJSON(data.Variations.ValueString()),
		Type:              data.Type.ValueString(),
		Project:           data.Project.ValueString(),
		Hypothesis:        data.Hypothesis.ValueString(),
		Description:       data.Description.ValueString(),
		Tags:              listToStrings(data.Tags),
		Owner:             data.Owner.ValueString(),
		Archived:          data.Archived.ValueBool(),
		Status:            data.Status.ValueString(),
		AutoRefresh:       data.AutoRefresh.ValueBool(),
		HashAttribute:     data.HashAttribute.ValueString(),
		FallbackAttribute: data.FallbackAttribute.ValueString(),
		DatasourceID:      data.DatasourceID.ValueString(),
		AssignmentQueryID: data.AssignmentQueryID.ValueString(),
		SegmentID:         data.SegmentID.ValueString(),
		Metrics:           listToStrings(data.Metrics),
		SecondaryMetrics:  listToStrings(data.SecondaryMetrics),
		GuardrailMetrics:  listToStrings(data.GuardrailMetrics),
		StatsEngine:       data.StatsEngine.ValueString(),
		Phases:            stringToRawJSON(data.Phases.ValueString()),
	}
}

func experimentModelFromAPI(data *experimentModel, item *growthbookapi.Experiment) {
	data.ID = types.StringValue(item.ID)
	data.Name = types.StringValue(item.Name)
	data.TrackingKey = types.StringValue(item.TrackingKey)
	data.Variations = types.StringValue(rawJSONToString(item.Variations))
	data.Type = types.StringValue(item.Type)
	data.Project = types.StringValue(item.Project)
	data.Hypothesis = types.StringValue(item.Hypothesis)
	data.Description = types.StringValue(item.Description)
	data.Tags = stringsToList(item.Tags)
	data.Owner = types.StringValue(item.Owner)
	data.Archived = types.BoolValue(item.Archived)
	data.Status = types.StringValue(item.Status)
	data.AutoRefresh = types.BoolValue(item.AutoRefresh)
	data.HashAttribute = types.StringValue(item.HashAttribute)
	data.FallbackAttribute = types.StringValue(item.FallbackAttribute)
	data.DatasourceID = types.StringValue(item.DatasourceID)
	data.AssignmentQueryID = types.StringValue(item.AssignmentQueryID)
	data.SegmentID = types.StringValue(item.SegmentID)
	data.Metrics = stringsToList(item.Metrics)
	data.SecondaryMetrics = stringsToList(item.SecondaryMetrics)
	data.GuardrailMetrics = stringsToList(item.GuardrailMetrics)
	data.StatsEngine = types.StringValue(item.StatsEngine)
	data.Phases = types.StringValue(rawJSONToString(item.Phases))
	data.DateCreated = types.StringValue(item.DateCreated)
	data.DateUpdated = types.StringValue(item.DateUpdated)
}

func (r *experimentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data experimentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.CreateExperiment(ctx, experimentFromModel(data))
	if err != nil {
		resp.Diagnostics.AddError("Error creating experiment", err.Error())
		return
	}

	experimentModelFromAPI(&data, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *experimentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data experimentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.GetExperiment(ctx, data.ID.ValueString())
	if err != nil {
		if errors.Is(err, growthbookapi.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading experiment", err.Error())
		return
	}

	experimentModelFromAPI(&data, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *experimentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data experimentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state experimentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.UpdateExperiment(ctx, state.ID.ValueString(), experimentFromModel(data))
	if err != nil {
		resp.Diagnostics.AddError("Error updating experiment", err.Error())
		return
	}

	experimentModelFromAPI(&data, item)
	data.ID = state.ID
	if data.DateCreated.ValueString() == "" {
		data.DateCreated = state.DateCreated
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *experimentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data experimentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteExperiment(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting experiment", err.Error())
	}
}

func (r *experimentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

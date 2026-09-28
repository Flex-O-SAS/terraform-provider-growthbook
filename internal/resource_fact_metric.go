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

var _ resource.Resource = &factMetricResource{}
var _ resource.ResourceWithImportState = &factMetricResource{}

func newFactMetricResource() resource.Resource {
	return &factMetricResource{}
}

type factMetricResource struct {
	client *growthbookapi.Client
}

type factMetricModel struct {
	ID                           types.String  `tfsdk:"id"`
	Name                         types.String  `tfsdk:"name"`
	MetricType                   types.String  `tfsdk:"metric_type"`
	Numerator                    types.String  `tfsdk:"numerator"`
	Datasource                   types.String  `tfsdk:"datasource"`
	Description                  types.String  `tfsdk:"description"`
	Owner                        types.String  `tfsdk:"owner"`
	Projects                     types.List    `tfsdk:"projects"`
	Tags                         types.List    `tfsdk:"tags"`
	Denominator                  types.String  `tfsdk:"denominator"`
	Inverse                      types.Bool    `tfsdk:"inverse"`
	CappingSettings              types.String  `tfsdk:"capping_settings"`
	WindowSettings               types.String  `tfsdk:"window_settings"`
	PriorSettings                types.String  `tfsdk:"prior_settings"`
	RegressionAdjustmentSettings types.String  `tfsdk:"regression_adjustment_settings"`
	RiskThresholdSuccess         types.Float64 `tfsdk:"risk_threshold_success"`
	RiskThresholdDanger          types.Float64 `tfsdk:"risk_threshold_danger"`
	MinPercentChange             types.Float64 `tfsdk:"min_percent_change"`
	MaxPercentChange             types.Float64 `tfsdk:"max_percent_change"`
	MinSampleSize                types.Float64 `tfsdk:"min_sample_size"`
	TargetMDE                    types.Float64 `tfsdk:"target_mde"`
	ManagedBy                    types.String  `tfsdk:"managed_by"`
	Archived                     types.Bool    `tfsdk:"archived"`
	DateCreated                  types.String  `tfsdk:"date_created"`
	DateUpdated                  types.String  `tfsdk:"date_updated"`
}

func (r *factMetricResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fact_metric"
}

func (r *factMetricResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                             schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"name":                           schema.StringAttribute{Required: true},
			"metric_type":                    schema.StringAttribute{Required: true},
			"numerator":                      schema.StringAttribute{Required: true, Validators: []validator.String{JSONString()}},
			"datasource":                     schema.StringAttribute{Optional: true, Computed: true},
			"description":                    schema.StringAttribute{Optional: true, Computed: true},
			"owner":                          schema.StringAttribute{Optional: true, Computed: true},
			"projects":                       schema.ListAttribute{ElementType: types.StringType, Optional: true, Computed: true},
			"tags":                           schema.ListAttribute{ElementType: types.StringType, Optional: true, Computed: true},
			"denominator":                    schema.StringAttribute{Optional: true, Computed: true, Validators: []validator.String{JSONString()}},
			"inverse":                        schema.BoolAttribute{Optional: true, Computed: true},
			"capping_settings":               schema.StringAttribute{Optional: true, Computed: true, Validators: []validator.String{JSONString()}},
			"window_settings":                schema.StringAttribute{Optional: true, Computed: true, Validators: []validator.String{JSONString()}},
			"prior_settings":                 schema.StringAttribute{Optional: true, Computed: true, Validators: []validator.String{JSONString()}},
			"regression_adjustment_settings": schema.StringAttribute{Optional: true, Computed: true, Validators: []validator.String{JSONString()}},
			"risk_threshold_success":         schema.Float64Attribute{Optional: true, Computed: true},
			"risk_threshold_danger":          schema.Float64Attribute{Optional: true, Computed: true},
			"min_percent_change":             schema.Float64Attribute{Optional: true, Computed: true},
			"max_percent_change":             schema.Float64Attribute{Optional: true, Computed: true},
			"min_sample_size":                schema.Float64Attribute{Optional: true, Computed: true},
			"target_mde":                     schema.Float64Attribute{Optional: true, Computed: true},
			"managed_by":                     schema.StringAttribute{Optional: true, Computed: true},
			"archived":                       schema.BoolAttribute{Optional: true, Computed: true},
			"date_created":                   schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"date_updated":                   schema.StringAttribute{Computed: true},
		},
	}
}

func (r *factMetricResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func factMetricFromModel(data factMetricModel) *growthbookapi.FactMetric {
	return &growthbookapi.FactMetric{
		Name:                         data.Name.ValueString(),
		MetricType:                   data.MetricType.ValueString(),
		Numerator:                    stringToRawJSON(data.Numerator.ValueString()),
		Datasource:                   data.Datasource.ValueString(),
		Description:                  data.Description.ValueString(),
		Owner:                        data.Owner.ValueString(),
		Projects:                     listToStrings(data.Projects),
		Tags:                         listToStrings(data.Tags),
		Denominator:                  stringToRawJSON(data.Denominator.ValueString()),
		Inverse:                      data.Inverse.ValueBool(),
		CappingSettings:              stringToRawJSON(data.CappingSettings.ValueString()),
		WindowSettings:               stringToRawJSON(data.WindowSettings.ValueString()),
		PriorSettings:                stringToRawJSON(data.PriorSettings.ValueString()),
		RegressionAdjustmentSettings: stringToRawJSON(data.RegressionAdjustmentSettings.ValueString()),
		RiskThresholdSuccess:         float64PointerValue(data.RiskThresholdSuccess),
		RiskThresholdDanger:          float64PointerValue(data.RiskThresholdDanger),
		MinPercentChange:             float64PointerValue(data.MinPercentChange),
		MaxPercentChange:             float64PointerValue(data.MaxPercentChange),
		MinSampleSize:                float64PointerValue(data.MinSampleSize),
		TargetMDE:                    float64PointerValue(data.TargetMDE),
		ManagedBy:                    data.ManagedBy.ValueString(),
		Archived:                     data.Archived.ValueBool(),
	}
}

func factMetricModelFromAPI(data *factMetricModel, item *growthbookapi.FactMetric) {
	data.ID = types.StringValue(item.ID)
	data.Name = types.StringValue(item.Name)
	data.MetricType = types.StringValue(item.MetricType)
	data.Numerator = types.StringValue(rawJSONToString(item.Numerator))
	data.Datasource = types.StringValue(item.Datasource)
	data.Description = types.StringValue(item.Description)
	data.Owner = types.StringValue(item.Owner)
	data.Projects = stringsToList(item.Projects)
	data.Tags = stringsToList(item.Tags)
	data.Denominator = types.StringValue(rawJSONToString(item.Denominator))
	data.Inverse = types.BoolValue(item.Inverse)
	data.CappingSettings = types.StringValue(rawJSONToString(item.CappingSettings))
	data.WindowSettings = types.StringValue(rawJSONToString(item.WindowSettings))
	data.PriorSettings = types.StringValue(rawJSONToString(item.PriorSettings))
	data.RegressionAdjustmentSettings = types.StringValue(rawJSONToString(item.RegressionAdjustmentSettings))
	data.RiskThresholdSuccess = float64ValuePointer(item.RiskThresholdSuccess)
	data.RiskThresholdDanger = float64ValuePointer(item.RiskThresholdDanger)
	data.MinPercentChange = float64ValuePointer(item.MinPercentChange)
	data.MaxPercentChange = float64ValuePointer(item.MaxPercentChange)
	data.MinSampleSize = float64ValuePointer(item.MinSampleSize)
	data.TargetMDE = float64ValuePointer(item.TargetMDE)
	data.ManagedBy = types.StringValue(item.ManagedBy)
	data.Archived = types.BoolValue(item.Archived)
	data.DateCreated = types.StringValue(item.DateCreated)
	data.DateUpdated = types.StringValue(item.DateUpdated)
}

func (r *factMetricResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data factMetricModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.CreateFactMetric(ctx, factMetricFromModel(data))
	if err != nil {
		resp.Diagnostics.AddError("Error creating fact metric", err.Error())
		return
	}

	factMetricModelFromAPI(&data, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *factMetricResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data factMetricModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.GetFactMetric(ctx, data.ID.ValueString())
	if err != nil {
		if errors.Is(err, growthbookapi.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading fact metric", err.Error())
		return
	}

	factMetricModelFromAPI(&data, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *factMetricResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data factMetricModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state factMetricModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.UpdateFactMetric(ctx, state.ID.ValueString(), factMetricFromModel(data))
	if err != nil {
		resp.Diagnostics.AddError("Error updating fact metric", err.Error())
		return
	}

	factMetricModelFromAPI(&data, item)
	data.ID = state.ID
	if data.DateCreated.ValueString() == "" {
		data.DateCreated = state.DateCreated
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *factMetricResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data factMetricModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteFactMetric(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting fact metric", err.Error())
	}
}

func (r *factMetricResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

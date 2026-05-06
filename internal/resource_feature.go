package internal

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-growthbook/internal/growthbookapi"
)

var _ resource.Resource = &featureResource{}
var _ resource.ResourceWithImportState = &featureResource{}

func newFeatureResource() resource.Resource {
	return &featureResource{}
}

type featureResource struct {
	client *growthbookapi.Client
}

// featureEnvironmentModel maps a single GrowthBook feature environment.
type featureEnvironmentModel struct {
	Enabled      types.Bool         `tfsdk:"enabled"`
	DefaultValue types.String       `tfsdk:"default_value"`
	Rules        []featureRuleModel `tfsdk:"rules"`
}

// featureRuleModel maps a single targeting rule (force / rollout / experiment-ref).
type featureRuleModel struct {
	ID            types.String            `tfsdk:"id"`
	Type          types.String            `tfsdk:"type"`
	Enabled       types.Bool              `tfsdk:"enabled"`
	Description   types.String            `tfsdk:"description"`
	Condition     types.String            `tfsdk:"condition"`
	Value         types.String            `tfsdk:"value"`
	Coverage      types.Float64           `tfsdk:"coverage"`
	HashAttribute types.String            `tfsdk:"hash_attribute"`
	ExperimentID  types.String            `tfsdk:"experiment_id"`
	Variations    []featureVariationModel `tfsdk:"variations"`
	Prerequisites []featurePrereqModel    `tfsdk:"prerequisites"`
}

// featureVariationModel maps a single experiment-ref variation.
type featureVariationModel struct {
	Value       types.String `tfsdk:"value"`
	VariationID types.String `tfsdk:"variation_id"`
}

// featurePrereqModel maps a rule prerequisite.
type featurePrereqModel struct {
	ID        types.String `tfsdk:"id"`
	Condition types.String `tfsdk:"condition"`
}

type featureModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Archived      types.Bool   `tfsdk:"archived"`
	Description   types.String `tfsdk:"description"`
	Owner         types.String `tfsdk:"owner"`
	Project       types.String `tfsdk:"project"`
	ValueType     types.String `tfsdk:"value_type"`
	DefaultValue  types.String `tfsdk:"default_value"`
	Tags          types.List   `tfsdk:"tags"`
	Environments  types.Map    `tfsdk:"environments"`
	Prerequisites types.List   `tfsdk:"prerequisites"`
}

// Attribute type definitions for nested objects (used by types.MapValueFrom).

func featurePrereqObjectType() types.ObjectType {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":        types.StringType,
		"condition": types.StringType,
	}}
}

func featureVariationObjectType() types.ObjectType {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"value":        types.StringType,
		"variation_id": types.StringType,
	}}
}

func featureRuleObjectType() types.ObjectType {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":             types.StringType,
		"type":           types.StringType,
		"enabled":        types.BoolType,
		"description":    types.StringType,
		"condition":      types.StringType,
		"value":          types.StringType,
		"coverage":       types.Float64Type,
		"hash_attribute": types.StringType,
		"experiment_id":  types.StringType,
		"variations":     types.ListType{ElemType: featureVariationObjectType()},
		"prerequisites":  types.ListType{ElemType: featurePrereqObjectType()},
	}}
}

func featureEnvObjectType() types.ObjectType {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"enabled":       types.BoolType,
		"default_value": types.StringType,
		"rules":         types.ListType{ElemType: featureRuleObjectType()},
	}}
}

func featureEnvironmentSchemaAttr() schema.MapNestedAttribute {
	return schema.MapNestedAttribute{
		Optional: true,
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"enabled": schema.BoolAttribute{
					Required: true,
				},
				"default_value": schema.StringAttribute{
					Optional: true,
					Computed: true,
				},
				"rules": schema.ListNestedAttribute{
					Optional: true,
					Computed: true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Optional: true,
								Computed: true,
							},
							"type": schema.StringAttribute{
								Optional: true,
								Computed: true,
							},
							"enabled": schema.BoolAttribute{
								Optional: true,
								Computed: true,
							},
							"description": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString(""),
							},
							"condition": schema.StringAttribute{
								Optional: true,
								Computed: true,
							},
							"value": schema.StringAttribute{
								Optional: true,
								Computed: true,
							},
							"coverage": schema.Float64Attribute{
								Optional: true,
								Computed: true,
								Default:  float64default.StaticFloat64(1.0),
							},
							"hash_attribute": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString(""),
							},
							"experiment_id": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Default:  stringdefault.StaticString(""),
							},
							"variations": schema.ListNestedAttribute{
								Optional: true,
								Computed: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"value": schema.StringAttribute{
											Required: true,
										},
										"variation_id": schema.StringAttribute{
											Required: true,
										},
									},
								},
								Default: listdefault.StaticValue(types.ListValueMust(types.ObjectType{AttrTypes: map[string]attr.Type{
									"value":        types.StringType,
									"variation_id": types.StringType,
								}}, []attr.Value{})),
							},
							"prerequisites": schema.ListNestedAttribute{
								Optional: true,
								Computed: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Required: true,
										},
										"condition": schema.StringAttribute{
											Required: true,
										},
									},
								},
								Default: listdefault.StaticValue(types.ListValueMust(types.ObjectType{AttrTypes: map[string]attr.Type{
									"id":        types.StringType,
									"condition": types.StringType,
								}}, []attr.Value{})),
							},
						},
					},
				},
			},
		},
	}
}

func (r *featureResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_feature"
}

func (r *featureResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"archived": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"owner": schema.StringAttribute{
				Required: true,
			},
			"project": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"value_type": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"default_value": schema.StringAttribute{
				Required: true,
			},
			"tags": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"environments": featureEnvironmentSchemaAttr(),
			"prerequisites": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *featureResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *featureResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data featureModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tags := []string{}
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	prereqs := []string{}
	if !data.Prerequisites.IsNull() && !data.Prerequisites.IsUnknown() {
		resp.Diagnostics.Append(data.Prerequisites.ElementsAs(ctx, &prereqs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var apiEnvs map[string]growthbookapi.FeatureEnvironmentConfig
	if !data.Environments.IsNull() && !data.Environments.IsUnknown() {
		var envModels map[string]featureEnvironmentModel
		resp.Diagnostics.Append(data.Environments.ElementsAs(ctx, &envModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiEnvs = envsToAPI(envModels)
	}

	feature := &growthbookapi.Feature{
		ID:            data.Name.ValueString(),
		Description:   data.Description.ValueString(),
		Owner:         data.Owner.ValueString(),
		Project:       data.Project.ValueString(),
		ValueType:     data.ValueType.ValueString(),
		DefaultValue:  data.DefaultValue.ValueString(),
		Tags:          tags,
		Prerequisites: prereqs,
		Environments:  apiEnvs,
	}

	created, err := r.client.CreateFeature(ctx, feature)
	if err != nil {
		resp.Diagnostics.AddError("Error creating feature", err.Error())
		return
	}

	resp.Diagnostics.Append(featureModelFromAPI(ctx, &data, created)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *featureResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data featureModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	feature, err := r.client.GetFeature(ctx, data.ID.ValueString())
	if err != nil {
		if errors.Is(err, growthbookapi.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading feature", err.Error())
		return
	}

	resp.Diagnostics.Append(featureModelFromAPI(ctx, &data, feature)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *featureResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data featureModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state featureModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tags := []string{}
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	prereqs := []string{}
	if !data.Prerequisites.IsNull() && !data.Prerequisites.IsUnknown() {
		resp.Diagnostics.Append(data.Prerequisites.ElementsAs(ctx, &prereqs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var apiEnvs map[string]growthbookapi.FeatureEnvironmentConfig
	if !data.Environments.IsNull() && !data.Environments.IsUnknown() {
		var envModels map[string]featureEnvironmentModel
		resp.Diagnostics.Append(data.Environments.ElementsAs(ctx, &envModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiEnvs = envsToAPI(envModels)
	}

	feature := &growthbookapi.Feature{
		Archived:      data.Archived.ValueBool(),
		Description:   data.Description.ValueString(),
		Owner:         data.Owner.ValueString(),
		Project:       data.Project.ValueString(),
		DefaultValue:  data.DefaultValue.ValueString(),
		Tags:          tags,
		Prerequisites: prereqs,
		Environments:  apiEnvs,
	}

	updated, err := r.client.UpdateFeature(ctx, state.ID.ValueString(), feature)
	if err != nil {
		resp.Diagnostics.AddError("Error updating feature", err.Error())
		return
	}

	resp.Diagnostics.Append(featureModelFromAPI(ctx, &data, updated)...)
	data.ID = state.ID // preserve original ID in case API returns a different casing
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *featureResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data featureModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteFeature(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting feature", err.Error())
	}
}

func (r *featureResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// featureModelFromAPI populates a featureModel from a GrowthBook API Feature.
func featureModelFromAPI(ctx context.Context, m *featureModel, f *growthbookapi.Feature) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(f.ID)
	m.Name = types.StringValue(f.ID)
	m.Archived = types.BoolValue(f.Archived)
	m.Description = types.StringValue(f.Description)
	m.Owner = types.StringValue(f.Owner)
	m.Project = types.StringValue(f.Project)
	m.ValueType = types.StringValue(f.ValueType)
	m.DefaultValue = types.StringValue(f.DefaultValue)
	m.Tags = stringsToList(f.Tags)
	m.Prerequisites = stringsToList(f.Prerequisites)

	envsMap := envsFromAPI(f.Environments)
	var d diag.Diagnostics
	m.Environments, d = types.MapValueFrom(ctx, featureEnvObjectType(), envsMap)
	diags.Append(d...)

	return diags
}

// envsFromAPI converts API environments to Terraform model environments.
func envsFromAPI(envs map[string]growthbookapi.FeatureEnvironmentConfig) map[string]featureEnvironmentModel {
	if envs == nil {
		return nil
	}
	result := make(map[string]featureEnvironmentModel, len(envs))
	for name, env := range envs {
		result[name] = featureEnvironmentModel{
			Enabled:      types.BoolValue(env.Enabled),
			DefaultValue: types.StringValue(env.DefaultValue),
			Rules:        rulesFromAPI(env.Rules),
		}
	}
	return result
}

func rulesFromAPI(rules []growthbookapi.FeatureRule) []featureRuleModel {
	out := make([]featureRuleModel, len(rules))
	for i, r := range rules {
		rm := featureRuleModel{
			ID:            types.StringValue(r.ID),
			Type:          types.StringValue(r.Type),
			Enabled:       types.BoolValue(r.Enabled),
			Description:   types.StringValue(r.Description),
			Condition:     types.StringValue(r.Condition),
			Value:         types.StringValue(r.Value),
			HashAttribute: types.StringValue(r.HashAttribute),
			ExperimentID:  types.StringValue(r.ExperimentID),
			Variations:    variationsFromAPI(r.Variations),
			Prerequisites: rulePrereqsFromAPI(r.Prerequisites),
		}
		if r.Coverage != nil {
			rm.Coverage = types.Float64Value(*r.Coverage)
		} else {
			rm.Coverage = types.Float64Null()
		}
		out[i] = rm
	}
	return out
}

func variationsFromAPI(vars []growthbookapi.FeatureVariation) []featureVariationModel {
	out := make([]featureVariationModel, len(vars))
	for i, v := range vars {
		out[i] = featureVariationModel{
			Value:       types.StringValue(v.Value),
			VariationID: types.StringValue(v.VariationID),
		}
	}
	return out
}

func rulePrereqsFromAPI(prereqs []growthbookapi.FeaturePrerequisite) []featurePrereqModel {
	out := make([]featurePrereqModel, len(prereqs))
	for i, p := range prereqs {
		out[i] = featurePrereqModel{
			ID:        types.StringValue(p.ID),
			Condition: types.StringValue(p.Condition),
		}
	}
	return out
}

// envsToAPI converts Terraform model environments to API environments.
func envsToAPI(envs map[string]featureEnvironmentModel) map[string]growthbookapi.FeatureEnvironmentConfig {
	if envs == nil {
		return nil
	}
	result := make(map[string]growthbookapi.FeatureEnvironmentConfig, len(envs))
	for name, env := range envs {
		result[name] = growthbookapi.FeatureEnvironmentConfig{
			Enabled:      env.Enabled.ValueBool(),
			DefaultValue: env.DefaultValue.ValueString(),
			Rules:        rulesToAPI(env.Rules),
		}
	}
	return result
}

func rulesToAPI(rules []featureRuleModel) []growthbookapi.FeatureRule {
	out := make([]growthbookapi.FeatureRule, len(rules))
	for i, r := range rules {
		ar := growthbookapi.FeatureRule{
			ID:            r.ID.ValueString(),
			Type:          r.Type.ValueString(),
			Enabled:       r.Enabled.ValueBool(),
			Description:   r.Description.ValueString(),
			Condition:     r.Condition.ValueString(),
			Value:         r.Value.ValueString(),
			HashAttribute: r.HashAttribute.ValueString(),
			ExperimentID:  r.ExperimentID.ValueString(),
			Variations:    variationsToAPI(r.Variations),
			Prerequisites: rulePrereqsToAPI(r.Prerequisites),
		}
		if !r.Coverage.IsNull() && !r.Coverage.IsUnknown() {
			v := r.Coverage.ValueFloat64()
			ar.Coverage = &v
		}
		out[i] = ar
	}
	return out
}

func variationsToAPI(vars []featureVariationModel) []growthbookapi.FeatureVariation {
	out := make([]growthbookapi.FeatureVariation, len(vars))
	for i, v := range vars {
		out[i] = growthbookapi.FeatureVariation{
			Value:       v.Value.ValueString(),
			VariationID: v.VariationID.ValueString(),
		}
	}
	return out
}

func rulePrereqsToAPI(prereqs []featurePrereqModel) []growthbookapi.FeaturePrerequisite {
	out := make([]growthbookapi.FeaturePrerequisite, len(prereqs))
	for i, p := range prereqs {
		out[i] = growthbookapi.FeaturePrerequisite{
			ID:        p.ID.ValueString(),
			Condition: p.Condition.ValueString(),
		}
	}
	return out
}

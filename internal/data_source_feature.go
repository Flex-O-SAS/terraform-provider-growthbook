package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-growthbook/internal/growthbookapi"
)

var _ datasource.DataSource = &featureDataSource{}

func newFeatureDataSource() datasource.DataSource {
	return &featureDataSource{}
}

type featureDataSource struct {
	client *growthbookapi.Client
}

type featureDataModel struct {
	ID            types.String `tfsdk:"id"`
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

func featureDataEnvironmentSchemaAttr() schema.MapNestedAttribute {
	return schema.MapNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"enabled": schema.BoolAttribute{
					Computed: true,
				},
				"default_value": schema.StringAttribute{
					Computed: true,
				},
				"rules": schema.ListNestedAttribute{
					Computed: true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"id":          schema.StringAttribute{Computed: true},
							"type":        schema.StringAttribute{Computed: true},
							"enabled":     schema.BoolAttribute{Computed: true},
							"description": schema.StringAttribute{Computed: true},
							"condition":   schema.StringAttribute{Computed: true},
							"value":       schema.StringAttribute{Computed: true},
							"coverage":    schema.Float64Attribute{Computed: true},
							"hash_attribute": schema.StringAttribute{Computed: true},
							"experiment_id":  schema.StringAttribute{Computed: true},
							"variations": schema.ListNestedAttribute{
								Computed: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"value":        schema.StringAttribute{Computed: true},
										"variation_id": schema.StringAttribute{Computed: true},
									},
								},
							},
							"prerequisites": schema.ListNestedAttribute{
								Computed: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"id":        schema.StringAttribute{Computed: true},
										"condition": schema.StringAttribute{Computed: true},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *featureDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_feature"
}

func (d *featureDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the GrowthBook feature.",
			},
			"archived": schema.BoolAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"owner": schema.StringAttribute{
				Computed: true,
			},
			"project": schema.StringAttribute{
				Computed: true,
			},
			"value_type": schema.StringAttribute{
				Computed: true,
			},
			"default_value": schema.StringAttribute{
				Computed: true,
			},
			"tags": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"environments": featureDataEnvironmentSchemaAttr(),
			"prerequisites": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

func (d *featureDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*growthbookapi.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data type", "Expected *growthbookapi.Client")
		return
	}
	d.client = client
}

func (d *featureDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data featureDataModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	feature, err := d.client.GetFeature(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to find GrowthBook feature by ID", err.Error())
		return
	}

	data.ID = types.StringValue(feature.ID)
	data.Archived = types.BoolValue(feature.Archived)
	data.Description = types.StringValue(feature.Description)
	data.Owner = types.StringValue(feature.Owner)
	data.Project = types.StringValue(feature.Project)
	data.ValueType = types.StringValue(feature.ValueType)
	data.DefaultValue = types.StringValue(feature.DefaultValue)
	data.Tags = stringsToList(ctx, feature.Tags)
	data.Prerequisites = stringsToList(ctx, feature.Prerequisites)

	envsMap := envsFromAPI(feature.Environments)
	var envDiags diag.Diagnostics
	data.Environments, envDiags = types.MapValueFrom(ctx, featureEnvObjectType(), envsMap)
	resp.Diagnostics.Append(envDiags...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

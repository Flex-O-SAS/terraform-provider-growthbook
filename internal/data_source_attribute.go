package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-growthbook/internal/growthbookapi"
)

var _ datasource.DataSource = &attributeDataSource{}

func newAttributeDataSource() datasource.DataSource {
	return &attributeDataSource{}
}

type attributeDataSource struct {
	client *growthbookapi.Client
}

type attributeDataModel struct {
	Property    types.String `tfsdk:"property"`
	DataType    types.String `tfsdk:"datatype"`
	Format      types.String `tfsdk:"format"`
	EnumValues  types.String `tfsdk:"enum_values"`
	Projects    types.List   `tfsdk:"projects"`
	Archived    types.Bool   `tfsdk:"archived"`
	Description types.String `tfsdk:"description"`
}

func (d *attributeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_attribute"
}

func (d *attributeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"property": schema.StringAttribute{
				Required: true,
			},
			"datatype": schema.StringAttribute{
				Computed: true,
			},
			"format": schema.StringAttribute{
				Computed: true,
			},
			"enum_values": schema.StringAttribute{
				Computed: true,
			},
			"projects": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Array of project IDs.",
			},
			"archived": schema.BoolAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *attributeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *attributeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data attributeDataModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	attribute, err := d.client.GetAttribute(ctx, data.Property.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to find GrowthBook attribute by property", err.Error())
		return
	}

	data.Property = types.StringValue(attribute.Property)
	data.DataType = types.StringValue(attribute.DataType)
	data.Format = types.StringValue(attribute.Format)
	data.EnumValues = types.StringValue(attribute.EnumValues)
	data.Projects = stringsToList(attribute.Projects)
	data.Archived = types.BoolValue(attribute.Archived)
	data.Description = types.StringValue(attribute.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

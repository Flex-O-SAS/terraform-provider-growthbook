package internal

import (
	"context"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-growthbook/internal/growthbookapi"
)

var _ resource.Resource = &attributeResource{}
var _ resource.ResourceWithImportState = &attributeResource{}

func newAttributeResource() resource.Resource {
	return &attributeResource{}
}

type attributeResource struct {
	client *growthbookapi.Client
}

type attributeModel struct {
	Property    types.String `tfsdk:"property"`
	DataType    types.String `tfsdk:"datatype"`
	Format      types.String `tfsdk:"format"`
	EnumValues  types.String `tfsdk:"enum_values"`
	Projects    types.List   `tfsdk:"projects"`
	Archived    types.Bool   `tfsdk:"archived"`
	Description types.String `tfsdk:"description"`
}

func (r *attributeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_attribute"
}

func (r *attributeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"property": schema.StringAttribute{
				Required: true,
			},
			"datatype": schema.StringAttribute{
				Required: true,
			},
			"format": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"enum_values": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"projects": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "Array of project IDs.",
			},
			"archived": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
		},
	}
}

func (r *attributeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

var validFormats = []string{"", "version", "date", "isoCountryCode"}

func (r *attributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data attributeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	format := data.Format.ValueString()
	if !slices.Contains(validFormats, format) {
		resp.Diagnostics.AddError(
			"Invalid format value",
			"Expected '' | 'version' | 'date' | 'isoCountryCode', received: "+format,
		)
		return
	}

	projects := []string{}
	if !data.Projects.IsNull() && !data.Projects.IsUnknown() {
		resp.Diagnostics.Append(data.Projects.ElementsAs(ctx, &projects, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	attribute := &growthbookapi.Attribute{
		Property:    data.Property.ValueString(),
		DataType:    data.DataType.ValueString(),
		Format:      format,
		EnumValues:  data.EnumValues.ValueString(),
		Projects:    projects,
		Archived:    data.Archived.ValueBool(),
		Description: data.Description.ValueString(),
	}

	created, err := r.client.CreateAttribute(ctx, attribute)
	if err != nil {
		resp.Diagnostics.AddError("Error creating attribute", err.Error())
		return
	}

	data.Property = types.StringValue(created.Property)
	data.DataType = types.StringValue(created.DataType)
	data.Format = types.StringValue(created.Format)
	data.EnumValues = types.StringValue(created.EnumValues)
	data.Projects = stringsToList(ctx, created.Projects)
	data.Archived = types.BoolValue(created.Archived)
	data.Description = types.StringValue(created.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *attributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data attributeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, err := r.client.GetAttribute(ctx, data.Property.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading attribute", err.Error())
		return
	}
	if out == nil {
		resp.Diagnostics.AddError("Attribute not found", "Attribute with property '"+data.Property.ValueString()+"' was not found.")
		return
	}

	data.Property = types.StringValue(out.Property)
	data.DataType = types.StringValue(out.DataType)
	data.Format = types.StringValue(out.Format)
	data.EnumValues = types.StringValue(out.EnumValues)
	data.Projects = stringsToList(ctx, out.Projects)
	data.Archived = types.BoolValue(out.Archived)
	data.Description = types.StringValue(out.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *attributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data attributeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projects := []string{}
	if !data.Projects.IsNull() && !data.Projects.IsUnknown() {
		resp.Diagnostics.Append(data.Projects.ElementsAs(ctx, &projects, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	attribute := &growthbookapi.Attribute{
		Property:    data.Property.ValueString(),
		DataType:    data.DataType.ValueString(),
		Format:      data.Format.ValueString(),
		EnumValues:  data.EnumValues.ValueString(),
		Projects:    projects,
		Archived:    data.Archived.ValueBool(),
		Description: data.Description.ValueString(),
	}

	_, err := r.client.UpdateAttribute(ctx, attribute.Property, attribute)
	if err != nil {
		resp.Diagnostics.AddError("Error updating attribute", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *attributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data attributeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteAttribute(ctx, data.Property.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting attribute", err.Error())
	}
}

func (r *attributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("property"), req, resp)
}

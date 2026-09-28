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

var _ resource.Resource = &namespaceResource{}
var _ resource.ResourceWithImportState = &namespaceResource{}

func newNamespaceResource() resource.Resource {
	return &namespaceResource{}
}

type namespaceResource struct {
	client *growthbookapi.Client
}

type namespaceModel struct {
	ID            types.String `tfsdk:"id"`
	DisplayName   types.String `tfsdk:"display_name"`
	Description   types.String `tfsdk:"description"`
	Status        types.String `tfsdk:"status"`
	Format        types.String `tfsdk:"format"`
	HashAttribute types.String `tfsdk:"hash_attribute"`
	Seed          types.String `tfsdk:"seed"`
}

func (r *namespaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace"
}

func (r *namespaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":             schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"display_name":   schema.StringAttribute{Required: true},
			"description":    schema.StringAttribute{Optional: true, Computed: true},
			"status":         schema.StringAttribute{Optional: true, Computed: true},
			"format":         schema.StringAttribute{Optional: true, Computed: true},
			"hash_attribute": schema.StringAttribute{Optional: true, Computed: true},
			"seed":           schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
		},
	}
}

func (r *namespaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func namespaceFromModel(data namespaceModel) *growthbookapi.Namespace {
	return &growthbookapi.Namespace{
		DisplayName:   data.DisplayName.ValueString(),
		Description:   data.Description.ValueString(),
		Status:        data.Status.ValueString(),
		Format:        data.Format.ValueString(),
		HashAttribute: data.HashAttribute.ValueString(),
	}
}

func namespaceModelFromAPI(data *namespaceModel, item *growthbookapi.Namespace) {
	data.ID = types.StringValue(item.ID)
	data.DisplayName = types.StringValue(item.DisplayName)
	data.Description = types.StringValue(item.Description)
	data.Status = types.StringValue(item.Status)
	data.Format = types.StringValue(item.Format)
	data.HashAttribute = types.StringValue(item.HashAttribute)
	data.Seed = types.StringValue(item.Seed)
}

func (r *namespaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data namespaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.CreateNamespace(ctx, namespaceFromModel(data))
	if err != nil {
		resp.Diagnostics.AddError("Error creating namespace", err.Error())
		return
	}

	namespaceModelFromAPI(&data, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *namespaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data namespaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.GetNamespace(ctx, data.ID.ValueString())
	if err != nil {
		if errors.Is(err, growthbookapi.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading namespace", err.Error())
		return
	}

	namespaceModelFromAPI(&data, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *namespaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data namespaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state namespaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.UpdateNamespace(ctx, state.ID.ValueString(), namespaceFromModel(data))
	if err != nil {
		resp.Diagnostics.AddError("Error updating namespace", err.Error())
		return
	}

	namespaceModelFromAPI(&data, item)
	data.ID = state.ID
	if data.Seed.ValueString() == "" {
		data.Seed = state.Seed
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *namespaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data namespaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteNamespace(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting namespace", err.Error())
	}
}

func (r *namespaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

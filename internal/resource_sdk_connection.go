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

var _ resource.Resource = &sdkConnectionResource{}
var _ resource.ResourceWithImportState = &sdkConnectionResource{}

func newSDKConnectionResource() resource.Resource {
	return &sdkConnectionResource{}
}

type sdkConnectionResource struct {
	client *growthbookapi.Client
}

type sdkConnectionModel struct {
	ID                          types.String `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	Environment                 types.String `tfsdk:"environment"`
	Language                    types.String `tfsdk:"language"`
	SdkVersion                  types.String `tfsdk:"sdk_version"`
	Projects                    types.List   `tfsdk:"projects"`
	EncryptPayload              types.Bool   `tfsdk:"encrypt_payload"`
	IncludeVisualExperiments    types.Bool   `tfsdk:"include_visual_experiments"`
	IncludeDraftExperiments     types.Bool   `tfsdk:"include_draft_experiments"`
	IncludeExperimentNames      types.Bool   `tfsdk:"include_experiment_names"`
	IncludeRedirectExperiments  types.Bool   `tfsdk:"include_redirect_experiments"`
	IncludeRuleIDs              types.Bool   `tfsdk:"include_rule_ids"`
	ProxyEnabled                types.Bool   `tfsdk:"proxy_enabled"`
	ProxyHost                   types.String `tfsdk:"proxy_host"`
	HashSecureAttributes        types.Bool   `tfsdk:"hash_secure_attributes"`
	RemoteEvalEnabled           types.Bool   `tfsdk:"remote_eval_enabled"`
	SavedGroupReferencesEnabled types.Bool   `tfsdk:"saved_group_references_enabled"`
	// computed
	Organization    types.String `tfsdk:"organization"`
	Key             types.String `tfsdk:"key"`
	ProxySigningKey types.String `tfsdk:"proxy_signing_key"`
	SseEnabled      types.Bool   `tfsdk:"sse_enabled"`
	EncryptionKey   types.String `tfsdk:"encryption_key"`
	DateCreated     types.String `tfsdk:"date_created"`
	DateUpdated     types.String `tfsdk:"date_updated"`
}

func (r *sdkConnectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdk_connection"
}

func (r *sdkConnectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"environment": schema.StringAttribute{
				Required: true,
			},
			"language": schema.StringAttribute{
				Required: true,
			},
			"sdk_version": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"projects": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"encrypt_payload": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"include_visual_experiments": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"include_draft_experiments": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"include_experiment_names": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"include_redirect_experiments": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"include_rule_ids": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"proxy_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"proxy_host": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"hash_secure_attributes": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"remote_eval_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"saved_group_references_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			// computed only
			"organization": schema.StringAttribute{
				Computed: true,
			},
			"key": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"proxy_signing_key": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"sse_enabled": schema.BoolAttribute{
				Computed: true,
			},
			"encryption_key": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"date_created": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"date_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *sdkConnectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func sdkConnToModel(ctx context.Context, conn *growthbookapi.SDKConnection) sdkConnectionModel {
	return sdkConnectionModel{
		ID:                          types.StringValue(conn.ID),
		Name:                        types.StringValue(conn.Name),
		Language:                    types.StringValue(conn.Language),
		Environment:                 types.StringValue(conn.Environment),
		SdkVersion:                  types.StringValue(conn.SdkVersion),
		Projects:                    stringsToList(conn.Projects),
		EncryptPayload:              types.BoolValue(conn.EncryptPayload),
		IncludeVisualExperiments:    types.BoolValue(conn.IncludeVisualExperiments),
		IncludeDraftExperiments:     types.BoolValue(conn.IncludeDraftExperiments),
		IncludeExperimentNames:      types.BoolValue(conn.IncludeExperimentNames),
		IncludeRedirectExperiments:  types.BoolValue(conn.IncludeRedirectExperiments),
		IncludeRuleIDs:              types.BoolValue(conn.IncludeRuleIDs),
		ProxyEnabled:                types.BoolValue(conn.ProxyEnabled),
		ProxyHost:                   types.StringValue(conn.ProxyHost),
		HashSecureAttributes:        types.BoolValue(conn.HashSecureAttributes),
		RemoteEvalEnabled:           types.BoolValue(conn.RemoteEvalEnabled),
		SavedGroupReferencesEnabled: types.BoolValue(conn.SavedGroupReferencesEnabled),
		Organization:                types.StringValue(conn.Organization),
		Key:                         types.StringValue(conn.Key),
		ProxySigningKey:             types.StringValue(conn.ProxySigningKey),
		SseEnabled:                  types.BoolValue(conn.SseEnabled),
		EncryptionKey:               types.StringValue(conn.EncryptionKey),
		DateCreated:                 types.StringValue(conn.DateCreated),
		DateUpdated:                 types.StringValue(conn.DateUpdated),
	}
}

func sdkConnFromPlan(ctx context.Context, data sdkConnectionModel, projects []string) *growthbookapi.SDKConnection {
	return &growthbookapi.SDKConnection{
		Name:                        data.Name.ValueString(),
		Language:                    data.Language.ValueString(),
		Environment:                 data.Environment.ValueString(),
		SdkVersion:                  data.SdkVersion.ValueString(),
		Projects:                    projects,
		EncryptPayload:              data.EncryptPayload.ValueBool(),
		IncludeVisualExperiments:    data.IncludeVisualExperiments.ValueBool(),
		IncludeDraftExperiments:     data.IncludeDraftExperiments.ValueBool(),
		IncludeExperimentNames:      data.IncludeExperimentNames.ValueBool(),
		IncludeRedirectExperiments:  data.IncludeRedirectExperiments.ValueBool(),
		IncludeRuleIDs:              data.IncludeRuleIDs.ValueBool(),
		ProxyEnabled:                data.ProxyEnabled.ValueBool(),
		ProxyHost:                   data.ProxyHost.ValueString(),
		HashSecureAttributes:        data.HashSecureAttributes.ValueBool(),
		RemoteEvalEnabled:           data.RemoteEvalEnabled.ValueBool(),
		SavedGroupReferencesEnabled: data.SavedGroupReferencesEnabled.ValueBool(),
	}
}

func (r *sdkConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data sdkConnectionModel
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

	created, err := r.client.CreateSDKConnection(ctx, sdkConnFromPlan(ctx, data, projects))
	if err != nil {
		resp.Diagnostics.AddError("Error creating SDK connection", err.Error())
		return
	}

	result := sdkConnToModel(ctx, created)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r *sdkConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data sdkConnectionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	conn, err := r.client.GetSDKConnection(ctx, data.ID.ValueString())
	if err != nil {
		if errors.Is(err, growthbookapi.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading SDK connection", err.Error())
		return
	}

	result := sdkConnToModel(ctx, conn)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r *sdkConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data sdkConnectionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state sdkConnectionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
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

	updated, err := r.client.UpdateSDKConnection(ctx, state.ID.ValueString(), sdkConnFromPlan(ctx, data, projects))
	if err != nil {
		resp.Diagnostics.AddError("Error updating SDK connection", err.Error())
		return
	}

	result := sdkConnToModel(ctx, updated)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r *sdkConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data sdkConnectionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSDKConnection(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting SDK connection", err.Error())
	}
}

func (r *sdkConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

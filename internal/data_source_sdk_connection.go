package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-growthbook/internal/growthbookapi"
)

var _ datasource.DataSource = &sdkConnectionDataSource{}

func newSDKConnectionDataSource() datasource.DataSource {
	return &sdkConnectionDataSource{}
}

type sdkConnectionDataSource struct {
	client *growthbookapi.Client
}

type sdkConnectionDataModel struct {
	ID                          types.String `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	Organization                types.String `tfsdk:"organization"`
	Language                    types.String `tfsdk:"language"`
	SdkVersion                  types.String `tfsdk:"sdk_version"`
	Environment                 types.String `tfsdk:"environment"`
	Projects                    types.List   `tfsdk:"projects"`
	EncryptPayload              types.Bool   `tfsdk:"encrypt_payload"`
	EncryptionKey               types.String `tfsdk:"encryption_key"`
	IncludeVisualExperiments    types.Bool   `tfsdk:"include_visual_experiments"`
	IncludeDraftExperiments     types.Bool   `tfsdk:"include_draft_experiments"`
	IncludeExperimentNames      types.Bool   `tfsdk:"include_experiment_names"`
	IncludeRedirectExperiments  types.Bool   `tfsdk:"include_redirect_experiments"`
	IncludeRuleIDs              types.Bool   `tfsdk:"include_rule_ids"`
	Key                         types.String `tfsdk:"key"`
	ProxyEnabled                types.Bool   `tfsdk:"proxy_enabled"`
	ProxyHost                   types.String `tfsdk:"proxy_host"`
	ProxySigningKey             types.String `tfsdk:"proxy_signing_key"`
	SseEnabled                  types.Bool   `tfsdk:"sse_enabled"`
	HashSecureAttributes        types.Bool   `tfsdk:"hash_secure_attributes"`
	RemoteEvalEnabled           types.Bool   `tfsdk:"remote_eval_enabled"`
	SavedGroupReferencesEnabled types.Bool   `tfsdk:"saved_group_references_enabled"`
	DateCreated                 types.String `tfsdk:"date_created"`
	DateUpdated                 types.String `tfsdk:"date_updated"`
}

func (d *sdkConnectionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdk_connection"
}

func (d *sdkConnectionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the SDK Connection.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier for the SDK Connection.",
			},
			"organization": schema.StringAttribute{
				Computed:    true,
				Description: "The organization associated with the SDK Connection.",
			},
			"language": schema.StringAttribute{
				Computed:    true,
				Description: "The programming language for the SDK.",
			},
			"sdk_version": schema.StringAttribute{
				Computed:    true,
				Description: "The version of the SDK.",
			},
			"environment": schema.StringAttribute{
				Computed:    true,
				Description: "The environment for the SDK Connection.",
			},
			"projects": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "The projects associated with the SDK Connection.",
			},
			"encrypt_payload": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether to encrypt the payload.",
			},
			"encryption_key": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The key used for encryption.",
			},
			"include_visual_experiments": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether to include visual experiments.",
			},
			"include_draft_experiments": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether to include draft experiments.",
			},
			"include_experiment_names": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether to include experiment names.",
			},
			"include_redirect_experiments": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether to include redirect experiments.",
			},
			"include_rule_ids": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether to include rule IDs.",
			},
			"key": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The key for the SDK Connection.",
			},
			"proxy_enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the proxy is enabled.",
			},
			"proxy_host": schema.StringAttribute{
				Computed:    true,
				Description: "The host of the proxy.",
			},
			"proxy_signing_key": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The signing key for the proxy.",
			},
			"sse_enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether server-sent events are enabled.",
			},
			"hash_secure_attributes": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether to hash secure attributes.",
			},
			"remote_eval_enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether remote evaluation is enabled.",
			},
			"saved_group_references_enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether saved group references are enabled.",
			},
			"date_created": schema.StringAttribute{
				Computed: true,
			},
			"date_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *sdkConnectionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *sdkConnectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data sdkConnectionDataModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	conn, err := d.client.FindSDKConnectionByName(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading SDK connection", err.Error())
		return
	}

	data.ID = types.StringValue(conn.ID)
	data.Name = types.StringValue(conn.Name)
	data.Organization = types.StringValue(conn.Organization)
	data.Language = types.StringValue(conn.Language)
	data.SdkVersion = types.StringValue(conn.SdkVersion)
	data.Environment = types.StringValue(conn.Environment)
	data.Projects = stringsToList(conn.Projects)
	data.EncryptPayload = types.BoolValue(conn.EncryptPayload)
	data.EncryptionKey = types.StringValue(conn.EncryptionKey)
	data.IncludeVisualExperiments = types.BoolValue(conn.IncludeVisualExperiments)
	data.IncludeDraftExperiments = types.BoolValue(conn.IncludeDraftExperiments)
	data.IncludeExperimentNames = types.BoolValue(conn.IncludeExperimentNames)
	data.IncludeRedirectExperiments = types.BoolValue(conn.IncludeRedirectExperiments)
	data.IncludeRuleIDs = types.BoolValue(conn.IncludeRuleIDs)
	data.Key = types.StringValue(conn.Key)
	data.ProxyEnabled = types.BoolValue(conn.ProxyEnabled)
	data.ProxyHost = types.StringValue(conn.ProxyHost)
	data.ProxySigningKey = types.StringValue(conn.ProxySigningKey)
	data.SseEnabled = types.BoolValue(conn.SseEnabled)
	data.HashSecureAttributes = types.BoolValue(conn.HashSecureAttributes)
	data.RemoteEvalEnabled = types.BoolValue(conn.RemoteEvalEnabled)
	data.SavedGroupReferencesEnabled = types.BoolValue(conn.SavedGroupReferencesEnabled)
	data.DateCreated = types.StringValue(conn.DateCreated)
	data.DateUpdated = types.StringValue(conn.DateUpdated)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

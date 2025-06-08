package internal

import (
	"context"
	"terraform-provider-growthbook/internal/growthbookapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSDKConnection() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSDKConnectionRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the SDK Connection.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for the SDK Connection.",
			},
			"organization": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The organization associated with the SDK Connection.",
			},
			"language": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The programming language for the SDK.",
			},
			"sdk_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the SDK.",
			},
			"environment": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The environment for the SDK Connection.",
			},
			"projects": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "The projects associated with the SDK Connection.",
			},
			"encrypt_payload": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to encrypt the payload.",
			},
			"encryption_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The key used for encryption.",
			},
			"include_visual_experiments": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to include visual experiments.",
			},
			"include_draft_experiments": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to include draft experiments.",
			},
			"include_experiment_names": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to include experiment names.",
			},
			"include_redirect_experiments": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to include redirect experiments.",
			},
			"include_rule_ids": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to include rule IDs.",
			},
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The key for the SDK Connection.",
			},
			"proxy_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the proxy is enabled.",
			},
			"proxy_host": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The host of the proxy.",
			},
			"proxy_signing_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The signing key for the proxy.",
			},
			"sse_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether server-sent events are enabled.",
			},
			"hash_secure_attributes": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to hash secure attributes.",
			},
			"remote_eval_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether remote evaluation is enabled.",
			},
			"saved_group_references_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether saved group references are enabled.",
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSDKConnectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	name := d.Get("name").(string)
	sdkConn, err := client.FindSDKConnectionByName(ctx, name)
	if err != nil {
		return diag.Errorf("error reading SDK connection: %s", err)
	}
	d.SetId(sdkConn.ID)
	d.Set("name", sdkConn.Name)
	d.Set("organization", sdkConn.Organization)
	d.Set("language", sdkConn.Language)
	d.Set("sdk_version", sdkConn.SdkVersion)
	d.Set("environment", sdkConn.Environment)
	d.Set("projects", sdkConn.Projects)
	d.Set("encrypt_payload", sdkConn.EncryptPayload)
	d.Set("encryption_key", sdkConn.EncryptionKey)
	d.Set("include_visual_experiments", sdkConn.IncludeVisualExperiments)
	d.Set("include_draft_experiments", sdkConn.IncludeDraftExperiments)
	d.Set("include_experiment_names", sdkConn.IncludeExperimentNames)
	d.Set("include_redirect_experiments", sdkConn.IncludeRedirectExperiments)
	d.Set("include_rule_ids", sdkConn.IncludeRuleIDs)
	d.Set("key", sdkConn.Key)
	d.Set("proxy_enabled", sdkConn.ProxyEnabled)
	d.Set("proxy_host", sdkConn.ProxyHost)
	d.Set("proxy_signing_key", sdkConn.ProxySigningKey)
	d.Set("sse_enabled", sdkConn.SseEnabled)
	d.Set("hash_secure_attributes", sdkConn.HashSecureAttributes)
	d.Set("remote_eval_enabled", sdkConn.RemoteEvalEnabled)
	d.Set("saved_group_references_enabled", sdkConn.SavedGroupReferencesEnabled)
	d.Set("date_created", sdkConn.DateCreated)
	d.Set("date_updated", sdkConn.DateUpdated)

	return nil
}

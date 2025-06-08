package internal

import (
	"context"
	"terraform-provider-growthbook/internal/growthbookapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSDKConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSDKConnectionCreate,
		ReadContext:   resourceSDKConnectionRead,
		UpdateContext: resourceSDKConnectionUpdate,
		DeleteContext: resourceSDKConnectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"environment": {
				Type:     schema.TypeString,
				Required: true,
			},
			"language": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sdk_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"projects": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"encrypt_payload": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"include_visual_experiments": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"include_draft_experiments": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"include_experiment_names": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"include_redirect_experiments": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"include_rule_ids": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"proxy_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"proxy_host": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hash_secure_attributes": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"remote_eval_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"saved_group_references_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			// computed
			"organization": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"proxy_signing_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"sse_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"encryption_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
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

func resourceSDKConnectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	sdkConn := &growthbookapi.SDKConnection{
		Name:        d.Get("name").(string),
		Language:    d.Get("language").(string),
		Environment: d.Get("environment").(string),
	}

	if v, ok := d.GetOk("sdk_version"); ok {
		sdkConn.SdkVersion = v.(string)
	}
	if v, ok := d.GetOk("projects"); ok {
		projects := []string{}
		for _, p := range v.([]interface{}) {
			projects = append(projects, p.(string))
		}
		sdkConn.Projects = projects
	}
	if v, ok := d.GetOk("encrypt_payload"); ok {
		sdkConn.EncryptPayload = v.(bool)
	}
	if v, ok := d.GetOk("include_visual_experiments"); ok {
		sdkConn.IncludeVisualExperiments = v.(bool)
	}
	if v, ok := d.GetOk("include_draft_experiments"); ok {
		sdkConn.IncludeDraftExperiments = v.(bool)
	}
	if v, ok := d.GetOk("include_experiment_names"); ok {
		sdkConn.IncludeExperimentNames = v.(bool)
	}
	if v, ok := d.GetOk("include_redirect_experiments"); ok {
		sdkConn.IncludeRedirectExperiments = v.(bool)
	}
	if v, ok := d.GetOk("include_rule_ids"); ok {
		sdkConn.IncludeRuleIDs = v.(bool)
	}
	if v, ok := d.GetOk("proxy_enabled"); ok {
		sdkConn.ProxyEnabled = v.(bool)
	}
	if v, ok := d.GetOk("proxy_host"); ok {
		sdkConn.ProxyHost = v.(string)
	}
	if v, ok := d.GetOk("hash_secure_attributes"); ok {
		sdkConn.HashSecureAttributes = v.(bool)
	}
	if v, ok := d.GetOk("remote_eval_enabled"); ok {
		sdkConn.RemoteEvalEnabled = v.(bool)
	}
	if v, ok := d.GetOk("saved_group_references_enabled"); ok {
		sdkConn.SavedGroupReferencesEnabled = v.(bool)
	}
	created, err := client.CreateSDKConnection(ctx, sdkConn)
	if err != nil {
		return diag.Errorf("error creating SDK connection: %s", err)
	}
	d.SetId(created.ID)

	return resourceSDKConnectionRead(ctx, d, m)
}

func resourceSDKConnectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	id := d.Id()
	sdkConn, err := client.GetSDKConnection(ctx, id)
	if err != nil {
		return diag.Errorf("error reading SDK connection: %s", err)
	}
	d.Set("name", sdkConn.Name)
	d.Set("language", sdkConn.Language)
	d.Set("sdk_version", sdkConn.SdkVersion)
	d.Set("environment", sdkConn.Environment)
	d.Set("projects", sdkConn.Projects)
	d.Set("encrypt_payload", sdkConn.EncryptPayload)
	d.Set("include_visual_experiments", sdkConn.IncludeVisualExperiments)
	d.Set("include_draft_experiments", sdkConn.IncludeDraftExperiments)
	d.Set("include_experiment_names", sdkConn.IncludeExperimentNames)
	d.Set("include_redirect_experiments", sdkConn.IncludeRedirectExperiments)
	d.Set("include_rule_ids", sdkConn.IncludeRuleIDs)
	d.Set("proxy_enabled", sdkConn.ProxyEnabled)
	d.Set("proxy_host", sdkConn.ProxyHost)
	d.Set("hash_secure_attributes", sdkConn.HashSecureAttributes)
	d.Set("remote_eval_enabled", sdkConn.RemoteEvalEnabled)
	d.Set("saved_group_references_enabled", sdkConn.SavedGroupReferencesEnabled)
	// computed
	d.Set("organization", sdkConn.Organization)
	d.Set("encryption_key", sdkConn.EncryptionKey)
	d.Set("key", sdkConn.Key)
	d.Set("proxy_signing_key", sdkConn.ProxySigningKey)
	d.Set("sse_enabled", sdkConn.SseEnabled)
	d.Set("date_created", sdkConn.DateCreated)
	d.Set("date_updated", sdkConn.DateUpdated)

	return nil
}

func resourceSDKConnectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	id := d.Id()
	sdkConn := &growthbookapi.SDKConnection{
		Name:        d.Get("name").(string),
		Language:    d.Get("language").(string),
		Environment: d.Get("environment").(string),
	}

	if v, ok := d.GetOk("sdk_version"); ok {
		sdkConn.SdkVersion = v.(string)
	}
	if v, ok := d.GetOk("projects"); ok {
		projects := []string{}
		for _, p := range v.([]interface{}) {
			projects = append(projects, p.(string))
		}
		sdkConn.Projects = projects
	}
	if v, ok := d.GetOk("encrypt_payload"); ok {
		sdkConn.EncryptPayload = v.(bool)
	}
	if v, ok := d.GetOk("include_visual_experiments"); ok {
		sdkConn.IncludeVisualExperiments = v.(bool)
	}
	if v, ok := d.GetOk("include_draft_experiments"); ok {
		sdkConn.IncludeDraftExperiments = v.(bool)
	}
	if v, ok := d.GetOk("include_experiment_names"); ok {
		sdkConn.IncludeExperimentNames = v.(bool)
	}
	if v, ok := d.GetOk("include_redirect_experiments"); ok {
		sdkConn.IncludeRedirectExperiments = v.(bool)
	}
	if v, ok := d.GetOk("include_rule_ids"); ok {
		sdkConn.IncludeRuleIDs = v.(bool)
	}
	if v, ok := d.GetOk("proxy_enabled"); ok {
		sdkConn.ProxyEnabled = v.(bool)
	}
	if v, ok := d.GetOk("proxy_host"); ok {
		sdkConn.ProxyHost = v.(string)
	}
	if v, ok := d.GetOk("hash_secure_attributes"); ok {
		sdkConn.HashSecureAttributes = v.(bool)
	}
	if v, ok := d.GetOk("remote_eval_enabled"); ok {
		sdkConn.RemoteEvalEnabled = v.(bool)
	}
	if v, ok := d.GetOk("saved_group_references_enabled"); ok {
		sdkConn.SavedGroupReferencesEnabled = v.(bool)
	}
	_, err := client.UpdateSDKConnection(ctx, id, sdkConn)
	if err != nil {
		return diag.Errorf("error updating SDK connection: %s", err)
	}

	return resourceSDKConnectionRead(ctx, d, m)
}

func resourceSDKConnectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	id := d.Id()
	if err := client.DeleteSDKConnection(ctx, id); err != nil {
		return diag.Errorf("error deleting SDK connection: %s", err)
	}
	d.SetId("")

	return nil
}

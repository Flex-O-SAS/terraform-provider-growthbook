package internal

import (
	"context"
	"terraform-provider-growthbook/internal/growthbookapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFeature() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFeatureCreate,
		ReadContext:   resourceFeatureRead,
		UpdateContext: resourceFeatureUpdate,
		DeleteContext: resourceFeatureDelete,
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
			"archived": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"value_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default_value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"environments": {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString}, // For simplicity, you may want to expand this for full nested support
				Optional: true,
			},
			"prerequisites": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

func resourceFeatureCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	name := d.Get("name").(string)
	feature := &growthbookapi.Feature{
		ID:           name,
		Description:  d.Get("description").(string),
		Archived:     d.Get("archived").(bool),
		Owner:        d.Get("owner").(string),
		Project:      d.Get("project").(string),
		ValueType:    d.Get("value_type").(string),
		DefaultValue: d.Get("default_value").(string),
	}
	if v, ok := d.GetOk("tags"); ok {
		tags := []string{}
		for _, t := range v.([]interface{}) {
			tags = append(tags, t.(string))
		}
		feature.Tags = tags
	}
	if v, ok := d.GetOk("prerequisites"); ok {
		prereqs := []string{}
		for _, p := range v.([]interface{}) {
			prereqs = append(prereqs, p.(string))
		}
		feature.Prerequisites = prereqs
	}
	// Note: environments is a complex map, skipping for now unless you want full support
	created, err := client.CreateFeature(feature)
	if err != nil {
		return diag.Errorf("Failed to create feature: %s", err)
	}
	d.SetId(created.ID)
	return resourceFeatureRead(ctx, d, m)
}

func resourceFeatureRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	id := d.Id()
	feature, err := client.GetFeature(id)
	if err != nil {
		return diag.Errorf("error reading feature: %v", err)
	}
	d.Set("name", id)
	d.Set("archived", feature.Archived)
	d.Set("description", feature.Description)
	d.Set("owner", feature.Owner)
	d.Set("project", feature.Project)
	d.Set("value_type", feature.ValueType)
	d.Set("default_value", feature.DefaultValue)
	d.Set("tags", feature.Tags)
	d.Set("environments", feature.Environments)
	d.Set("prerequisites", feature.Prerequisites)
	return nil
}

func resourceFeatureUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	id := d.Id()
	feature := &growthbookapi.Feature{
		Description:  d.Get("description").(string),
		Archived:     d.Get("archived").(bool),
		Owner:        d.Get("owner").(string),
		Project:      d.Get("project").(string),
		ValueType:    d.Get("value_type").(string),
		DefaultValue: d.Get("default_value").(string),
	}
	if v, ok := d.GetOk("tags"); ok {
		tags := []string{}
		for _, t := range v.([]interface{}) {
			tags = append(tags, t.(string))
		}
		feature.Tags = tags
	}
	if v, ok := d.GetOk("prerequisites"); ok {
		prereqs := []string{}
		for _, p := range v.([]interface{}) {
			prereqs = append(prereqs, p.(string))
		}
		feature.Prerequisites = prereqs
	}
	_, err := client.UpdateFeature(id, feature)
	if err != nil {
		return diag.Errorf("Failed to update feature: %s", err)
	}
	return resourceFeatureRead(ctx, d, m)
}

func resourceFeatureDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	id := d.Id()
	if err := client.DeleteFeature(id); err != nil {
		return diag.Errorf("Failed to delete feature: %s", err)
	}
	d.SetId("")
	return nil
}

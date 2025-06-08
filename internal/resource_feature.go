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
				ForceNew: true,
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
				Type: schema.TypeMap,
				// For simplicity, you may want to expand this for full nested support
				Elem:     &schema.Schema{Type: schema.TypeString},
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
	feature := &growthbookapi.Feature{
		ID:           d.Get("name").(string),
		Description:  d.Get("description").(string),
		Owner:        d.Get("owner").(string),
		Project:      d.Get("project").(string),
		ValueType:    d.Get("value_type").(string),
		DefaultValue: d.Get("default_value").(string),
		Tags:         expandStringList(d.Get("tags")),
	}
	created, err := client.CreateFeature(ctx, feature)
	if err != nil {
		return diag.Errorf("error creating feature: %v", err)
	}
	d.SetId(created.ID)
	return resourceFeatureRead(ctx, d, m)
}

func resourceFeatureRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	id := d.Id()
	feature, err := client.GetFeature(ctx, id)
	if err != nil {
		return diag.Errorf("error reading feature: %v", err)
	}
	d.Set("id", feature.ID)
	d.Set("name", feature.ID)
	d.Set("description", feature.Description)
	d.Set("owner", feature.Owner)
	d.Set("project", feature.Project)
	d.Set("value_type", feature.ValueType)
	d.Set("default_value", feature.DefaultValue)
	d.Set("tags", feature.Tags)
	return nil
}

func resourceFeatureUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	id := d.Id()
	feature := &growthbookapi.Feature{
		Description:  d.Get("description").(string),
		Owner:        d.Get("owner").(string),
		Project:      d.Get("project").(string),
		ValueType:    d.Get("value_type").(string),
		DefaultValue: d.Get("default_value").(string),
		Tags:         expandStringList(d.Get("tags")),
	}
	_, err := client.UpdateFeature(ctx, id, feature)
	if err != nil {
		return diag.Errorf("error updating feature: %v", err)
	}
	return resourceFeatureRead(ctx, d, m)
}

func resourceFeatureDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	id := d.Id()
	if err := client.DeleteFeature(ctx, id); err != nil {
		return diag.Errorf("Failed to delete feature: %s", err)
	}
	d.SetId("")
	return nil
}

func expandStringList(input interface{}) []string {
	if input == nil {
		return nil
	}
	list := []string{}
	for _, v := range input.([]interface{}) {
		list = append(list, v.(string))
	}
	return list
}

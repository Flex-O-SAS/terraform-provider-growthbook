package internal

import (
	"context"
	"terraform-provider-growthbook/internal/growthbookapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFeature() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFeatureRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the GrowthBook feature.",
			},
			"archived": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"environments": {
				Type: schema.TypeMap,
				// For simplicity, you may want to expand this for full nested support
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"prerequisites": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"json_schema": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceFeatureRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	id := d.Get("id").(string)
	feature, err := client.GetFeature(ctx, id)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to find GrowthBook feature by id",
				Detail:   err.Error(),
			},
		}
	}
	d.SetId(feature.ID)
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

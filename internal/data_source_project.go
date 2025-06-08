package internal

import (
	"context"
	"terraform-provider-growthbook/internal/growthbookapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProjectRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the GrowthBook project.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the GrowthBook project.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the GrowthBook project.",
			},
			"stats_engine": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The stats engine used by the project.",
			},
			"date_created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation date of the project.",
			},
			"date_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last update date of the project.",
			},
		},
	}
}

func dataSourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	name := d.Get("name").(string)
	project, err := client.FindProjectByName(ctx, name)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to find GrowthBook project by name",
				Detail:   err.Error(),
			},
		}
	}
	d.SetId(project.ID)
	d.Set("name", project.Name)
	d.Set("description", project.Description)
	d.Set("stats_engine", project.Settings.StatsEngine)
	d.Set("date_created", project.DateCreated)
	d.Set("date_updated", project.DateUpdated)
	return nil
}

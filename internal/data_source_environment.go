package internal

import (
	"context"
	"terraform-provider-growthbook/internal/growthbookapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEnvironment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnvironmentRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the GrowthBook environment.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the GrowthBook environment.",
			},
			"toggle_on_list": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Show toggle on feature list.",
			},
			"default_state": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Default state for new features.",
			},
			"projects": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Projects associated with the environment.",
			},
		},
	}
}

func dataSourceEnvironmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	id := d.Get("id").(string)
	env, err := client.FindEnvironmentByID(ctx, id)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to find GrowthBook environment by ID",
				Detail:   err.Error(),
			},
		}
	}
	d.SetId(env.ID)
	d.Set("description", env.Description)
	d.Set("toggle_on_list", env.ToggleOnList)
	d.Set("default_state", env.DefaultState)
	d.Set("projects", env.Projects)

	return nil
}

package internal

import (
	"context"
	"terraform-provider-growthbook/internal/growthbookapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEnvironments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnvironmentsRead,
		Schema: map[string]*schema.Schema{
			"environments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
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
				},
			},
		},
	}
}

func dataSourceEnvironmentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	envs, err := client.ListEnvironments(ctx)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to list GrowthBook environments",
				Detail:   err.Error(),
			},
		}
	}

	items := make([]map[string]interface{}, len(envs))
	for i, env := range envs {
		items[i] = map[string]interface{}{
			"id":             env.ID,
			"description":    env.Description,
			"toggle_on_list": env.ToggleOnList,
			"default_state":  env.DefaultState,
			"projects":       env.Projects,
		}
	}
	d.SetId("growthbook_environments")
	d.Set("environments", items)

	return nil
}

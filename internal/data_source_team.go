package internal

import (
	"context"
	"terraform-provider-growthbook/internal/growthbookapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTeam() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTeamRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the GrowthBook team.",
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The global role for team members.",
			},
			"limit_access_by_environment": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"environments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"project_roles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"limit_access_by_environment": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"environments": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"members": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"managed_by_idp": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
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

func dataSourceTeamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	name := d.Get("name").(string)
	team, err := client.FindTeamByName(ctx, name)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to find GrowthBook team by name",
				Detail:   err.Error(),
			},
		}
	}
	d.SetId(team.ID)
	d.Set("name", team.Name)
	d.Set("description", team.Description)
	d.Set("role", team.Role)
	d.Set("limit_access_by_environment", team.LimitAccessByEnvironment)
	d.Set("environments", team.Environments)
	d.Set("project_roles", flattenProjectRoles(team.ProjectRoles))
	d.Set("members", team.Members)
	d.Set("managed_by_idp", team.ManagedByIdp)
	d.Set("created_by", team.CreatedBy)
	d.Set("date_created", team.DateCreated)
	d.Set("date_updated", team.DateUpdated)
	return nil
}

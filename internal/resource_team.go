package internal

import (
	"context"
	"terraform-provider-growthbook/internal/growthbookapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func projectRolesSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"project": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The project ID.",
				},
				"role": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The role for this project.",
				},
				"limit_access_by_environment": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Whether to limit access by environment.",
				},
				"environments": {
					Type:     schema.TypeList,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
					Description: "List of environments the role applies to. " +
						"Empty means all environments.",
				},
			},
		},
	}
}

func expandProjectRoles(raw []interface{}) []growthbookapi.ProjectRole {
	roles := make([]growthbookapi.ProjectRole, len(raw))
	for i, v := range raw {
		m := v.(map[string]interface{})
		roles[i] = growthbookapi.ProjectRole{
			Project:                  m["project"].(string),
			Role:                     m["role"].(string),
			LimitAccessByEnvironment: m["limit_access_by_environment"].(bool),
			Environments:             expandStringList(m["environments"].([]interface{})),
		}
	}
	return roles
}

func flattenProjectRoles(roles []growthbookapi.ProjectRole) []map[string]interface{} {
	result := make([]map[string]interface{}, len(roles))
	for i, r := range roles {
		result[i] = map[string]interface{}{
			"project":                     r.Project,
			"role":                        r.Role,
			"limit_access_by_environment": r.LimitAccessByEnvironment,
			"environments":                r.Environments,
		}
	}
	return result
}

func resourceTeam() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTeamCreate,
		ReadContext:   resourceTeamRead,
		UpdateContext: resourceTeamUpdate,
		DeleteContext: resourceTeamDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The team name.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The team description.",
			},
			"role": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The global role for team members (e.g. readonly, collaborator, engineer, analyst, experimenter, admin).",
			},
			"limit_access_by_environment": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to restrict team access to specific environments.",
			},
			"environments": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of environment IDs the team has access to. Empty means all.",
			},
			"project_roles": projectRolesSchema(),
			"members": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Set of member IDs in this team.",
			},
			"managed_by_idp": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the team is managed by an identity provider.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user who created the team.",
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

func resourceTeamCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	team := &growthbookapi.Team{
		Name:                     d.Get("name").(string),
		Description:              d.Get("description").(string),
		Role:                     d.Get("role").(string),
		LimitAccessByEnvironment: d.Get("limit_access_by_environment").(bool),
		Environments:             expandStringList(d.Get("environments").([]interface{})),
		ProjectRoles:             expandProjectRoles(d.Get("project_roles").([]interface{})),
		Members:                  expandStringSet(d.Get("members").(*schema.Set)),
	}
	created, err := client.CreateTeam(ctx, team)
	if err != nil {
		return diag.Errorf("error creating team: %v", err)
	}
	d.SetId(created.ID)

	return resourceTeamRead(ctx, d, m)
}

func resourceTeamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	team, err := client.GetTeam(ctx, d.Id())
	if err != nil {
		return diag.Errorf("error reading team: %v", err)
	}
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

func resourceTeamUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	team := &growthbookapi.Team{
		Name:                     d.Get("name").(string),
		Description:              d.Get("description").(string),
		Role:                     d.Get("role").(string),
		LimitAccessByEnvironment: d.Get("limit_access_by_environment").(bool),
		Environments:             expandStringList(d.Get("environments").([]interface{})),
		ProjectRoles:             expandProjectRoles(d.Get("project_roles").([]interface{})),
		Members:                  expandStringSet(d.Get("members").(*schema.Set)),
	}
	_, err := client.UpdateTeam(ctx, d.Id(), team)
	if err != nil {
		return diag.Errorf("error updating team: %v", err)
	}

	return resourceTeamRead(ctx, d, m)
}

func resourceTeamDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	if err := client.DeleteTeam(ctx, d.Id()); err != nil {
		return diag.Errorf("error deleting team: %v", err)
	}
	d.SetId("")

	return nil
}

func expandStringSet(set *schema.Set) []string {
	list := set.List()
	result := make([]string, len(list))
	for i, v := range list {
		result[i] = v.(string)
	}
	return result
}

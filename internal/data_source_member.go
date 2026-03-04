package internal

import (
	"context"
	"terraform-provider-growthbook/internal/growthbookapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMember() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMemberRead,
		Schema: map[string]*schema.Schema{
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The email of the GrowthBook member.",
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The member's global role.",
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
			"teams": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"managed_by_idp": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"last_login_date": {
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

func dataSourceMemberRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	email := d.Get("email").(string)
	member, err := client.FindMemberByEmail(ctx, email)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to find GrowthBook member by email",
				Detail:   err.Error(),
			},
		}
	}
	d.SetId(member.ID)
	d.Set("email", member.Email)
	d.Set("name", member.Name)
	d.Set("role", member.Role)
	d.Set("limit_access_by_environment", member.LimitAccessByEnvironment)
	d.Set("environments", member.Environments)
	d.Set("project_roles", flattenProjectRoles(member.ProjectRoles))
	d.Set("teams", member.Teams)
	d.Set("managed_by_idp", member.ManagedByIdp)
	d.Set("last_login_date", member.LastLoginDate)
	d.Set("date_created", member.DateCreated)
	d.Set("date_updated", member.DateUpdated)
	return nil
}

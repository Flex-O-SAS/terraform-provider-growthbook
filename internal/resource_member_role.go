package internal

import (
	"context"
	"terraform-provider-growthbook/internal/growthbookapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMemberRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMemberRoleCreate,
		ReadContext:   resourceMemberRoleRead,
		UpdateContext: resourceMemberRoleUpdate,
		DeleteContext: resourceMemberRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The member's email address. Used to look up the existing member.",
			},
			"role": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The member's global role (e.g. readonly, collaborator, engineer, analyst, experimenter, admin).",
			},
			"limit_access_by_environment": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to restrict access to specific environments.",
			},
			"environments": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of environments the member has access to. Empty means all.",
			},
			"project_roles": projectRolesSchema(),
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The member's display name.",
			},
			"teams": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Team IDs the member belongs to.",
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

func resourceMemberRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	email := d.Get("email").(string)

	// Look up the member by email
	member, err := client.FindMemberByEmail(ctx, email)
	if err != nil {
		return diag.Errorf("error finding member with email %s: %v", email, err)
	}

	// Update the member's role
	update := &growthbookapi.MemberRoleUpdate{
		Role:                     d.Get("role").(string),
		LimitAccessByEnvironment: d.Get("limit_access_by_environment").(bool),
		Environments:             expandStringList(d.Get("environments").([]interface{})),
		ProjectRoles:             expandProjectRoles(d.Get("project_roles").([]interface{})),
	}
	_, err = client.UpdateMemberRole(ctx, member.ID, update)
	if err != nil {
		return diag.Errorf("error updating member role: %v", err)
	}
	d.SetId(member.ID)

	return resourceMemberRoleRead(ctx, d, m)
}

func resourceMemberRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	member, err := client.GetMember(ctx, d.Id())
	if err != nil {
		return diag.Errorf("error reading member: %v", err)
	}
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

func resourceMemberRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	update := &growthbookapi.MemberRoleUpdate{
		Role:                     d.Get("role").(string),
		LimitAccessByEnvironment: d.Get("limit_access_by_environment").(bool),
		Environments:             expandStringList(d.Get("environments").([]interface{})),
		ProjectRoles:             expandProjectRoles(d.Get("project_roles").([]interface{})),
	}
	_, err := client.UpdateMemberRole(ctx, d.Id(), update)
	if err != nil {
		return diag.Errorf("error updating member role: %v", err)
	}

	return resourceMemberRoleRead(ctx, d, m)
}

func resourceMemberRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	if err := client.DeleteMember(ctx, d.Id()); err != nil {
		return diag.Errorf("error removing member from organization: %v", err)
	}
	d.SetId("")

	return nil
}

package internal

import (
	"context"
	"terraform-provider-growthbook/internal/growthbookapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentCreate,
		ReadContext:   resourceEnvironmentRead,
		UpdateContext: resourceEnvironmentUpdate,
		DeleteContext: resourceEnvironmentDelete,
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
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"toggle_on_list": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"default_state": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"projects": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	env := &growthbookapi.Environment{
		ID:           d.Get("name").(string),
		Description:  d.Get("description").(string),
		ToggleOnList: d.Get("toggle_on_list").(bool),
		DefaultState: d.Get("default_state").(bool),
	}
	if v, ok := d.GetOk("projects"); ok {
		projects := []string{}
		for _, p := range v.([]interface{}) {
			projects = append(projects, p.(string))
		}
		env.Projects = projects
	}
	created, err := client.CreateEnvironment(ctx, env)
	if err != nil {
		return diag.Errorf("error creating environment: %v", err)
	}
	d.SetId(created.ID)

	return resourceEnvironmentRead(ctx, d, m)
}

func resourceEnvironmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	id := d.Id()
	env, err := client.FindEnvironmentByID(ctx, id)
	if err != nil {
		return diag.Errorf("error reading environment: %v", err)
	}
	d.Set("name", env.ID)
	d.Set("description", env.Description)
	d.Set("toggle_on_list", env.ToggleOnList)
	d.Set("default_state", env.DefaultState)
	d.Set("projects", env.Projects)

	return nil
}

func resourceEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	id := d.Id()
	env := &growthbookapi.Environment{
		Description:  d.Get("description").(string),
		ToggleOnList: d.Get("toggle_on_list").(bool),
		DefaultState: d.Get("default_state").(bool),
	}
	if v, ok := d.GetOk("projects"); ok {
		projects := []string{}
		for _, p := range v.([]interface{}) {
			projects = append(projects, p.(string))
		}
		env.Projects = projects
	}
	_, err := client.UpdateEnvironment(ctx, id, env)
	if err != nil {
		return diag.Errorf("error updating environment: %v", err)
	}

	return resourceEnvironmentRead(ctx, d, m)
}

func resourceEnvironmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	id := d.Id()
	if err := client.DeleteEnvironment(ctx, id); err != nil {
		return diag.Errorf("Failed to delete environment: %s", err)
	}
	d.SetId("")

	return nil
}

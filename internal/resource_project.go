package internal

import (
	"context"
	"terraform-provider-growthbook/internal/growthbookapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,
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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"stats_engine": {
				Type:     schema.TypeString,
				Optional: true,
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

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	project := &growthbookapi.Project{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
	if v, ok := d.GetOk("stats_engine"); ok {
		project.Settings.StatsEngine = v.(string)
	}
	created, err := client.CreateProject(ctx, project)
	if err != nil {
		return diag.Errorf("error creating project: %v", err)
	}
	d.SetId(created.ID)

	return resourceProjectRead(ctx, d, m)
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	id := d.Id()
	project, err := client.GetProject(ctx, id)
	if err != nil {
		return diag.Errorf("error reading project: %v", err)
	}
	d.Set("name", project.Name)
	d.Set("description", project.Description)
	if project.Settings.StatsEngine != "" {
		d.Set("stats_engine", project.Settings.StatsEngine)
	}
	d.Set("date_created", project.DateCreated)
	d.Set("date_updated", project.DateUpdated)

	return nil
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	id := d.Id()
	if d.Get("name").(string) == "" {
		return diag.Errorf("Name must be provided for the project.")
	}
	project := &growthbookapi.Project{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
	if v, ok := d.GetOk("stats_engine"); ok {
		project.Settings.StatsEngine = v.(string)
	}
	updated, err := client.UpdateProject(ctx, id, project)
	if err != nil {
		return diag.Errorf("error updating project: %v", err)
	}
	d.Set("date_updated", updated.DateUpdated)

	return resourceProjectRead(ctx, d, m)
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	id := d.Id()
	if err := client.DeleteProject(ctx, id); err != nil {
		return diag.Errorf("Failed to delete project: %s", err)
	}
	d.SetId("")

	return nil
}

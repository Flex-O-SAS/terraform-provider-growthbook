package internal

import (
	"context"

	"terraform-provider-growthbook/internal/growthbookapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAttribute() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAttributeRead,		
		Schema: map[string]*schema.Schema{
			"property": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
			},
			"datatype": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			"format": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			"enum_values": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
			"projects": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
				Description: "Array of project Id",
			},
			"archived": &schema.Schema{
				Type: schema.TypeBool,
				Computed: true,
			},
			"description": &schema.Schema{
				Type: schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAttributeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	property := d.Get("property").(string)
	attribute, err := client.GetAttribute(ctx, property)
	if (err != nil) {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to find GrowthBook attribute by property",
				Detail:   err.Error(),
			},
		}
	}
	d.SetId(attribute.Property)
	d.Set("property",		attribute.Property)
	d.Set("datatype",		attribute.DataType)
	d.Set("format",			attribute.Format)
	d.Set("enum_values",	attribute.EnumValues)
	d.Set("projects",		attribute.Projects)
	d.Set("archived",		attribute.Archived)
	d.Set("description",	attribute.Description)
	return nil
}
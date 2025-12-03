package internal

import (
	"context"
	"terraform-provider-growthbook/internal/growthbookapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"sync"
	"slices"
)

var attributeMutex sync.Mutex

func	resourceAttribute()	*schema.Resource {
	return &schema.Resource{
		CreateContext:	resourceAttributeCreate,
		ReadContext: resourceAttributeRead,
		UpdateContext: resourceAttributeUpdate,
		DeleteContext: resourceAttributeDelete,
		
		Schema: map[string]*schema.Schema{
			"property": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
			},
			"datatype": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
			},
			"format": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
			},
			"enum_values": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
			},
			"projects": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Description: "Array of project Id",
			},
			"archived": &schema.Schema{
				Type: schema.TypeBool,
				Optional: true,
			},
			"description": &schema.Schema{
				Type: schema.TypeString,
				Optional: true,
			},
		},
	}
}

func	resourceAttributeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	attribute := &growthbookapi.Attribute{
		Property: 		d.Get("property").(string),
		DataType: 		d.Get("datatype").(string),
		Format:			d.Get("format").(string),
		EnumValues:		d.Get("enum_values").(string),
		Projects:		retrieveStrings(d.Get("projects").([]interface{})),
		Archived:		d.Get("archived").(bool),
		Description:	d.Get("description").(string),
	}
	out, err := client.GetAttribute(ctx, attribute.Property)
	if err != nil {
		return diag.Errorf("error reading attribute: %v", err)
	}
	if out == nil {
		return diag.Errorf("error reading attribute: %v", attribute.Property)
	}
	d.Set("property", 		out.Property)
	d.Set("datatype", 		out.DataType)
	d.Set("format",   		out.Format)
	d.Set("enum_values",	out.EnumValues)
	d.Set("projects",		out.Projects)
	d.Set("archived",		out.Archived)
	d.Set("description",	out.Description)
	return nil
}

func	resourceAttributeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	attribute := &growthbookapi.Attribute{
		Property:		d.Get("property").(string),
		DataType:		d.Get("datatype").(string),
		Format:			d.Get("format").(string),
		EnumValues:		d.Get("enum_values").(string),
		Projects:		retrieveStrings(d.Get("projects").([]interface{})),
		Archived:		d.Get("archived").(bool),
		Description:	d.Get("description").(string),
	}

	FormatAcceptedValue := []string{"", "version", "date", "isoCountryCode"}
	if slices.Contains(FormatAcceptedValue, attribute.Format) == false {
		return diag.Errorf("[format] Invalid value. Expected '' | 'version' | 'date' | 'isoCountryCode', received %v", attribute.Format)
	}
	
	// GrowthBook does not handle concurrent attribute creation properly and may enter a data race state. 
	// To prevent this issue, we enforce serialized execution using a mutex, ensuring attributes are created sequentially.
	attributeMutex.Lock()
	defer attributeMutex.Unlock()
	created, err := client.CreateAttribute(ctx, attribute)

	if err != nil {
		return diag.Errorf("error creating attribute %v %v", err, created)
	}
	d.SetId(created.Property)
	return resourceAttributeRead(ctx, d, m)
}

func	resourceAttributeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	attribute := &growthbookapi.Attribute{
		Property: 		d.Get("property").(string),
		DataType: 		d.Get("datatype").(string),
		Format:	  		d.Get("format").(string),
		EnumValues:		d.Get("enum_values").(string),
		Projects:		retrieveStrings(d.Get("projects").([]interface{})),
		Archived:		d.Get("archived").(bool),
		Description:	d.Get("description").(string),
	}
	_, err := client.UpdateAttribute(ctx, attribute.Property, attribute)
	if err != nil {
		return diag.Errorf("error updating attribute: %v", err)
	}
	return resourceAttributeRead(ctx, d, m)
}

func	resourceAttributeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*growthbookapi.Client)
	attribute := &growthbookapi.Attribute{
		Property: 		d.Get("property").(string),
		DataType: 		d.Get("datatype").(string),
		Format:	  		d.Get("format").(string),
		EnumValues:		d.Get("enum_values").(string),
		Projects:		retrieveStrings(d.Get("projects").([]interface{})),
		Archived:		d.Get("archived").(bool),
		Description:	d.Get("description").(string),
	}
	if err := client.DeleteAttribute(ctx, attribute.Property); err != nil {
		return diag.Errorf("Failed to delete attribute: %s", err)
	}
	d.SetId("")
	return nil
}

func	retrieveStrings(p []interface{}) ([]string) {
	projects := make([]string, len(p))
	for i := 0; i < len(p); i++ {
		projects[i] = p[i].(string)
	}
	return projects
}
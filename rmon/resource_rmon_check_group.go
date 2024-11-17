package rmon

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCheckGroup() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceCheckGroupCreate,
		ReadWithoutTimeout:   resourceCheckGroupRead,
		UpdateWithoutTimeout: resourceCheckGroupUpdate,
		DeleteWithoutTimeout: resourceCheckGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Description: "Managing Check groups.",

		Schema: map[string]*schema.Schema{
			NameField: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the check group.",
			},
			GroupIDField: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The user group ID.",
			},
		},
	}
}

func resourceCheckGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	name := d.Get(NameField).(string)

	requestBody := map[string]interface{}{NameField: name, GroupIDField: d.Get(GroupIDField).(int)}
	resp, err := client.doRequest("POST", "/api/v1.0/rmon/check-group", requestBody)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("API response: %s", resp)

	// Assuming the response contains an ID field with the unique identifier
	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return diag.FromErr(err)
	}

	id, ok := result[IDField].(float64) // ID возвращается как число
	if !ok {
		return diag.Errorf("unable to find ID in response: %v", result)
	}

	d.SetId(fmt.Sprintf("%d", int(id))) // Преобразование ID в строку
	return resourceCheckGroupRead(ctx, d, m)
}

func resourceCheckGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	// Implement API call to read the resource
	resp, err := client.doRequest("GET", fmt.Sprintf("/api/v1.0/rmon/check-group/%s", id), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	// Process response and set data
	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return diag.FromErr(err)
	}

	if name, ok := result[NameField].(string); ok {
		d.Set(NameField, name)
	}

	if groupId, ok := result[GroupIDField].(int); ok {
		d.Set(GroupIDField, groupId)
	}

	return nil
}

func resourceCheckGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	requestBody := map[string]interface{}{NameField: d.Get(NameField).(string), GroupIDField: d.Get(GroupIDField).(int)}

	_, err := client.doRequest("PUT", fmt.Sprintf("/api/v1.0/rmon/check-group/%s", id), requestBody)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCheckGroupRead(ctx, d, m)
}

func resourceCheckGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	// Implement API call to delete the resource
	_, err := client.doRequest("DELETE", fmt.Sprintf("/api/v1.0/rmon/check-group/%s", id), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

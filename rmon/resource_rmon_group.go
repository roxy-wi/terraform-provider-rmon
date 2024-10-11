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

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceGroupCreate,
		ReadWithoutTimeout:   resourceGroupRead,
		UpdateWithoutTimeout: resourceGroupUpdate,
		DeleteWithoutTimeout: resourceGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Description: "Represent roxy-wi group. All servers managed via RMON can be included in groups. A group is a user-defined pool of servers. It is up to you decide how you will group your servers. You can create groups because of their type, purpose, etc. By default, all your servers are included in the common group named Default. All other groups are created within this group.",

		Schema: map[string]*schema.Schema{
			NameField: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the group.",
			},
			DescriptionField: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the group.",
			},
		},
	}
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	name := d.Get(NameField).(string)
	description := d.Get(DescriptionField).(string)

	requestBody := map[string]string{NameField: name, DescriptionField: description}
	resp, err := client.doRequest("POST", "/api/v1.0/group", requestBody)
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
	return resourceGroupRead(ctx, d, m)
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	// Implement API call to read the resource
	resp, err := client.doRequest("GET", fmt.Sprintf("/api/v1.0/group/%s", id), nil)
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

	if description, ok := result[DescriptionField].(string); ok {
		d.Set(DescriptionField, description)
	}

	return nil
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	requestBody := map[string]string{NameField: d.Get(NameField).(string)}
	if d.HasChange(DescriptionField) {
		requestBody[DescriptionField] = d.Get(DescriptionField).(string)
	}

	_, err := client.doRequest("PUT", fmt.Sprintf("/api/v1.0/group/%s", id), requestBody)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGroupRead(ctx, d, m)
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	// Implement API call to delete the resource
	_, err := client.doRequest("DELETE", fmt.Sprintf("/api/v1.0/group/%s", id), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

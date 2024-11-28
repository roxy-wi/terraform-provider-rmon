package rmon

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCountry() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceCountryCreate,
		ReadWithoutTimeout:   resourceCountryRead,
		UpdateWithoutTimeout: resourceCountryUpdate,
		DeleteWithoutTimeout: resourceCountryDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Description: "This resource manages Countries in RMON.",

		Schema: map[string]*schema.Schema{
			NameField: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Country.",
			},
			DescriptionField: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the Country.",
			},
			EnabledField: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enabled state of the Country.",
			},
			SharedField: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Is the Country shared with other groups?.",
			},
		},
	}
}

func resourceCountryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client

	description := strings.ReplaceAll(d.Get(DescriptionField).(string), "'", "")
	name := strings.ReplaceAll(d.Get(NameField).(string), "'", "")
	shred := boolToInt(d.Get(SharedField).(bool))
	enabled := boolToInt(d.Get(EnabledField).(bool))

	server := map[string]interface{}{
		SharedField:      shred,
		DescriptionField: description,
		EnabledField:     enabled,
		NameField:        name,
	}

	resp, err := client.doRequest("POST", "/api/v1.0/rmon/country", server)
	if err != nil {
		return diag.FromErr(err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return diag.FromErr(err)
	}

	id, ok := result["id"].(float64)
	if !ok {
		return diag.Errorf("unable to find ID in response: %v", result)
	}

	d.SetId(fmt.Sprintf("%d", int(id)))
	return resourceCountryRead(ctx, d, m)
}

func resourceCountryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	resp, err := client.doRequest("GET", fmt.Sprintf("/api/v1.0/rmon/country/%s", id), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return diag.FromErr(err)
	}

	description := strings.ReplaceAll(result[DescriptionField].(string), "'", "")
	name := strings.ReplaceAll(result[NameField].(string), "'", "")
	d.Set(DescriptionField, description)
	d.Set(EnabledField, result[EnabledField])
	d.Set(SharedField, result[SharedField])
	d.Set(NameField, name)

	return nil
}

func resourceCountryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	description := strings.ReplaceAll(d.Get(DescriptionField).(string), "'", "")
	name := strings.ReplaceAll(d.Get(NameField).(string), "'", "")
	shred := boolToInt(d.Get(SharedField).(bool))
	enabled := boolToInt(d.Get(EnabledField).(bool))

	server := map[string]interface{}{
		SharedField:      shred,
		DescriptionField: description,
		EnabledField:     enabled,
		NameField:        name,
	}

	_, err := client.doRequest("PUT", fmt.Sprintf("/api/v1.0/rmon/country/%s", id), server)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCountryRead(ctx, d, m)
}

func resourceCountryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	_, err := client.doRequest("DELETE", fmt.Sprintf("/api/v1.0/rmon/country/%s", id), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

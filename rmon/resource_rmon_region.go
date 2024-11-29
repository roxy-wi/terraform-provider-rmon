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

const (
	CountryField = "country_id"
)

func resourceRegion() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceRegionCreate,
		ReadWithoutTimeout:   resourceRegionRead,
		UpdateWithoutTimeout: resourceRegionUpdate,
		DeleteWithoutTimeout: resourceRegionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Description: "This resource manages Regions in RMON.",

		Schema: map[string]*schema.Schema{
			NameField: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Region.",
			},
			DescriptionField: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the Region.",
			},
			EnabledField: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enabled state of the Region.",
			},
			SharedField: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Is the Region shared with other groups?.",
			},
			CountryField: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Country ID to what the Region belongs to.",
			},
			GroupIDField: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Group ID.",
			},
		},
	}
}

func resourceRegionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		CountryField:     d.Get(CountryField).(int),
		GroupIDField:     d.Get(GroupIDField).(int),
	}

	resp, err := client.doRequest("POST", "/api/v1.0/rmon/region", server)
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
	return resourceRegionRead(ctx, d, m)
}

func resourceRegionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	resp, err := client.doRequest("GET", fmt.Sprintf("/api/v1.0/rmon/region/%s", id), nil)
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
	d.Set(CountryField, result[CountryField])
	d.Set(GroupIDField, result[GroupIDField])

	return nil
}

func resourceRegionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		CountryField:     d.Get(CountryField).(int),
		GroupIDField:     d.Get(GroupIDField).(int),
	}

	_, err := client.doRequest("PUT", fmt.Sprintf("/api/v1.0/rmon/region/%s", id), server)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRegionRead(ctx, d, m)
}

func resourceRegionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	_, err := client.doRequest("DELETE", fmt.Sprintf("/api/v1.0/rmon/region/%s", id), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

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
	ServerIdField    = "server_id"
	RegionIdFiled    = "region_id"
	SharedAgentField = "shared"
)

func resourceAgent() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAgentCreate,
		ReadWithoutTimeout:   resourceAgentRead,
		UpdateWithoutTimeout: resourceAgentUpdate,
		DeleteWithoutTimeout: resourceAgentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Description: "This resource manages Agents in RMON.",

		Schema: map[string]*schema.Schema{
			NameField: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Agent.",
			},
			DescriptionField: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the Agent.",
			},
			EnabledField: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enabled state of the Agent.",
			},
			SharedAgentField: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Is the Agent shared with other groups?.",
			},
			ServerIdField: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID of the server where Agent will be installed.",
			},
			PortField: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Port number for binding Agent.",
			},
			RegionIdFiled: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the region to which the agent belongs.",
			},
		},
	}
}

func resourceAgentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client

	description := strings.ReplaceAll(d.Get(DescriptionField).(string), "'", "")
	name := strings.ReplaceAll(d.Get(NameField).(string), "'", "")
	shred := boolToInt(d.Get(SharedAgentField).(bool))
	enabled := boolToInt(d.Get(EnabledField).(bool))

	server := map[string]interface{}{
		SharedAgentField: shred,
		DescriptionField: description,
		EnabledField:     enabled,
		NameField:        name,
		ServerIdField:    d.Get(ServerIdField).(int),
		PortField:        d.Get(PortField).(int),
		RegionIdFiled:    d.Get(RegionIdFiled).(int),
	}

	resp, err := client.doRequest("POST", "/api/v1.0/rmon/agent", server)
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
	return resourceAgentRead(ctx, d, m)
}

func resourceAgentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	resp, err := client.doRequest("GET", fmt.Sprintf("/api/v1.0/rmon/agent/%s", id), nil)
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
	d.Set(EnabledField, intToBool(result[EnabledField].(float64)))
	d.Set(SharedAgentField, intToBool(result[SharedAgentField].(float64)))
	d.Set(NameField, name)
	d.Set(ServerIdField, result[ServerIdField])
	d.Set(PortField, result[PortField])
	d.Set(RegionIdFiled, result[RegionIdFiled])

	return nil
}

func resourceAgentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	description := strings.ReplaceAll(d.Get(DescriptionField).(string), "'", "")
	name := strings.ReplaceAll(d.Get(NameField).(string), "'", "")

	server := map[string]interface{}{
		DescriptionField: description,
		EnabledField:     boolToInt(d.Get(EnabledField).(bool)),
		SharedAgentField: boolToInt(d.Get(SharedAgentField).(bool)),
		NameField:        name,
		ServerIdField:    d.Get(ServerIdField).(int),
		PortField:        d.Get(PortField).(int),
		RegionIdFiled:    d.Get(RegionIdFiled).(int),
	}

	_, err := client.doRequest("PUT", fmt.Sprintf("/api/v1.0/rmon/agent/%s", id), server)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAgentRead(ctx, d, m)
}

func resourceAgentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	_, err := client.doRequest("DELETE", fmt.Sprintf("/api/v1.0/rmon/agent/%s", id), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

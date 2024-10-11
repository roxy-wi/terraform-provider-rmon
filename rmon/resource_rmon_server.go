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
	CredIDField   = "cred_id"
	EnabledField  = "enabled"
	GroupIDField  = "group_id"
	HostnameField = "hostname"
	IPField       = "ip"
	PortField     = "port"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceServerCreate,
		ReadWithoutTimeout:   resourceServerRead,
		UpdateWithoutTimeout: resourceServerUpdate,
		DeleteWithoutTimeout: resourceServerDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Description: "This resource manages server entries in RMON. It allows you to define servers with specific credentials, group associations, IP addresses, and other configurations.",

		Schema: map[string]*schema.Schema{
			CredIDField: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Credentials ID.",
			},
			DescriptionField: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the server.",
			},
			EnabledField: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enabled state of the server.",
			},
			GroupIDField: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Group ID.",
			},
			HostnameField: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Hostname of the server.",
			},
			IPField: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "IP address of the server.",
			},
			PortField: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Port number.",
			},
		},
	}
}

func resourceServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client

	description := strings.ReplaceAll(d.Get(DescriptionField).(string), "'", "")
	hostname := strings.ReplaceAll(d.Get(HostnameField).(string), "'", "")

	server := map[string]interface{}{
		CredIDField:      d.Get(CredIDField).(int),
		DescriptionField: description,
		EnabledField:     boolToInt(d.Get(EnabledField).(bool)),
		GroupIDField:     d.Get(GroupIDField).(int),
		HostnameField:    hostname,
		IPField:          d.Get(IPField).(string),
		PortField:        d.Get(PortField).(int),
	}

	resp, err := client.doRequest("POST", "/api/v1.0/server", server)
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
	return resourceServerRead(ctx, d, m)
}

func resourceServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	resp, err := client.doRequest("GET", fmt.Sprintf("/api/v1.0/server/%s", id), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return diag.FromErr(err)
	}

	d.Set(CredIDField, result[CredIDField])
	description := strings.ReplaceAll(result[DescriptionField].(string), "'", "")
	hostname := strings.ReplaceAll(result[HostnameField].(string), "'", "")
	d.Set(DescriptionField, description)
	d.Set(EnabledField, intToBool(result[EnabledField].(float64)))
	d.Set(GroupIDField, result[GroupIDField])
	d.Set(HostnameField, hostname)
	d.Set(IPField, result[IPField])
	d.Set(PortField, result[PortField])

	return nil
}

func resourceServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	description := strings.ReplaceAll(d.Get(DescriptionField).(string), "'", "")
	hostname := strings.ReplaceAll(d.Get(HostnameField).(string), "'", "")

	server := map[string]interface{}{
		CredIDField:      d.Get(CredIDField).(int),
		DescriptionField: description,
		EnabledField:     boolToInt(d.Get(EnabledField).(bool)),
		GroupIDField:     d.Get(GroupIDField).(int),
		HostnameField:    hostname,
		IPField:          d.Get(IPField).(string),
		PortField:        d.Get(PortField).(int),
	}

	_, err := client.doRequest("PUT", fmt.Sprintf("/api/v1.0/server/%s", id), server)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceServerRead(ctx, d, m)
}

func resourceServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	_, err := client.doRequest("DELETE", fmt.Sprintf("/api/v1.0/server/%s", id), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

package rmon

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceCheckDns() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceCheckDnsCreate,
		ReadWithoutTimeout:   resourceCheckDnsRead,
		UpdateWithoutTimeout: resourceCheckDnsUpdate,
		DeleteWithoutTimeout: resourceCheckDnsDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Description: "This resource manages DNS check in RMON.",

		Schema: map[string]*schema.Schema{
			NameField: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the CheckDns.",
			},
			DescriptionField: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the CheckDns.",
			},
			EnabledField: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enabled state of the Check DNS.",
			},
			CheckGroupIdFiled: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the check group for group DNS checks.",
			},
			PlaceField: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Port number for binding Check DNS.",
				ValidateFunc: validation.StringInSlice([]string{
					"all",
					"country",
					"region",
					"agent",
				}, false),
			},
			EntitiesField: {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of entities where check must be created.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			IntervalField: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Interval in seconds between checks.",
			},
			TimeoutField: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Answer timeout in seconds.",
			},
			TelegramField: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Telegram channel ID for alerts.",
			},
			SlackField: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Slack channel ID for alerts.",
			},
			MMField: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Mattermost channel ID for alerts.",
			},
			PDField: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "PagerDuty channel ID for alerts.",
			},
			IPField: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IP address or domain name for check.",
			},
			PortField: {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "Packet size in bytes.",
				ValidateFunc: validation.IsPortNumber,
			},
			ResolverField: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "DNS server where resolve DNS query.",
			},
			RecordTypeField: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "DNS record type.",
				ValidateFunc: validation.StringInSlice([]string{
					"a",
					"aaa",
					"caa",
					"cname",
					"mx",
					"ns",
					"ptr",
					"sao",
					"src",
					"txt",
				}, false),
			},
			RetriesField: {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "Number of retries before check is marked down.",
				ValidateFunc: validation.IntAtLeast(0),
				Default:      3,
			},
		},
	}
}

func resourceCheckDnsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client

	description := strings.ReplaceAll(d.Get(DescriptionField).(string), "'", "")
	name := strings.ReplaceAll(d.Get(NameField).(string), "'", "")
	enabled := boolToInt(d.Get(EnabledField).(bool))

	server := map[string]interface{}{
		DescriptionField:  description,
		EnabledField:      enabled,
		NameField:         name,
		CheckGroupIdFiled: d.Get(CheckGroupIdFiled).(string),
		PlaceField:        d.Get(PlaceField).(string),
		EntitiesField:     d.Get(EntitiesField),
		IntervalField:     d.Get(IntervalField).(int),
		TimeoutField:      d.Get(TimeoutField).(int),
		TelegramField:     d.Get(TelegramField).(int),
		SlackField:        d.Get(SlackField).(int),
		MMField:           d.Get(MMField).(int),
		PDField:           d.Get(PDField).(int),
		PortField:         d.Get(PortField).(int),
		IPField:           d.Get(IPField),
		ResolverField:     d.Get(ResolverField),
		RecordTypeField:   d.Get(RecordTypeField),
		RetriesField:      d.Get(RetriesField).(int),
	}

	resp, err := client.doRequest("POST", "/api/v1.0/rmon/check/dns", server)
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
	return resourceCheckDnsRead(ctx, d, m)
}

func resourceCheckDnsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	resp, err := client.doRequest("GET", fmt.Sprintf("/api/v1.0/rmon/check/dns/%s", id), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return diag.FromErr(err)
	}

	entities := result[EntitiesField].([]interface{})
	if result[PlaceField] == "all" {
		entities = []interface{}{}
	}

	description := strings.ReplaceAll(result[DescriptionField].(string), "'", "")
	name := strings.ReplaceAll(result[NameField].(string), "'", "")
	d.Set(DescriptionField, description)
	d.Set(EnabledField, intToBool(result[EnabledField].(float64)))
	d.Set(NameField, name)
	d.Set(CheckGroupIdFiled, result[CheckGroupIdFiled])
	d.Set(PlaceField, result[PlaceField])
	d.Set(EntitiesField, entities)
	d.Set(IntervalField, result[IntervalField])
	d.Set(TimeoutField, result[TimeoutField])
	d.Set(TelegramField, result[TelegramField])
	d.Set(SlackField, result[SlackField])
	d.Set(MMField, result[MMField])
	d.Set(PDField, result[PDField])
	d.Set(PortField, result[PortField])
	d.Set(IPField, result[IPField])
	d.Set(ResolverField, result[ResolverField])
	d.Set(RecordTypeField, result[RecordTypeField])
	d.Set(RetriesField, result[RetriesField])

	return nil
}

func resourceCheckDnsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	description := strings.ReplaceAll(d.Get(DescriptionField).(string), "'", "")
	name := strings.ReplaceAll(d.Get(NameField).(string), "'", "")
	enabled := boolToInt(d.Get(EnabledField).(bool))

	server := map[string]interface{}{
		DescriptionField:  description,
		EnabledField:      enabled,
		NameField:         name,
		CheckGroupIdFiled: d.Get(CheckGroupIdFiled).(string),
		PlaceField:        d.Get(PlaceField).(string),
		EntitiesField:     d.Get(EntitiesField),
		IntervalField:     d.Get(IntervalField).(int),
		TimeoutField:      d.Get(TimeoutField).(int),
		TelegramField:     d.Get(TelegramField).(int),
		SlackField:        d.Get(SlackField).(int),
		MMField:           d.Get(MMField).(int),
		PDField:           d.Get(PDField).(int),
		PortField:         d.Get(PortField).(int),
		IPField:           d.Get(IPField),
		ResolverField:     d.Get(ResolverField),
		RecordTypeField:   d.Get(RecordTypeField),
		RetriesField:      d.Get(RetriesField).(int),
	}

	_, err := client.doRequest("PUT", fmt.Sprintf("/api/v1.0/rmon/check/dns/%s", id), server)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCheckDnsRead(ctx, d, m)
}

func resourceCheckDnsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	_, err := client.doRequest("DELETE", fmt.Sprintf("/api/v1.0/rmon/check/dns/%s", id), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

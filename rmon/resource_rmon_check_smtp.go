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

func resourceCheckSmtp() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceCheckSmtpCreate,
		ReadWithoutTimeout:   resourceCheckSmtpRead,
		UpdateWithoutTimeout: resourceCheckSmtpUpdate,
		DeleteWithoutTimeout: resourceCheckSmtpDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Description: "This resource manages SMTP check in RMON.",

		Schema: map[string]*schema.Schema{
			NameField: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the CheckSmtp.",
			},
			DescriptionField: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the Check SMTP.",
			},
			EnabledField: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enabled state of the Check SMTP.",
			},
			CheckGroupIdFiled: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the check group for group SMTP checks.",
			},
			PlaceField: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Port number for binding CheckSmtp.",
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
				Required:    true,
				Description: "IP address or domain name of SMTP server for check.",
			},
			PortField: {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "SMTP server port.",
				Default:      587,
				ValidateFunc: validation.IsPortNumber,
			},
			IgnoreSslErrorField: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Ignore TLS/SSL error.",
			},
			UserNameField: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User name for authenticating to SMTP server.",
			},
			PasswordField: {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password for authenticating to SMTP server.",
			},
		},
	}
}

func resourceCheckSmtpCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client

	description := strings.ReplaceAll(d.Get(DescriptionField).(string), "'", "")
	name := strings.ReplaceAll(d.Get(NameField).(string), "'", "")
	enabled := boolToInt(d.Get(EnabledField).(bool))

	server := map[string]interface{}{
		DescriptionField:    description,
		EnabledField:        enabled,
		NameField:           name,
		CheckGroupIdFiled:   d.Get(CheckGroupIdFiled).(string),
		PlaceField:          d.Get(PlaceField).(string),
		EntitiesField:       d.Get(EntitiesField),
		IntervalField:       d.Get(IntervalField).(int),
		TimeoutField:        d.Get(TimeoutField).(int),
		TelegramField:       d.Get(TelegramField).(int),
		SlackField:          d.Get(SlackField).(int),
		MMField:             d.Get(MMField).(int),
		PDField:             d.Get(PDField).(int),
		PortField:           d.Get(PortField).(int),
		IPField:             d.Get(IPField),
		IgnoreSslErrorField: boolToInt(d.Get(IgnoreSslErrorField).(bool)),
		UserNameField:       d.Get(UserNameField),
		PasswordField:       d.Get(PasswordField),
	}

	resp, err := client.doRequest("POST", "/api/v1.0/rmon/check/smtp", server)
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
	return resourceCheckSmtpRead(ctx, d, m)
}

func resourceCheckSmtpRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	resp, err := client.doRequest("GET", fmt.Sprintf("/api/v1.0/rmon/check/smtp/%s", id), nil)
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
	d.Set(NameField, name)
	d.Set(CheckGroupIdFiled, result[CheckGroupIdFiled])
	d.Set(PlaceField, result[PlaceField])
	d.Set(EntitiesField, result[EntitiesField])
	d.Set(IntervalField, result[IntervalField])
	d.Set(TimeoutField, result[TimeoutField])
	d.Set(TelegramField, result[TelegramField])
	d.Set(SlackField, result[SlackField])
	d.Set(MMField, result[MMField])
	d.Set(PDField, result[PDField])
	d.Set(PortField, result[PortField])
	d.Set(IPField, result[IPField])
	d.Set(UserNameField, result[UserNameField])
	d.Set(PasswordField, result[PasswordField])
	d.Set(IgnoreSslErrorField, intToBool(result[IgnoreSslErrorField].(float64)))

	return nil
}

func resourceCheckSmtpUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	description := strings.ReplaceAll(d.Get(DescriptionField).(string), "'", "")
	name := strings.ReplaceAll(d.Get(NameField).(string), "'", "")
	enabled := boolToInt(d.Get(EnabledField).(bool))

	server := map[string]interface{}{
		DescriptionField:    description,
		EnabledField:        enabled,
		NameField:           name,
		CheckGroupIdFiled:   d.Get(CheckGroupIdFiled).(string),
		PlaceField:          d.Get(PlaceField).(string),
		EntitiesField:       d.Get(EntitiesField),
		IntervalField:       d.Get(IntervalField).(int),
		TimeoutField:        d.Get(TimeoutField).(int),
		TelegramField:       d.Get(TelegramField).(int),
		SlackField:          d.Get(SlackField).(int),
		MMField:             d.Get(MMField).(int),
		PDField:             d.Get(PDField).(int),
		PortField:           d.Get(PortField).(int),
		IPField:             d.Get(IPField),
		IgnoreSslErrorField: boolToInt(d.Get(IgnoreSslErrorField).(bool)),
		UserNameField:       d.Get(UserNameField),
		PasswordField:       d.Get(PasswordField),
	}

	_, err := client.doRequest("PUT", fmt.Sprintf("/api/v1.0/rmon/check/smtp/%s", id), server)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCheckSmtpRead(ctx, d, m)
}

func resourceCheckSmtpDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	_, err := client.doRequest("DELETE", fmt.Sprintf("/api/v1.0/rmon/check/smtp/%s", id), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
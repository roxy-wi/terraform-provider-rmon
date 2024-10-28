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

func resourceCheckHttp() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceCheckHttpCreate,
		ReadWithoutTimeout:   resourceCheckHttpRead,
		UpdateWithoutTimeout: resourceCheckHttpUpdate,
		DeleteWithoutTimeout: resourceCheckHttpDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Description: "This resource manages HTTP(s) check in RMON.",

		Schema: map[string]*schema.Schema{
			NameField: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Check Http.",
			},
			DescriptionField: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the Check HTTP(s).",
			},
			EnabledField: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enabled state of the Check HTTP(s).",
			},
			CheckGroupIdFiled: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the check group for group HTTP(s) checks.",
			},
			PlaceField: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Port number for binding Check HTTP(s).",
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
			UrlField: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "URL what must be checked.",
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			HttpMethodField: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "HTTP method for HTTP(s) check.",
				ValidateFunc: validation.StringInSlice([]string{
					"get",
					"post",
					"put",
					"delete",
					"options",
					"head",
				}, false),
			},
			IgnoreSslErrorField: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Ignore TLS/SSL error.",
			},
			AcceptedStatusCodesField: {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "Expected status code (default to 200, optional).",
				ValidateFunc: validation.IntBetween(100, 599),
			},
			BodyField: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Check body answer.",
			},
			BodyRequestField: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Send body to server. In JSON.",
				ValidateFunc: validation.StringIsJSON,
			},
			HeaderRequestField: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Send headers to server. In JSON.",
				ValidateFunc: validation.StringIsJSON,
			},
		},
	}
}

func resourceCheckHttpCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client

	description := strings.ReplaceAll(d.Get(DescriptionField).(string), "'", "")
	name := strings.ReplaceAll(d.Get(NameField).(string), "'", "")
	enabled := boolToInt(d.Get(EnabledField).(bool))

	server := map[string]interface{}{
		DescriptionField:         description,
		EnabledField:             enabled,
		NameField:                name,
		CheckGroupIdFiled:        d.Get(CheckGroupIdFiled).(string),
		PlaceField:               d.Get(PlaceField).(string),
		EntitiesField:            d.Get(EntitiesField),
		IntervalField:            d.Get(IntervalField).(int),
		TimeoutField:             d.Get(TimeoutField).(int),
		TelegramField:            d.Get(TelegramField).(int),
		SlackField:               d.Get(SlackField).(int),
		MMField:                  d.Get(MMField).(int),
		PDField:                  d.Get(PDField).(int),
		UrlField:                 d.Get(UrlField),
		HttpMethodField:          d.Get(HttpMethodField),
		IgnoreSslErrorField:      boolToInt(d.Get(IgnoreSslErrorField).(bool)),
		AcceptedStatusCodesField: d.Get(AcceptedStatusCodesField),
		BodyField:                d.Get(BodyField),
		BodyRequestField:         d.Get(BodyRequestField),
		HeaderRequestField:       d.Get(HeaderRequestField),
	}

	resp, err := client.doRequest("POST", "/api/v1.0/rmon/check/http", server)
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
	return resourceCheckHttpRead(ctx, d, m)
}

func resourceCheckHttpRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	resp, err := client.doRequest("GET", fmt.Sprintf("/api/v1.0/rmon/check/http/%s", id), nil)
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
	d.Set(UrlField, result[UrlField])
	d.Set(HttpMethodField, result[HttpMethodField])
	d.Set(IgnoreSslErrorField, intToBool(result[IgnoreSslErrorField].(float64)))
	d.Set(AcceptedStatusCodesField, result[AcceptedStatusCodesField])
	d.Set(BodyField, result[BodyField])
	d.Set(BodyRequestField, result[BodyRequestField])
	d.Set(HeaderRequestField, result[HeaderRequestField])

	return nil
}

func resourceCheckHttpUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	description := strings.ReplaceAll(d.Get(DescriptionField).(string), "'", "")
	name := strings.ReplaceAll(d.Get(NameField).(string), "'", "")
	enabled := boolToInt(d.Get(EnabledField).(bool))

	server := map[string]interface{}{
		DescriptionField:         description,
		EnabledField:             enabled,
		NameField:                name,
		CheckGroupIdFiled:        d.Get(CheckGroupIdFiled).(string),
		PlaceField:               d.Get(PlaceField).(string),
		EntitiesField:            d.Get(EntitiesField),
		IntervalField:            d.Get(IntervalField).(int),
		TimeoutField:             d.Get(TimeoutField).(int),
		TelegramField:            d.Get(TelegramField).(int),
		SlackField:               d.Get(SlackField).(int),
		MMField:                  d.Get(MMField).(int),
		PDField:                  d.Get(PDField).(int),
		UrlField:                 d.Get(UrlField),
		HttpMethodField:          d.Get(HttpMethodField),
		IgnoreSslErrorField:      boolToInt(d.Get(IgnoreSslErrorField).(bool)),
		AcceptedStatusCodesField: d.Get(AcceptedStatusCodesField),
		BodyField:                d.Get(BodyField),
		BodyRequestField:         d.Get(BodyRequestField),
		HeaderRequestField:       d.Get(HeaderRequestField),
	}

	if d.HasChange(PortField) {
		server[ReconfigureField] = true
	}

	_, err := client.doRequest("PUT", fmt.Sprintf("/api/v1.0/rmon/check/http/%s", id), server)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCheckHttpRead(ctx, d, m)
}

func resourceCheckHttpDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	_, err := client.doRequest("DELETE", fmt.Sprintf("/api/v1.0/rmon/check/http/%s", id), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

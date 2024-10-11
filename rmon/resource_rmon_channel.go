package rmon

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	ReceiverField          = "receiver"
	ChannelField           = "channel"
	TokenField             = "token"
	ReceiverTypeTelegram   = "telegram"
	ReceiverTypeSlack      = "slack"
	ReceiverTypePagerDuty  = "pd"
	ReceiverTypeMattermost = "mm"
)

func resourceChannel() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceChannelCreate,
		ReadWithoutTimeout:   resourceChannelRead,
		UpdateWithoutTimeout: resourceChannelUpdate,
		DeleteWithoutTimeout: resourceChannelDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Description: "Represents a communication channel such as Telegram, Slack, PagerDuty, or Mattermost.",

		Schema: map[string]*schema.Schema{
			ReceiverField: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  fmt.Sprintf("The type of the receiver. Only `%s`, `%s`, `%s`, `%s` are allowed.", ReceiverTypeTelegram, ReceiverTypeSlack, ReceiverTypePagerDuty, ReceiverTypeMattermost),
				ValidateFunc: validation.StringInSlice([]string{ReceiverTypeTelegram, ReceiverTypeSlack, ReceiverTypePagerDuty, ReceiverTypeMattermost}, true),
			},
			ChannelField: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The channel identifier.",
			},
			GroupIDField: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the group to which the channel belongs.",
			},
			TokenField: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The token used for the channel.",
			},
		},
	}
}

func resourceChannelCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client

	receiver := d.Get(ReceiverField).(string)

	channel := map[string]interface{}{
		ReceiverField: receiver,
		ChannelField:  strings.ReplaceAll(d.Get(ChannelField).(string), "'", ""),
		GroupIDField:  d.Get(GroupIDField).(int),
		TokenField:    d.Get(TokenField).(string),
	}

	resp, err := client.doRequest("POST", fmt.Sprintf("/api/v1.0/channel/%s", receiver), channel)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("API response: %s", resp)

	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return diag.FromErr(err)
	}

	id, ok := result["id"].(float64)
	if !ok {
		return diag.Errorf("unable to find ID in response: %v", result)
	}

	d.SetId(fmt.Sprintf("%d", int(id)))
	return resourceChannelRead(ctx, d, m)
}

func resourceChannelRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()
	receiver := d.Get(ReceiverField).(string)

	resp, err := client.doRequest("GET", fmt.Sprintf("/api/v1.0/channel/%s/%s", receiver, id), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return diag.FromErr(err)
	}

	if receiverValue, ok := result[ReceiverField].(string); ok && receiverValue != "" {
		d.Set(ReceiverField, receiverValue)
	}

	if channelValue, ok := result[ChannelField].(string); ok && channelValue != "" {
		channel := strings.ReplaceAll(channelValue, "'", "")
		d.Set(ChannelField, channel)
	}

	if groupIDValue, ok := result[GroupIDField].(float64); ok {
		d.Set(GroupIDField, int(groupIDValue))
	}

	if tokenValue, ok := result[TokenField].(string); ok && tokenValue != "" {
		d.Set(TokenField, tokenValue)
	}

	return nil
}

func resourceChannelUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()
	receiver := d.Get(ReceiverField).(string)

	channel := map[string]interface{}{
		ReceiverField: d.Get(ReceiverField).(string),
		ChannelField:  strings.ReplaceAll(d.Get(ChannelField).(string), "'", ""),
		GroupIDField:  d.Get(GroupIDField).(int),
		TokenField:    d.Get(TokenField).(string),
	}

	_, err := client.doRequest("PUT", fmt.Sprintf("/api/v1.0/channel/%s/%s", receiver, id), channel)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceChannelRead(ctx, d, m)
}

func resourceChannelDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()
	receiver := d.Get(ReceiverField).(string)

	_, err := client.doRequest("DELETE", fmt.Sprintf("/api/v1.0/channel/%s/%s", receiver, id), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

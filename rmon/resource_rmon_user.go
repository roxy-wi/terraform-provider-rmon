package rmon

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

const (
	UserEmailField    = "email"
	UserEnabledField  = "enabled"
	UserPasswordField = "password"
	UserUsernameField = "username"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceUserCreate,
		ReadWithoutTimeout:   resourceUserRead,
		UpdateWithoutTimeout: resourceUserUpdate,
		DeleteWithoutTimeout: resourceUserDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Description: "This resource manages user entries in RMON. It allows you to define users with specific email addresses, usernames, passwords, and enabled statuses.",

		Schema: map[string]*schema.Schema{
			UserEmailField: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The email of the user.",
				ValidateFunc: validateEmail,
			},
			UserEnabledField: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the user is enabled (true for enabled, false for disabled).",
			},
			UserPasswordField: {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The password of the user.",
			},
			UserUsernameField: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The username of the user.",
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client

	user := map[string]interface{}{
		UserEmailField:    d.Get(UserEmailField).(string),
		UserEnabledField:  boolToInt(d.Get(UserEnabledField).(bool)),
		UserPasswordField: d.Get(UserPasswordField).(string),
		UserUsernameField: d.Get(UserUsernameField).(string),
	}

	resp, err := client.doRequest("POST", "/api/v1.0/v1.0/user", user)
	if err != nil {
		return diag.FromErr(err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return diag.FromErr(err)
	}

	id, ok := result["id"]
	if !ok {
		return diag.Errorf("unable to extract user ID from response: %v", result)
	}

	switch v := id.(type) {
	case string:
		d.SetId(v)
	case float64:
		d.SetId(fmt.Sprintf("%.0f", v))
	default:
		return diag.Errorf("unexpected type for user ID: %T", id)
	}

	return resourceUserRead(ctx, d, m)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client

	resp, err := client.doRequest("GET", fmt.Sprintf("/api/v1.0/v1.0/user/%s", d.Id()), nil)
	if err != nil {
		if isNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(UserEmailField, result[UserEmailField]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(UserEnabledField, intToBool(result[UserEnabledField].(float64))); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(UserUsernameField, result[UserUsernameField]); err != nil {
		return diag.FromErr(err)
	}
	// Note: Password is not set here for security reasons

	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client

	user := map[string]interface{}{
		UserEmailField:    d.Get(UserEmailField).(string),
		UserEnabledField:  boolToInt(d.Get(UserEnabledField).(bool)),
		UserPasswordField: d.Get(UserPasswordField).(string),
		UserUsernameField: d.Get(UserUsernameField).(string),
	}

	_, err := client.doRequest("PUT", fmt.Sprintf("/api/v1.0/v1.0/user/%s", d.Id()), user)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client

	_, err := client.doRequest("DELETE", fmt.Sprintf("/api/v1.0/v1.0/user/%s", d.Id()), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

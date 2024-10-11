package rmon

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
	"time"
)

const (
	KeyEnabledField = "key_enabled"
	UsernameField   = "username"
	PassPhraseField = "passphrase"
	PrivateKeyField = "private_key"
	SharedField     = "shared"
)

func resourceSSHCredential() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSSHCredentialCreate,
		ReadWithoutTimeout:   resourceSSHCredentialRead,
		UpdateWithoutTimeout: resourceSSHCredentialUpdate,
		DeleteWithoutTimeout: resourceSSHCredentialDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Description: "Manages SSH credentials for RMON.",

		Schema: map[string]*schema.Schema{
			GroupIDField: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Group ID.",
			},
			KeyEnabledField: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Key enabled. `true` you want use private_key instead of password, `false` otherwise.",
			},
			NameField: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the credentials.",
			},
			PasswordField: {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Password for the SSH credentials.",
			},
			UsernameField: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username for the SSH credentials.",
			},
			PassPhraseField: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Passphrase for the SSH credentials.",
			},
			PrivateKeyField: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Private key in Base64 for the SSH credentials. Only ecdsa and rsa is supported.",
			},
			SharedField: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if the credentials are shared.",
			},
		},
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			username := d.Get(UsernameField).(string)
			password := d.Get(PasswordField).(string)
			passphrase := d.Get(PassPhraseField).(string)
			privateKey := d.Get(PrivateKeyField).(string)
			keyEnabled := d.Get(KeyEnabledField).(bool)

			if username != "" && password != "" {
				if passphrase != "" || privateKey != "" {
					return fmt.Errorf("`%s` and `%s` cannot be set when `%s` are provided",
						PassPhraseField, PrivateKeyField, PasswordField)
				}
			}

			if (passphrase != "" || privateKey != "") && !keyEnabled {
				return fmt.Errorf("`%s` must be true when `%s` or `%s` is set", KeyEnabledField, PassPhraseField, PrivateKeyField)
			}

			return nil
		},
	}
}

func resourceSSHCredentialCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client

	sshCred := map[string]interface{}{
		GroupIDField:    d.Get(GroupIDField).(int),
		KeyEnabledField: boolToInt(d.Get(KeyEnabledField).(bool)),
		NameField:       strings.ReplaceAll(d.Get(NameField).(string), "'", ""),
		PasswordField:   d.Get(PasswordField).(string),
		UsernameField:   strings.ReplaceAll(d.Get(UsernameField).(string), "'", ""),
		SharedField:     boolToInt(d.Get(SharedField).(bool)),
	}

	resp, err := client.doRequest("POST", "/api/v1.0/server/cred", sshCred)
	if err != nil {
		return diag.FromErr(err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return diag.Errorf("unexpected response format, could not unmarshal: %s", string(resp))
	}

	id, ok := result["id"].(float64)
	if !ok {
		return diag.Errorf("unable to find ID in response: %v", result)
	}
	d.SetId(fmt.Sprintf("%d", int(id)))

	status, ok := result["status"].(string)
	if !ok || status != "Ok" {
		return diag.Errorf("unexpected status in response: %v", result)
	}

	if d.Get(PassPhraseField).(string) != "" || d.Get(PrivateKeyField).(string) != "" {
		patchData := map[string]interface{}{}
		if d.Get(PassPhraseField).(string) != "" {
			patchData[PassPhraseField] = d.Get(PassPhraseField).(string)
		}
		if d.Get(PrivateKeyField).(string) != "" {
			patchData[PrivateKeyField] = d.Get(PrivateKeyField).(string)
		}

		resp, err := client.doRequest("PATCH", fmt.Sprintf("/api/v1.0/server/cred/%s", d.Id()), patchData)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := json.Unmarshal(resp, &result); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceSSHCredentialRead(ctx, d, m)
}

func resourceSSHCredentialRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	resp, err := client.doRequest("GET", fmt.Sprintf("/api/v1.0/server/cred/%s", id), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var resultArray []map[string]interface{}
	if err := json.Unmarshal(resp, &resultArray); err != nil {
		return diag.Errorf("unexpected response format, could not unmarshal: %s", string(resp))
	}

	if len(resultArray) == 0 {
		return diag.Errorf("empty array in response")
	}
	result := resultArray[0]

	d.Set(GroupIDField, result[GroupIDField])
	d.Set(KeyEnabledField, intToBool(result[KeyEnabledField].(float64)))
	name := strings.ReplaceAll(result[NameField].(string), "'", "")
	d.Set(NameField, name)
	d.Set(PasswordField, result[PasswordField])
	username := strings.ReplaceAll(result[UsernameField].(string), "'", "")
	d.Set(UsernameField, username)
	d.Set(PassPhraseField, result[PassPhraseField])
	d.Set(PrivateKeyField, result[PrivateKeyField])
	d.Set(SharedField, intToBool(result[SharedField].(float64)))

	return nil
}

func resourceSSHCredentialUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	sshCred := map[string]interface{}{
		GroupIDField:    d.Get(GroupIDField).(int),
		KeyEnabledField: boolToInt(d.Get(KeyEnabledField).(bool)),
		NameField:       strings.ReplaceAll(d.Get(NameField).(string), "'", ""),
		PasswordField:   d.Get(PasswordField).(string),
		UsernameField:   strings.ReplaceAll(d.Get(UsernameField).(string), "'", ""),
		SharedField:     boolToInt(d.Get(SharedField).(bool)),
	}

	if d.Get(KeyEnabledField).(bool) {
		privateKey := d.Get(PrivateKeyField).(string)
		if privateKey == "" {
			return diag.Errorf("`%s` must be provided when `%s` is true", PrivateKeyField, KeyEnabledField)
		}
		sshCred[PrivateKeyField] = privateKey
	}

	resp, err := client.doRequest("PUT", fmt.Sprintf("/api/v1.0/server/cred/%s", id), sshCred)
	if err != nil {
		return diag.FromErr(err)
	}

	var statusResponse map[string]interface{}
	if err := json.Unmarshal(resp, &statusResponse); err == nil {
		if status, ok := statusResponse["status"]; ok && status == "Ok" {
			patchData := map[string]interface{}{}
			if d.HasChange(PassPhraseField) {
				patchData[PassPhraseField] = d.Get(PassPhraseField).(string)
			}
			if d.HasChange(PrivateKeyField) {
				patchData[PrivateKeyField] = d.Get(PrivateKeyField).(string)
			}

			if len(patchData) > 0 {
				resp, err := client.doRequest("PATCH", fmt.Sprintf("/api/v1.0/server/cred/%s", d.Id()), patchData)
				if err != nil {
					return diag.FromErr(err)
				}

				if err := json.Unmarshal(resp, &statusResponse); err != nil {
					return diag.FromErr(err)
				}
			}

			return resourceSSHCredentialRead(ctx, d, m)
		}
	}

	var resultArray []map[string]interface{}
	if err := json.Unmarshal(resp, &resultArray); err != nil {
		return diag.Errorf("unexpected response format, could not unmarshal: %s", string(resp))
	}

	if len(resultArray) == 0 {
		return diag.Errorf("empty array in response")
	}

	return resourceSSHCredentialRead(ctx, d, m)
}

func resourceSSHCredentialDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	id := d.Id()

	resp, err := client.doRequest("DELETE", fmt.Sprintf("/api/v1.0/server/cred/%s", id), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(resp) == 0 {
		d.SetId("")
		return nil
	}

	var resultArray []map[string]interface{}
	if err := json.Unmarshal(resp, &resultArray); err != nil {
		return diag.Errorf("unexpected response format, could not unmarshal: %s", string(resp))
	}

	if len(resultArray) > 0 {
		d.SetId("")
		return nil
	}

	return diag.Errorf("unexpected response format during deletion")
}

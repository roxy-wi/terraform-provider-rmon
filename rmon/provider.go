package rmon

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Config struct {
	Client *Client
}

const (
	ProviderBaseURL = "base_url"
	LoginField      = "login"
	PasswordField   = "password"
)

func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			LoginField: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username for RMON.",
				DefaultFunc: schema.EnvDefaultFunc("RMON_USERNAME", nil),
			},
			PasswordField: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Password for RMON.",
				DefaultFunc: schema.EnvDefaultFunc("RMON_PASSWORD", nil),
			},
			ProviderBaseURL: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "URL to connect for RMON.",
				DefaultFunc: schema.EnvDefaultFunc("RMON_BASE_URL", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"rmon_group":             resourceGroup(),
			"rmon_user":              resourceUser(),
			"rmon_user_role_binding": resourceUserRoleBinding(),
			"rmon_server":            resourceServer(),
			"rmon_channel":           resourceChannel(),
			"rmon_ssh_credential":    resourceSSHCredential(),
			"rmon_agent":             resourceAgent(),
			"rmon_region":            resourceRegion(),
			"rmon_country":           resourceCountry(),
			"rmon_check_ping":        resourceCheckPing(),
			"rmon_check_tcp":         resourceCheckTcp(),
			"rmon_check_dns":         resourceCheckDns(),
			"rmon_check_http":        resourceCheckHttp(),
			"rmon_check_smtp":        resourceCheckSmtp(),
			"rmon_check_rabbitmq":    resourceCheckRabbitmq(),
			"rmon_check_group":       resourceCheckGroup(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"rmon_group":     dataSourceGroup(),
			"rmon_user_role": dataSourceUserRole(),
		},
	}

	p.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			terraformVersion = "1.0+compatible"
		}
		return providerConfigure(ctx, d, terraformVersion)
	}

	return p
}

func providerConfigure(
	_ context.Context,
	d *schema.ResourceData,
	terraformVersion string,
) (interface{}, diag.Diagnostics) {
	username := d.Get(LoginField).(string)
	password := d.Get(PasswordField).(string)
	apiEndpoint := d.Get(ProviderBaseURL).(string)

	userAgent := fmt.Sprintf("terraform/%s", terraformVersion)

	var diags diag.Diagnostics

	client, err := NewClient(apiEndpoint, username, password, userAgent)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	config := &Config{
		Client: client,
	}

	return config, diags
}

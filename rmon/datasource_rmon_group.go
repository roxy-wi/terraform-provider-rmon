package rmon

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	IDField          = "id"
	NameField        = "name"
	DescriptionField = "description"
)

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupRead,
		Description: "Represent roxy-wi group. All servers managed via RMON can be included in groups. A group is a user-defined pool of servers. By default, all your servers are included in the common group named Default.",

		Schema: map[string]*schema.Schema{
			IDField: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{IDField, NameField},
				Description:  "ID of the group.",
			},
			NameField: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{IDField, NameField},
				Description:  "The name of the group.",
			},
			DescriptionField: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the group.",
			},
		},
	}
}

func dataSourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client
	if id, ok := d.GetOk(IDField); ok {
		return readGroupByID(ctx, d, client, id.(string))
	} else if name, ok := d.GetOk(NameField); ok {
		return readGroupByName(ctx, d, client, name.(string))
	}

	return diag.Errorf("either 'id' or 'name' must be specified")
}

func readGroupByID(ctx context.Context, d *schema.ResourceData, client *Client, id string) diag.Diagnostics {
	resp, err := client.doRequest("GET", fmt.Sprintf("/api/v1.0/group/%s", id), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	// Process response and set data
	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return diag.FromErr(err)
	}

	if name, ok := result[NameField].(string); ok {
		d.Set(NameField, name)
	}

	if description, ok := result[DescriptionField].(string); ok {
		d.Set(DescriptionField, description)
	}

	d.SetId(id)
	return nil
}

func readGroupByName(ctx context.Context, d *schema.ResourceData, client *Client, name string) diag.Diagnostics {
	resp, err := client.doRequest("GET", "/api/v1.0/groups", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	// Process response and find the group by name
	var groups []map[string]interface{}
	if err := json.Unmarshal(resp, &groups); err != nil {
		return diag.FromErr(err)
	}

	for _, group := range groups {
		if groupName, ok := group[NameField].(string); ok && groupName == name {
			if id, ok := group["group_id"].(float64); ok {
				d.SetId(fmt.Sprintf("%d", int(id)))
				d.Set(NameField, groupName)

				if description, ok := group[DescriptionField].(string); ok {
					d.Set(DescriptionField, description)
				}

				return nil
			}
		}
	}

	return diag.Errorf("group with name '%s' not found", name)
}

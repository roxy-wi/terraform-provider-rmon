package rmon

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	UserIDField = "user_id"
)

func resourceUserRoleBinding() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceUserRoleBindingCreate,
		ReadWithoutTimeout:   resourceUserRoleBindingRead,
		UpdateWithoutTimeout: resourceUserRoleBindingUpdate,
		DeleteWithoutTimeout: resourceUserRoleBindingDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Description: "The resource manages the binding of a user to a specific role within a group in RMON. This resource allows you to assign roles to users, facilitating role-based access control and management.",

		Schema: map[string]*schema.Schema{
			UserIDField: {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the user.",
			},
			RoleIDField: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the role.",
			},
			GroupIDField: {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the group.",
			},
		},
	}
}

func resourceUserRoleBindingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client

	userID := d.Get(UserIDField).(int)
	groupID := d.Get(GroupIDField).(int)

	binding := map[string]interface{}{
		RoleIDField: d.Get(RoleIDField).(int),
	}

	resp, err := client.doRequest("POST", fmt.Sprintf("/api/v1.0/user/%d/groups/%d", userID, groupID), binding)
	if err != nil {
		return diag.FromErr(err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d-%d", userID, groupID))

	return resourceUserRoleBindingRead(ctx, d, m)
}

func resourceUserRoleBindingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client

	id := d.Id()
	ids := strings.Split(id, "-")
	if len(ids) != 2 {
		return diag.Errorf("invalid ID format for user role binding: %s", id)
	}
	userIDStr := ids[0]
	groupIDStr := ids[1]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return diag.FromErr(err)
	}
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.doRequest("GET", fmt.Sprintf("/api/v1.0/user/%d/groups", userID), nil)
	if err != nil {
		if isNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return diag.FromErr(err)
	}

	var found bool
	for _, group := range result {
		if gid, ok := group["user_group_id"].(float64); ok && int(gid) == groupID {
			if err := d.Set(UserIDField, userID); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set(RoleIDField, group["user_role_id"]); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set(GroupIDField, groupID); err != nil {
				return diag.FromErr(err)
			}
			found = true
			break
		}
	}

	if !found {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceUserRoleBindingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client

	userID := d.Get(UserIDField).(int)
	groupID := d.Get(GroupIDField).(int)

	binding := map[string]interface{}{
		RoleIDField: d.Get(RoleIDField).(int),
	}

	_, err := client.doRequest("PUT", fmt.Sprintf("/api/v1.0/user/%d/groups/%d", userID, groupID), binding)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceUserRoleBindingRead(ctx, d, m)
}

func resourceUserRoleBindingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client

	id := d.Id()
	ids := strings.Split(id, "-")
	if len(ids) != 2 {
		return diag.Errorf("invalid ID format for user role binding: %s", id)
	}
	userIDStr := ids[0]
	groupIDStr := ids[1]

	_, err := client.doRequest("DELETE", fmt.Sprintf("/api/v1.0/user/%s/groups/%s", userIDStr, groupIDStr), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

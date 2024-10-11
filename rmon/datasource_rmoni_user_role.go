package rmon

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	ctyjson "github.com/hashicorp/go-cty/cty/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	RoleIDField          = "role_id"
	RoleNameField        = "name"
	RoleDescriptionField = "description"
	RolesField           = "roles"
)

func dataSourceUserRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRoleRead,

		Description: "The data source allows you to retrieve information about user roles in RMON. This data source fetches a list of roles, providing details about each role's ID, name, and description.",

		Schema: map[string]*schema.Schema{
			RolesField: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of roles.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						RoleIDField: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the role.",
						},
						RoleNameField: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the role.",
						},
						RoleDescriptionField: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the role.",
						},
					},
				},
			},
		},
	}
}

func dataSourceUserRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config).Client

	resp, err := client.doRequest("GET", "/api/v1.0user/roles", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	rolesType := cty.List(cty.Object(map[string]cty.Type{
		RoleIDField:          cty.Number,
		RoleNameField:        cty.String,
		RoleDescriptionField: cty.String,
	}))

	rolesVal, err := ctyjson.Unmarshal(resp, rolesType)
	if err != nil {
		return diag.FromErr(err)
	}

	if rolesVal.IsNull() || rolesVal.LengthInt() == 0 {
		return diag.Errorf("No roles found")
	}

	convertedRoles, err := convertRoles(rolesVal.AsValueSlice())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(RolesField, convertedRoles); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("roles")
	return nil
}

func convertRoles(roles []cty.Value) ([]map[string]interface{}, error) {
	var convertedRoles []map[string]interface{}
	for _, roleVal := range roles {
		role := roleVal.AsValueMap()
		convertedRole := make(map[string]interface{})

		roleID, err := extractRoleField(role, RoleIDField)
		if err != nil {
			return nil, err
		}
		convertedRole[RoleIDField] = roleID

		roleName, err := extractRoleField(role, RoleNameField)
		if err != nil {
			return nil, err
		}
		convertedRole[RoleNameField] = roleName

		roleDescription, err := extractRoleField(role, RoleDescriptionField)
		if err != nil {
			return nil, err
		}
		convertedRole[RoleDescriptionField] = roleDescription

		convertedRoles = append(convertedRoles, convertedRole)
	}
	return convertedRoles, nil
}

func extractRoleField(role map[string]cty.Value, field string) (string, error) {
	if val, ok := role[field]; ok && !val.IsNull() {
		switch field {
		case RoleIDField:
			int64Value, _ := val.AsBigFloat().Int64()
			return fmt.Sprintf("%d", int64(int64Value)), nil
		default:
			return val.AsString(), nil
		}
	}
	return "", fmt.Errorf("unexpected type or null value for %s field", field)
}

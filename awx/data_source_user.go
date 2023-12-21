/*
Use this resource to query a user in AWX/AT with specified role entitlements

# Example Usage
```hcl

	data "awx_user" "my_user" {
	    username = "My_User"
	}

	output "my_user" {
	    value = data.awx_user.my_user
	}

```
*/
package awx

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/sharathrnair87/goawx/client"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to query a user in AWX/AT with specified role entitlements",
		ReadContext: dataSourceUserRead,

		Schema: map[string]*schema.Schema{
			"username": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id"},
				ExactlyOneOf:  []string{"username", "id"},
			},
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"username"},
				ExactlyOneOf:  []string{"username", "id"},
			},
			"first_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_superuser": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_system_auditor": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"role_entitlement": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Set of role IDs of the role entitlements",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"resource_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*awx.AWX)
	awxService := client.UserService
	params := make(map[string]string)

	if userName, okName := d.GetOk("username"); okName {
		params["username"] = userName.(string)
	}

	if userID, okID := d.GetOk("id"); okID {
		params["id"] = strconv.Itoa(userID.(int))
	}

	if len(params) == 0 {
		return buildDiagnosticsMessage(
			"Get: Missing Parameters",
			"Please use one of the selectors (username or id)",
		)
	}

	users, _, err := awxService.ListUsers(params)
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Failed to fetch Team",
			"Failed to find the team got: %s",
			err.Error(),
		)
	}
	if len(users) > 1 {
		return buildDiagnosticsMessage(
			"Get: found more than one Element",
			"The Query Returns more than one user, %d",
			len(users),
		)
	}
	if len(users) == 0 {
		return buildDiagnosticsMessage(
			"Get: No User found",
			"The Query Returns no user matching filter, %v",
			len(users),
		)
	}

	user := users[0]

	entitlements, _, err := awxService.ListUserRoleEntitlements(user.ID, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail("user roles", user.ID, err)
	}

	d.SetId(strconv.Itoa(user.ID))
	d.Set("username", user.Username)
	d.Set("first_name", user.FirstName)
	d.Set("last_name", user.LastName)
	d.Set("email", user.Email)
	d.Set("is_superuser", user.IsSuperUser)
	d.Set("is_system_auditor", user.IsSystemAuditor)

	var entlist []interface{}
	for _, v := range entitlements {
		elem := make(map[string]interface{})
		elem["role_id"] = v.ID
		elem["resource_name"] = v.Summary.ResourceName
		elem["resource_type"] = v.Summary.ResourceType
		entlist = append(entlist, elem)
	}
	f := schema.HashResource(&schema.Resource{
		Schema: map[string]*schema.Schema{
			"role_id":       {Type: schema.TypeInt},
			"resource_name": {Type: schema.TypeString},
			"resource_type": {Type: schema.TypeString},
		}})

	ent := schema.NewSet(f, entlist)

	d.Set("role_entitlement", ent)

	return diags
}

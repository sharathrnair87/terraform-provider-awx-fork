/*
Use this data source to lookup a Credential Role in AWX/AT

# Example Usage

```hcl

	resource "awx_credential" "my_credential" {
	  name = "My Credential"
	}

	data "awx_credential_role" "credential_admin_role" {
	  role_name     = "Admin"
	  credential_id = data.awx_credential.my_credential.id
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

func dataSourceCredentialRole() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to lookup a Credential Role in AWX/AT",
		ReadContext: dataSourceCredentialRoleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"id"},
			},
			"credential_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func dataSourceCredentialRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)

	cred_id := d.Get("credential_id").(int)
	cred, err := client.CredentialsService.GetCredentialsByID(cred_id, params)
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Failed to fetch credential",
			"Failed to find the credential, got: %s",
			err.Error(),
		)
	}

	roleslist := []*awx.ApplyRole{
		cred.SummaryFields.ObjectRoles.AdminRole,
		cred.SummaryFields.ObjectRoles.ReadRole,
		cred.SummaryFields.ObjectRoles.ExecuteRole,
	}

	if roleID, okID := d.GetOk("id"); okID {
		id := roleID.(int)
		for _, v := range roleslist {
			if v != nil && id == v.ID {
				d = setCredentialRoleData(d, v)
				return diags
			}
		}
	}

	if roleName, okName := d.GetOk("name"); okName {
		name := roleName.(string)

		for _, v := range roleslist {
			if v != nil && name == v.Name {
				d = setCredentialRoleData(d, v)
				return diags
			}
		}
	}

	return buildDiagnosticsMessage(
		"Failed to fetch credential role - Not Found",
		"The project role was not found",
	)
}

func setCredentialRoleData(d *schema.ResourceData, r *awx.ApplyRole) *schema.ResourceData {
	d.Set("name", r.Name)
	d.SetId(strconv.Itoa(r.ID))
	return d
}

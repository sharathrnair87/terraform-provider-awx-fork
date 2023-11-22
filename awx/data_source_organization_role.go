/*
Use this data source to query the ID of a given Organization Role in AWX/AT

# Example Usage

```hcl

	resource "awx_organization" "myorg" {
	  name = "My AWX Org"
	  // Truncated //
	}

	data "awx_organization_role" "org_admins" {
	  name            = "Admin"
	  organization_id = resource.awx_organization.myorg.id
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

func dataSourceOrganizationRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrganizationRolesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"organization_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func dataSourceOrganizationRolesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)

	org_id := d.Get("organization_id").(int)

	organization, err := client.OrganizationsService.GetOrganizationsByID(org_id, params)
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Failed to fetch organization role",
			"Failed to find the organization role got: %s",
			err.Error(),
		)
	}

	roleslist := []*awx.ApplyRole{
		organization.SummaryFields.ObjectRoles.AdhocRole,
		organization.SummaryFields.ObjectRoles.AdminRole,
		organization.SummaryFields.ObjectRoles.ApprovalRole,
		organization.SummaryFields.ObjectRoles.AuditorRole,
		organization.SummaryFields.ObjectRoles.CredentialAdminRole,
		organization.SummaryFields.ObjectRoles.ExecuteRole,
		organization.SummaryFields.ObjectRoles.InventoryAdminRole,
		organization.SummaryFields.ObjectRoles.JobTemplateAdminRole,
		organization.SummaryFields.ObjectRoles.MemberRole,
		organization.SummaryFields.ObjectRoles.NotificationAdminRole,
		organization.SummaryFields.ObjectRoles.ProjectAdminRole,
		organization.SummaryFields.ObjectRoles.ReadRole,
		organization.SummaryFields.ObjectRoles.UpdateRole,
		organization.SummaryFields.ObjectRoles.UseRole,
		organization.SummaryFields.ObjectRoles.WorkflowAdminRole,
	}

	if roleID, okID := d.GetOk("id"); okID {
		id := roleID.(int)
		for _, v := range roleslist {
			if v != nil && id == v.ID {
				d = setOrganizationRoleData(d, v)
				return diags
			}
		}
	}

	if roleName, okName := d.GetOk("name"); okName {
		name := roleName.(string)

		for _, v := range roleslist {
			if v != nil && name == v.Name {
				d = setOrganizationRoleData(d, v)
				return diags
			}
		}
	}

	return buildDiagnosticsMessage(
		"Failed to fetch organization role - Not Found",
		"The organization role was not found",
	)
}

func setOrganizationRoleData(d *schema.ResourceData, r *awx.ApplyRole) *schema.ResourceData {
	d.Set("name", r.Name)
	d.SetId(strconv.Itoa(r.ID))
	return d
}

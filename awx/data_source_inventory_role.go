/*
Use this data source to query the ID of a given Inventory Role in AWX/AT

# Example Usage

```hcl

		resource "awx_inventory" "myinv" {
		  name = "My Inventory"
          // Truncated //
		}

		data "awx_inventory_role" "inv_admin_role" {
		  name         = "Admin"
		  inventory_id = data.awx_inventory.myinv.id
		}

	    output "inv_admin_role_id" {
	        value = data.awx_inventory_role.inv_admin_role.id
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

func dataSourceInventoryRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInventoryRoleRead,
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
			"inventory_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func dataSourceInventoryRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)

	inv_id := d.Get("inventory_id").(int)
	inventory, err := client.InventoriesService.GetInventoryByID(inv_id, params)
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Failed to fetch Inventory",
			"Failed to find the inventory, got: %s",
			err.Error(),
		)
	}

	roleslist := []*awx.ApplyRole{
		inventory.SummaryFields.ObjectRoles.UseRole,
		inventory.SummaryFields.ObjectRoles.AdminRole,
		inventory.SummaryFields.ObjectRoles.AdhocRole,
		inventory.SummaryFields.ObjectRoles.UpdateRole,
		inventory.SummaryFields.ObjectRoles.ReadRole,
		inventory.SummaryFields.ObjectRoles.ExecuteRole,
	}

	if roleID, okID := d.GetOk("id"); okID {
		id := roleID.(int)
		for _, v := range roleslist {
			if v != nil && id == v.ID {
				d = setInventoryRoleData(d, v)
				return diags
			}
		}
	}

	if roleName, okName := d.GetOk("name"); okName {
		name := roleName.(string)

		for _, v := range roleslist {
			if v != nil && name == v.Name {
				d = setInventoryRoleData(d, v)
				return diags
			}
		}
	}

	return buildDiagnosticsMessage(
		"Failed to fetch inventory role - Not Found",
		"The project role was not found",
	)
}

func setInventoryRoleData(d *schema.ResourceData, r *awx.ApplyRole) *schema.ResourceData {
	d.Set("name", r.Name)
	d.SetId(strconv.Itoa(r.ID))
	return d
}

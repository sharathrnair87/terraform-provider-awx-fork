/*
Use this dataSource to query an Inventory Source for an existing Inventory in AWX/AT

# Example Usage

```hcl

	data "awx_inventory" "db_inventory" {
	  name            = "DB_Inventory"
	  organization_id = data.awx_organization.db.id
	}

	data "awx_inventory_source" "db_inventory_source" {
	  name         = "DB_AZ_Inventory"
	  inventory_id = data.awx_inventory.db_inventory.id
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

func dataSourceInventorySource() *schema.Resource {
	return &schema.Resource{
		Description: "Use this dataSource to query an Inventory Source for an existing Inventory in AWX/AT",
		ReadContext: dataSourceInventorySourceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled_var": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"overwrite": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"overwrite_vars": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"update_on_launch": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"inventory_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"credential_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_vars": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_filter": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_cache_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"verbosity": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			// obsolete schema added so terraform doesn't break
			// these don't do anything in later versions of AWX! Update your code.
			"source_regions": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_filters": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_project_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"source_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceInventorySourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.InventorySourcesService

	res, _, err := awxService.ListInventorySources(map[string]string{
		"name":         d.Get("name").(string),
		"inventory_id": strconv.Itoa(d.Get("inventory_id").(int)),
	})

	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Failed to fetch Inventory Source",
			"Failed to find the Inventory Source got: %s",
			err.Error(),
		)
	}
	if len(res) > 1 {
		return buildDiagnosticsMessage(
			"Get: found more than one Inventory Source",
			"The Query Returns more than one Inventory Source, %d",
			len(res),
		)
	}
	result := res[0]

	d.SetId(strconv.Itoa(result.ID))
	d = setInventorySourceResourceData(d, result)
	return nil
}

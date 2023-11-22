/*
Use this dataSource to query a managed host in AWX/AT.

# Example Usage

```hcl

	data "awx_inventory" "db_inventory" {
	  name            = "DB Inventory"
	  organization_id = data.awx_organization.default.id
	}

	data "awx_host" "db_host" {
	  name         = "prddbsrvr01.example.com"
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

func dataSourceHost() *schema.Resource {
	return &schema.Resource{
		Description: "Use this dataSource to query a managed host in AWX/AT.",
		ReadContext: dataSourceHostRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"inventory_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"group_ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"variables": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceHostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.HostService

	res, _, err := awxService.ListHosts(map[string]string{
		"name":         d.Get("name").(string),
		"inventory_id": strconv.Itoa(d.Get("inventory_id").(int)),
	})

	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Failed to fetch Host",
			"Failed to find the Host got: %s",
			err.Error(),
		)
	}
	if len(res) > 1 {
		return buildDiagnosticsMessage(
			"Get: found more than one Element",
			"The Query Returns more than one Host, %d",
			len(res),
		)
	}
	result := res[0]

	d.SetId(strconv.Itoa(result.ID))
	d = setHostResourceData(d, result)
	return nil
}

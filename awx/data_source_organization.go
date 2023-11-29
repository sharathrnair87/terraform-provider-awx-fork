/*
Use this data source to query an AWX/AT Organization

# Example Usage

```hcl

	data "awx_organization" "default" {
	  name = "Default"
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

func dataSourceOrganization() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to query an AWX/AT Organization",
		ReadContext: dataSourceOrganizationRead,
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
			"max_hosts": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"custom_virtualenv": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_environment": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceOrganizationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)
	if groupName, okName := d.GetOk("name"); okName {
		params["name"] = groupName.(string)
	}

	if groupID, okGroupID := d.GetOk("id"); okGroupID {
		params["id"] = strconv.Itoa(groupID.(int))
	}

	if len(params) == 0 {
		return buildDiagnosticsMessage(
			"Get: Missing Parameters",
			"Please use one of the selectors (name or group_id)",
		)
	}
	organizations, err := client.OrganizationsService.ListOrganizations(params)
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Failed to fetch organization",
			"Failed to find the organization got: %s",
			err.Error(),
		)
	}
	if len(organizations) > 1 {
		return buildDiagnosticsMessage(
			"Get: found more than one Element",
			"The Query Returns more than one organization, %d",
			len(organizations),
		)
	}
	if len(organizations) == 0 {
		return buildDiagnosticsMessage(
			"Get: Organization does not exist",
			"The Query Returns no Organization matching filter, %v",
			params,
		)
	}

	organization := organizations[0]
	d = setOrganizationsResourceData(d, organization)
	return diags
}

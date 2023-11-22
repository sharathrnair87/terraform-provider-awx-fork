/*
Use this data source to lookup a Team by name or ID in AWX/AT

# Example Usage

```hcl

	data "awx_team" "default" {
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

func dataSourceTeam() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to lookup a Team by name or ID in AWX/AT",
		ReadContext: dataSourceTeamsRead,
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
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organization_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"role_entitlement": {
				Type:     schema.TypeSet,
				Computed: true,
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
	}
}

func dataSourceTeamsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)
	if teamName, okName := d.GetOk("name"); okName {
		params["name"] = teamName.(string)
	}

	if teamID, okTeamID := d.GetOk("id"); okTeamID {
		params["id"] = strconv.Itoa(teamID.(int))
	}

	if len(params) == 0 {
		return buildDiagnosticsMessage(
			"Get: Missing Parameters",
			"Please use one of the selectors (name or id)",
		)
	}
	Teams, _, err := client.TeamService.ListTeams(params)
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Failed to fetch Team",
			"Failed to find the team got: %s",
			err.Error(),
		)
	}
	if len(Teams) > 1 {
		return buildDiagnosticsMessage(
			"Get: found more than one Element",
			"The Query Returns more than one team, %d",
			len(Teams),
		)
	}

	Team := Teams[0]
	Entitlements, _, err := client.TeamService.ListTeamRoleEntitlements(Team.ID, make(map[string]string))
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Failed to fetch team role entitlements",
			"Failed to retrieve team role entitlements got: %s",
			err.Error(),
		)
	}

	d = setTeamResourceData(d, Team, Entitlements)
	return diags
}

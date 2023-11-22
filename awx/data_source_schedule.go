/*
Use this resource to query a Job Template schedule in AWX/AT

# Example Usage

```hcl

	data "awx_schedule" "default" {
	  name            = "private_services"
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

func dataSourceSchedule() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to query a Job Template schedule in AWX/AT",
		ReadContext: dataSourceSchedulesRead,
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
			"rrule": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"unified_job_template_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"inventory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"extra_data": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSchedulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	params := make(map[string]string)
	if groupName, okName := d.GetOk("name"); okName {
		params["name"] = groupName.(string)
	}

	if groupID, okID := d.GetOk("id"); okID {
		params["id"] = strconv.Itoa(groupID.(int))
	}

	if len(params) == 0 {
		return buildDiagnosticsMessage(
			"Get: Missing Parameters",
			"Please use one of the selectors (name or id)",
		)
	}

	schedules, _, err := client.ScheduleService.List(params)
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Failed to fetch Schedule Group",
			"Failed to find the group got: %s",
			err.Error(),
		)
	}
	if len(schedules) > 1 {
		return buildDiagnosticsMessage(
			"Get: found more than one Element",
			"The Query Returns more than one Group, %d",
			len(schedules),
		)
	}

	schedule := schedules[0]
	d = setScheduleResourceData(d, schedule)
	return diags
}

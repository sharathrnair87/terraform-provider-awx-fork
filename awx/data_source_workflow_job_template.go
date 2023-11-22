/*
Use this resource to lookup a Workflow Job Template in AWX/AT

# Example Usage

```hcl

	        // Lookup by `name`
			data "awx_workflow_job_template" "default" {
		        name = "Default"
			}

	        // Lookup by `id`
			data "awx_workflow_job_template" "default" {
		        id = var.default_workflow_job_template_id
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

func dataSourceWorkflowJobTemplate() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to lookup a Workflow Job Template in AWX/AT",
		ReadContext: dataSourceWorkflowJobTemplateRead,
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
			"variables": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organization_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"survey_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"allow_simultaneous": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ask_variables_on_launch": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"inventory_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"limit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scm_branch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ask_inventory_on_launch": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ask_scm_branch_on_launch": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ask_limit_on_launch": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"webhook_service": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"webhook_credential": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceWorkflowJobTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	workflowJobTemplate, _, err := client.WorkflowJobTemplateService.ListWorkflowJobTemplates(params)
	if err != nil {
		return buildDiagnosticsMessage(
			"Get: Failed to fetch Inventory Group",
			"Failed to find the group got: %s",
			err.Error(),
		)
	}
	if groupName, okName := d.GetOk("name"); okName {
		for _, template := range workflowJobTemplate {
			if template.Name == groupName {
				d = setWorkflowJobTemplateResourceData(d, template)
				return diags
			}
		}
	}
	if _, okGroupID := d.GetOk("id"); okGroupID {
		if len(workflowJobTemplate) != 1 {
			return buildDiagnosticsMessage(
				"Get: found more than one Element",
				"The Query Returns more than one Group, %d",
				len(workflowJobTemplate),
			)
		}
		d = setWorkflowJobTemplateResourceData(d, workflowJobTemplate[0])
		return diags
	}
	return buildDiagnosticsMessage(
		"Get: found more than one Element",
		"The Query Returns more than one Group, %d",
		len(workflowJobTemplate),
	)
}

/*
Use this resource to create an AWX/AT Schedule for a Workflow Job Template

# Example Usage

```hcl

	resource "awx_workflow_job_template_schedule" "default" {
	  workflow_job_template_id      = awx_workflow_job_template.default.id

	  name                      = "schedule-test"
	  rrule                     = "DTSTART;TZID=Europe/Paris:20211214T120000 RRULE:INTERVAL=1;FREQ=DAILY"
	  extra_data                = <<EOL

organization_name: testorg
EOL
}
```
*/
package awx

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/sharathrnair87/goawx/client"
)

func resourceWorkflowJobTemplateSchedule() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create an AWX/AT Schedule for a Workflow Job Template",
		CreateContext: resourceWorkflowJobTemplateScheduleCreate,
		ReadContext:   resourceScheduleRead,
		UpdateContext: resourceScheduleUpdate,
		DeleteContext: resourceScheduleDelete,
		Schema: map[string]*schema.Schema{

			"workflow_job_template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The workflow_job_template id for this schedule",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rrule": {
				Type:     schema.TypeString,
				Required: true,
			},
			"unified_job_template_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"inventory": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Inventory applied as a prompt, assuming job template prompts for inventory (id, default=``)",
			},
			"extra_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Extra data to be pass for the schedule (YAML format)",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceWorkflowJobTemplateScheduleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.WorkflowJobTemplateScheduleService

	workflowJobTemplateID := d.Get("workflow_job_template_id").(int)

	result, err := awxService.CreateWorkflowJobTemplateSchedule(workflowJobTemplateID, map[string]interface{}{
		"name":        d.Get("name").(string),
		"rrule":       d.Get("rrule").(string),
		"description": d.Get("description").(string),
		"enabled":     d.Get("enabled").(bool),
		"inventory":   AtoipOr(d.Get("inventory").(string), nil),
		"extra_data":  unmarshalYaml(d.Get("extra_data").(string)),
	}, map[string]string{})
	if err != nil {
		log.Printf("Failed to Create Schedule for WorkflowJobTemplate %d: %v", workflowJobTemplateID, err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Schedule",
			Detail:   fmt.Sprintf("Schedule failed to create %s", err.Error()),
		})
		return diags
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceScheduleRead(ctx, d, m)
}

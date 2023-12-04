/*
Use this resource to associate a Notification template with a Workflow Job Template, triggered in case of a successful execution

# Example Usage

```hcl

	resource "awx_workflow_job_template_notification_template_success" "baseconfig" {
	  workflow_job_template_id   = awx_workflow_job_template.baseconfig.id
	  notification_template_id   = awx_notification_template.default.id
	}

```
*/
package awx

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceWorkflowJobTemplateNotificationTemplateSuccess() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to associate a Notification template with a Workflow Job Template, triggered in case of a successful execution",
		CreateContext: resourceWorkflowJobTemplateNotificationTemplateCreateForType("success"),
		DeleteContext: resourceWorkflowJobTemplateNotificationTemplateDeleteForType("success"),
		ReadContext:   resourceWorkflowJobTemplateNotificationTemplateRead,

		Schema: map[string]*schema.Schema{
			"workflow_job_template_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"notification_template_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

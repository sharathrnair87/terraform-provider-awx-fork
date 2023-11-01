/*
*TBD*

# Example Usage

```hcl
resource "random_uuid" "workflow_node_k3s_uuid" {}

	resource "awx_workflow_job_template_node_success" "k3s" {
	  workflow_job_template_id      = awx_workflow_job_template.default.id
	  workflow_job_template_node_id = awx_workflow_job_template_node.default.id
	  unified_job_template_id       = awx_job_template.k3s.id
	  inventory_id                  = awx_inventory.default.id
	  identifier                    = random_uuid.workflow_node_k3s_uuid.result
	}

```
*/
package awx

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/sharathrnair87/goawx/client"
)

func resourceWorkflowJobTemplateNodeSuccess() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkflowJobTemplateNodeSuccessCreate,
		ReadContext:   resourceWorkflowJobTemplateNodeRead,
		UpdateContext: resourceWorkflowJobTemplateNodeUpdate,
		DeleteContext: resourceWorkflowJobTemplateNodeDelete,
		Schema:        workflowJobNodeSchema,
	}
}

func resourceWorkflowJobTemplateNodeSuccessCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*awx.AWX)
	awxService := client.WorkflowJobTemplateNodeSuccessService
	return createNodeForWorkflowJob(awxService, ctx, d, m)
}

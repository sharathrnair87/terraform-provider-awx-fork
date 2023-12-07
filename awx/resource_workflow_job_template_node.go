/*
Use this resource to create a node in a workflow job template, for more details see [Workflow Visualizer](https://docs.ansible.com/automation-controller/latest/html/userguide/workflow_templates.html#build-a-workflow)

# Example Usage

```hcl
resource "random_uuid" "workflow_node_base_uuid" {}

	resource "awx_workflow_job_template_node" "default" {
	  workflow_job_template_id = awx_workflow_job_template.default.id
	  unified_job_template_id  = awx_job_template.baseconfig.id
	  inventory_id             = awx_inventory.default.id
	  identifier               = random_uuid.workflow_node_base_uuid.result
	}

```
*/
package awx

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awx "github.com/sharathrnair87/goawx/client"
)

func resourceWorkflowJobTemplateNode() *schema.Resource {
	return &schema.Resource{
		Description:   "Use this resource to create a node in a workflow job template, for more details see [Workflow Visualizer](https://docs.ansible.com/automation-controller/latest/html/userguide/workflow_templates.html#build-a-workflow)",
		CreateContext: resourceWorkflowJobTemplateNodeCreate,
		ReadContext:   resourceWorkflowJobTemplateNodeRead,
		UpdateContext: resourceWorkflowJobTemplateNodeUpdate,
		DeleteContext: resourceWorkflowJobTemplateNodeDelete,

		Schema: map[string]*schema.Schema{

			"extra_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "",
				StateFunc:   normalizeJsonYaml,
			},
			"inventory_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Inventory applied as a prompt, assuming job template prompts for inventory.",
			},
			"scm_branch": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"job_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "run",
			},
			"job_tags": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"skip_tags": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"diff_mode": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"verbosity": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"workflow_job_template_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"unified_job_template_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"all_parents_must_converge": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"identifier": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
	}
}

func resourceWorkflowJobTemplateNodeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.WorkflowJobTemplateNodeService

	result, err := awxService.CreateWorkflowJobTemplateNode(map[string]interface{}{
		"extra_data":                d.Get("extra_data").(string),
		"inventory":                 d.Get("inventory_id").(int),
		"scm_branch":                d.Get("scm_branch").(string),
		"skip_tags":                 d.Get("skip_tags").(string),
		"job_type":                  d.Get("job_type").(string),
		"job_tags":                  d.Get("job_tags").(string),
		"limit":                     d.Get("limit").(string),
		"diff_mode":                 d.Get("diff_mode").(bool),
		"verbosity":                 d.Get("verbosity").(int),
		"workflow_job_template":     d.Get("workflow_job_template_id").(int),
		"unified_job_template":      d.Get("unified_job_template_id").(int),
		"all_parents_must_converge": d.Get("all_parents_must_converge").(bool),
		"identifier":                d.Get("identifier").(string),
	}, map[string]string{})
	if err != nil {
		log.Printf("Failed to Create Template %v", err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create WorkflowJobTemplateNode",
			Detail:   fmt.Sprintf("WorkflowJobTemplateNode with JobTemplateID %d and WorkflowID: %d failed to create %s", d.Get("unified_job_template_id").(int), d.Get("workflow_job_template_id").(int), err.Error()),
		})
		return diags
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceWorkflowJobTemplateNodeRead(ctx, d, m)
}

func resourceWorkflowJobTemplateNodeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.WorkflowJobTemplateNodeService
	id, diags := convertStateIDToNumeric("Update WorkflowJobTemplateNode", d)
	if diags.HasError() {
		return diags
	}

	params := make(map[string]string)
	_, err := awxService.GetWorkflowJobTemplateNodeByID(id, params)
	if err != nil {
		return buildDiagNotFoundFail("workflow job template node", id, err)
	}

	_, err = awxService.UpdateWorkflowJobTemplateNode(id, map[string]interface{}{
		"extra_data":                d.Get("extra_data").(string),
		"inventory":                 d.Get("inventory_id").(int),
		"scm_branch":                d.Get("scm_branch").(string),
		"skip_tags":                 d.Get("skip_tags").(string),
		"job_type":                  d.Get("job_type").(string),
		"job_tags":                  d.Get("job_tags").(string),
		"limit":                     d.Get("limit").(string),
		"diff_mode":                 d.Get("diff_mode").(bool),
		"verbosity":                 d.Get("verbosity").(int),
		"workflow_job_template":     d.Get("workflow_job_template_id").(int),
		"unified_job_template":      d.Get("unified_job_template_id").(int),
		"all_parents_must_converge": d.Get("all_parents_must_converge").(bool),
		"identifier":                d.Get("identifier").(string),
	}, map[string]string{})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to update WorkflowJobTemplateNode",
			Detail:   fmt.Sprintf("WorkflowJobTemplateNode with name %s in the project id %d failed to update %s", d.Get("name").(string), d.Get("project_id").(int), err.Error()),
		})
		return diags
	}

	return resourceWorkflowJobTemplateNodeRead(ctx, d, m)
}

func resourceWorkflowJobTemplateNodeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*awx.AWX)
	awxService := client.WorkflowJobTemplateNodeService
	id, diags := convertStateIDToNumeric("Read WorkflowJobTemplateNode", d)
	if diags.HasError() {
		return diags
	}

	res, err := awxService.GetWorkflowJobTemplateNodeByID(id, make(map[string]string))
	if err != nil {
		return buildDiagNotFoundFail("workflow job template node", id, err)

	}
	d = setWorkflowJobTemplateNodeResourceData(d, res)
	return nil
}

func resourceWorkflowJobTemplateNodeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*awx.AWX)
	awxService := client.WorkflowJobTemplateNodeService
	id, diags := convertStateIDToNumeric(diagElementHostTitle, d)
	if diags.HasError() {
		return diags
	}

	if _, err := awxService.DeleteWorkflowJobTemplateNode(id); err != nil {
		return buildDiagDeleteFail(
			diagElementHostTitle,
			fmt.Sprintf("id %v, got %s ",
				id, err.Error()))
	}
	d.SetId("")
	return nil
}

func setWorkflowJobTemplateNodeResourceData(d *schema.ResourceData, r *awx.WorkflowJobTemplateNode) *schema.ResourceData {

	d.Set("extra_data", normalizeJsonYaml(r.ExtraData))
	d.Set("inventory_id", strconv.Itoa(r.Inventory))
	d.Set("scm_branch", r.ScmBranch)
	d.Set("job_type", r.JobType)
	d.Set("job_tags", r.JobTags)
	d.Set("skip_tags", r.SkipTags)
	d.Set("limit", r.Limit)
	d.Set("diff_mode", r.DiffMode)
	d.Set("verbosity", r.Verbosity)
	d.Set("workflow_job_template_id", strconv.Itoa(r.WorkflowJobTemplate))
	d.Set("unified_job_template_id", strconv.Itoa(r.UnifiedJobTemplate))
	d.Set("all_parents_must_converge", r.AllParentsMustConverge)
	d.Set("identifier", r.Identifier)

	d.SetId(strconv.Itoa(r.ID))
	return d
}

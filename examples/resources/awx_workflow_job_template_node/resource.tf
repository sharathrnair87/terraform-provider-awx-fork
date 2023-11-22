resource "random_uuid" "workflow_node_base_uuid" {}

resource "awx_workflow_job_template_node" "default" {
  workflow_job_template_id = awx_workflow_job_template.default.id
  unified_job_template_id  = awx_job_template.baseconfig.id
  inventory_id             = awx_inventory.default.id
  identifier               = random_uuid.workflow_node_base_uuid.result
}


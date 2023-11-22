resource "random_uuid" "workflow_node_k3s_uuid" {}

resource "awx_workflow_job_template_node_success" "k3s" {
  workflow_job_template_id      = awx_workflow_job_template.default.id
  workflow_job_template_node_id = awx_workflow_job_template_node.default.id
  unified_job_template_id       = awx_job_template.k3s.id
  inventory_id                  = awx_inventory.default.id
  identifier                    = random_uuid.workflow_node_k3s_uuid.result
}


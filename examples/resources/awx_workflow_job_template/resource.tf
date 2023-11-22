
resource "awx_workflow_job_template" "default" {
  name            = "workflow-job"
  organization_id = var.organization_id
  inventory_id    = awx_inventory.default.id
}


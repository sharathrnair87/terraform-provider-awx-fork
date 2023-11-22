
// Lookup by `name`
data "awx_workflow_job_template" "default" {
  name = "Default"
}

// Lookup by `id`
data "awx_workflow_job_template" "default" {
  id = var.default_workflow_job_template_id
}


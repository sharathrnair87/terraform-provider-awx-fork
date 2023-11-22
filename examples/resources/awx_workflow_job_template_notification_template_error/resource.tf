
resource "awx_workflow_job_template_notification_template_error" "baseconfig" {
  workflow_job_template_id = awx_workflow_job_template.baseconfig.id
  notification_template_id = awx_notification_template.default.id
}


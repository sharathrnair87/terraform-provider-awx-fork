
resource "awx_job_template_notification_template_success" "baseconfig" {
  job_template_id          = awx_job_template.baseconfig.id
  notification_template_id = awx_notification_template.default.id
}



resource "awx_workflow_job_template_schedule" "default" {
  workflow_job_template_id = awx_workflow_job_template.default.id

  name       = "schedule-test"
  rrule      = "DTSTART;TZID=Europe/Paris:20211214T120000 RRULE:INTERVAL=1;FREQ=DAILY"
  extra_data = <<EOL

organization_name: testorg
EOL
}


resource "awx_schedule" "default" {
  name                    = "schedule-test"
  rrule                   = "DTSTART;TZID=Europe/Paris:20211214T120000 RRULE:INTERVAL=1;FREQ=DAILY"
  unified_job_template_id = awx_job_template.baseconfig.id
  extra_data              = <<EOL

organization_name: testorg
EOL
}

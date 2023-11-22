
resource "awx_job_template_credential" "baseconfig" {
  job_template_id = awx_job_template.baseconfig.id
  credential_id   = awx_credential_machine.pi_connection.id
}


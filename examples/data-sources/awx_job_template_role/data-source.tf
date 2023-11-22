
resource "awx_job_template" "my_job_template" {
  name = "My JobTemplate"
}

data "awx_job_template_role" "job_template_admin_role" {
  role_name       = "Admin"
  job_template_id = data.awx_job_template.my_job_template.id
}


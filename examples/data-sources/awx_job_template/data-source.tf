
data "awx_job_template" "default" {
  name = "Default"
}

output "def_job_templ_playbook" {
  value = data.awx_job_template.default.playbook
}


data "awx_organization" "default" {
  name = "acc-test"
}

data "awx_job_template" "template" {
  name = "acc-job-template"
}

output "job" {
  value = data.awx_job_template.template
}

data "awx_project" "this" {
  name = "acc-project"
}

output "project" {
  value = data.awx_project.this
}

#data "awx_workflow_job_template" "template" {
#  name           = "acc-workflow-job"
#}
#
#output "workflow" {
#    value = data.awx_workflow_job_template.template
#}

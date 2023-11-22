
data "awx_execution_environment" "default" {
  name = "Default"
}

output "default_ee" {
  value = data.awx_execution_environment.default.id
}


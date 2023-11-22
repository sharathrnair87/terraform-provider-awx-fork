
data "awx_credential_machine" "my_machine_creds" {
  credential_id = var.my_machine_creds
}

output "my_machine_creds" {
  value     = data.awx_credential_machine.my_machine_creds
  sensitive = true
}


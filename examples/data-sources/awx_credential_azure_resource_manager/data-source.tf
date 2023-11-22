
data "awx_credential_azure_resource_manager" "my_azrm_creds" {
  credential_id = var.my_azrm_cred_id
}

output "my_azrm_creds" {
  value     = data.awx_credential_azure_resource_manager.my_azrm_creds
  sensitive = true
}


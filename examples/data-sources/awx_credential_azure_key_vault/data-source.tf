
data "awx_credential_azure_key_vault" "my_kv_lookup" {
  credential_id = var.my_kv_id
}

output "kv" {
  value     = data.awx_credential_azure_key_vault.my_kv_lookup
  sensitive = true
}


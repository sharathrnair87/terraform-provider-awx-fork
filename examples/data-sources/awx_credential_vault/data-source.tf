
data "awx_credential_vault" "my_vault_cred" {
  credential_id = var.my_vault_cred_id
}

output "my_vault_cred_vault_id" {
  value = data.awx_credential_vault.my_vault_cred.vault_id
}


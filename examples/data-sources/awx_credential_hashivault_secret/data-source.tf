
data "awx_credential_hashivault_secret" "my_hashi_secret_lookup" {
  credential_id = var.my_hashi_secret_lookup_id
}

output "my_hashi_secret_lookup" {
  value     = data.awx_credential_hashivault_secret.my_hashi_secret_lookup
  sensitive = true
}


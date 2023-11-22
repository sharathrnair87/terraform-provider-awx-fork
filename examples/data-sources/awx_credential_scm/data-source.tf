
data "awx_credential_scm" "my_scm_cred" {
  credential_id = var.my_scm_cred_id
}

output "my_scm_cred" {
  value     = data.awx_credential_scm.my_scm_cred
  sensitive = true
}


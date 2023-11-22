
data "awx_credential_hashivault_signed_ssh" "my_hashi_signed_ssh" {
  credential_id = var.my_hashi_signed_ssh_id
}

output "my_hashi_signed_ssh" {
  value     = data.awx_credential_hashivault_signed_ssh.my_hashi_signed_ssh
  sensitive = true
}


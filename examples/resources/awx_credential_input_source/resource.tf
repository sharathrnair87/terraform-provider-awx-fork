
data "awx_credential_machine" "ssh_key" {
  name = "iaas-ssh-key"
}

data "awx_credential_hashivault_signed_ssh" "hcv_signer" {
  name = "hcv-sig-ssh"
}

resource "awx_credential_input_source" "hcv-iaas-sig-map" {
  source      = data.awx_credential_hashivault_signed_ssh.hcv_signer.id
  target      = data.awx_credential_machine.ssh_key.id
  description = "Mapping Unsigned Machine Cred with HashiVault Signer"
  metadata = {
    role             = var.vault_role
    public_key       = var.rsa_public_key
    secret_path      = var.vault_secret_path
    valid_principals = var.vault_valid_principals
  }
}


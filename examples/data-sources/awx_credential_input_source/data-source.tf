
data "awx_credential_machine" "ssh_key" {
  name = "iaas-ssh-key"
}

data "awx_credential_hashivault_signed_ssh" "hcv_signer" {
  name = "hcv-sig-ssh"
}

data "awx_credential_input_source" "hcv-iaas-sig-map" {
  input_source_id = var.id_of_credential_input_source
}


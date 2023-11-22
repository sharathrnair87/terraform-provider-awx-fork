
data "awx_organization" "default" {
  name = "Default"
}

resource "awx_credential_machine" "machine_cred" {
  name            = "Machine Credential"
  organization_id = data.awx_organization.default.id
  username        = "testuser"
  password        = "securepassword"
  become_method   = "sudo"
  become_username = "root"
  become_password = "securepassword"
}


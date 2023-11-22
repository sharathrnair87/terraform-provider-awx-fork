
data "awx_organization" "my_org" {
  name = "My Org"
}

resource "awx_credential" "my_creds" {
  name               = "My Creds"
  description        = "My Machine Credentials"
  organization_id    = data.awx_organization.my_org.id
  credential_type_id = 1 // SSH Machine Credential
  inputs = jsonencode({
    username        = "testuser",
    password        = "securepassword",
    become_method   = "sudo",
    become_username = "root",
    become_password = "securepasssword"
  })
}


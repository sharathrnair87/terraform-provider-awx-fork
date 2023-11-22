
resource "awx_credential" "my_credential" {
  name = "My Credential"
}

data "awx_credential_role" "credential_admin_role" {
  role_name     = "Admin"
  credential_id = data.awx_credential.my_credential.id
}


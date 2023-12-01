
data "awx_organization" "default" {
  name = "Default"
}

data "awx_organization_role" "orga_read" {
  name            = "Read"
  organization_id = awx_organization.default.id
}

resource "awx_user" "my_user" {
  username = "my_user"
  password = "my_password"
  role_entitlement {
    role_id = data.awx_organization_role.orga_read.id
  }
}


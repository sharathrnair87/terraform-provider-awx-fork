
data "awx_organization" "default" {
  name = "Default"
}

resource "awx_organization" "admin_organization" {
  name            = "Admins"
  organization_id = data.awx_organization.default.id
}

resource "awx_settings_saml_organization_map" "admin_organization_map" {
  name          = resource.awx_organization.admin_organization.name
  users         = ["myorg-infra-users"]
  admins        = ["myorg-infra-admins"]
  organization  = data.awx_organization.default.name
  remove_users  = true
  remove_admins = true
}


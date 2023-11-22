
data "awx_organization" "default" {
  name = "Default"
}

resource "awx_settings_saml_team_attributes" "global" {
  saml_attr = "groups"
  remove    = true
  team_org_map {
    team         = var.saml_provider_team_id // The team ID as it is displayed in your SAML Auth Provider
    organization = data.awx_organization.default.name
    team_alias   = "Admin"
  }
}


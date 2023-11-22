
data "awx_organization" "default" {
  name = "Default"
}

resource "awx_team" "admin_team" {
  name            = "Admins"
  organization_id = data.awx_organization.default.id
}

resource "awx_settings_ldap_team_map" "admin_team_map" {
  name         = resource.awx_team.admin_team.name
  users        = ["CN=MyTeam,OU=Groups,DC=example,DC=com"]
  organization = data.awx_organization.default.name
  remove       = true
}


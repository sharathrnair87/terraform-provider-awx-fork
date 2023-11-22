
resource "awx_organization" "myorg" {
  name = "My AWX Org"
  // Truncated //
}

data "awx_organization_role" "org_admins" {
  name            = "Admin"
  organization_id = resource.awx_organization.myorg.id
}


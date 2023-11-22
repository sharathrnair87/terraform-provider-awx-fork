
data "awx_inventory" "default" {
  name            = "private_services"
  organization_id = data.awx_organization.default.id
}


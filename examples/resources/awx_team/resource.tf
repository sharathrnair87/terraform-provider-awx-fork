
data "awx_organization" "default" {
  name = "Default"
}

data "awx_inventory" "myinv" {
  name = "My Inventory"
}

data "awx_inventory_role" "myinv_admins" {
  name         = "Admin"
  inventory_id = data.awx_inventory.myinv.id
}

data "awx_project" "myproj" {
  name = "My Project"
}

data "awx_project_role" "myproj_admins" {
  name       = "Admin"
  project_id = data.awx_project.myproj.id
}

resource "awx_team" "admins_team" {
  name            = "admins-team"
  organization_id = data.awx_organization.default.id

  role_entitlement {
    role_id = data.awx_inventory_role.myinv_admins.id
  }
  role_entitlement {
    role_id = data.awx_project_role.myproj_admins.id
  }
}


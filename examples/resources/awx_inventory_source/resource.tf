
data "awx_inventory" "db_inventory" {
  name            = "DB_Inventory"
  organization_id = data.awx_organization.db.id
}

data "awx_project" "db_project" {
  name = "DB_Infra"
}

resource "awx_inventory_source" "db_inventory_source" {
  name              = "DB Inventory Src"
  inventory_id      = data.awx_inventory.db_inventory.id
  source_project_id = data.awx_project.db_project.id
  source_path       = "inventory/db-hosts.yml"
}



data "awx_inventory" "db_inventory" {
  name            = "DB_Inventory"
  organization_id = data.awx_organization.db.id
}

data "awx_inventory_source" "db_inventory_source" {
  name         = "DB_AZ_Inventory"
  inventory_id = data.awx_inventory.db_inventory.id
}


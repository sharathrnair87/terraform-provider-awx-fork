
data "awx_inventory" "db_inventory" {
  name            = "DB Inventory"
  organization_id = data.awx_organization.default.id
}

data "awx_host" "db_host" {
  name         = "prddbsrvr01.example.com"
  inventory_id = data.awx_inventory.db_inventory.id
}


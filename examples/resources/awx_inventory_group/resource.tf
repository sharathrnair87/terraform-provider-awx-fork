
data "awx_inventory" "db_inventory" {
  name            = "DB_Inventory"
  organization_id = data.awx_organization.db.id
}

resource "awx_inventory_group" "db_inventory_grp" {
  name         = "DB Inventory Grp"
  inventory_id = data.awx_inventory.db_inventory.id
  variables    = <<YAML

---
key: value
YAML
}


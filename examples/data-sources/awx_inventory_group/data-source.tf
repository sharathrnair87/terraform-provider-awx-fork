
data "awx_inventory_group" "default" {
  name         = "k3sPrimary"
  inventory_id = data.awx_inventory.default.id
}


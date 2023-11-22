
resource "awx_inventory" "myinv" {
  name = "My Inventory"
  // Truncated //
}

data "awx_inventory_role" "inv_admin_role" {
  name         = "Admin"
  inventory_id = data.awx_inventory.myinv.id
}

output "inv_admin_role_id" {
  value = data.awx_inventory_role.inv_admin_role.id
}


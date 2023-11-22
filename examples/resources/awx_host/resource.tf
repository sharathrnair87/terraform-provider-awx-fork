
resource "awx_host" "k3snode1" {
  name         = "k3snode1"
  description  = "pi node 1"
  inventory_id = data.awx_inventory.default.id
  group_ids = [
    data.awx_inventory_group.default.id,
    data.awx_inventory_group.pinodes.id,
  ]
  enabled   = true
  variables = <<YAML

---
ansible_host: 192.168.178.29
YAML
}

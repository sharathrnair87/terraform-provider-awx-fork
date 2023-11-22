
data "awx_organization" "default" {
  name = "Default"
}

resource "awx_inventory" "default" {
  name            = "acc-test"
  organization_id = data.awx_organization.default.id
  variables       = <<YAML

---
system_supporters:
  - pi

YAML
}

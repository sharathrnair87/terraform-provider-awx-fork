
data "awx_inventory" "default" {
  name            = "private_services"
  organization_id = data.awx_organization.default.id
}

resource "awx_job_template" "baseconfig" {
  name           = "baseconfig"
  job_type       = "run"
  inventory_id   = data.awx_inventory.default.id
  project_id     = awx_project.base_service_config.id
  playbook       = "master-configure-system.yml"
  become_enabled = true
}



data "awx_organization" "default" {
  name = "Default"
}

resource "awx_project" "base_service_config" {
  name                 = "base-service-configuration"
  scm_type             = "git"
  scm_url              = "https://github.com/nolte/ansible_playbook-baseline-online-server"
  scm_branch           = "feature/centos8-v2"
  scm_update_on_launch = true
  organization_id      = data.awx_organization.default.id
}


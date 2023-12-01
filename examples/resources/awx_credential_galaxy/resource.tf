
data "awx_organization" "devops" {
  name = "DevOps"
}

resource "awx_credential_galaxy" "devops_galaxy_cred_01" {
  name            = "DevOps_Cred_01"
  description     = "Galaxy Cred for DevOps Org"
  organization_id = data.awx_organization.devops
  url             = "https://galaxy.ansible.com"
}


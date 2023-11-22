
data "awx_organization" "cybersec" {
  name = "CyberSec"
}

resource "awx_credential_hashivault_secret" "hv_cyber" {
  name            = "HV Cyber"
  organization_id = data.awx_organization.cybersec.id
  url             = var.hashicorp_vault_url
  token           = var.hashicorp_vault_token
}



data "awx_organization" "cybersec" {
  name = "CyberSec"
}

resource "awx_credential_hashivault_signed_ssh" "hv_cyber_signed_ssh" {
  name            = "HV Cyber Sig SSH"
  organization_id = data.awx_organization.cybersec.id
  url             = var.hashicorp_vault_url
  token           = var.hashicorp_vault_token
}


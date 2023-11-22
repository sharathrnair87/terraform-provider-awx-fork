
data "awx_organization" "default" {
  name = "Default"
}

resource "awx_credential_vault" "at_iaas_vault" {
  name            = "at-iaas-vault"
  organization_id = data.awx_organization.default.id
  vault_id        = "at-iaas-vault"
  vault_password  = "securetvaultpassword"
}


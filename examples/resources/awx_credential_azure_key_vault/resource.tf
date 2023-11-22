
data "awx_organization" "infra" {
  name = "Infrastructure"
}

resource "awx_credential_azure_key_vault" "kv" {
  name            = "Infra KV"
  description     = "Azure KV for Infra project"
  organization_id = data.awx_organization.infra.id
  url             = "https://infra-vault-example.vault.azure.net"
  client          = var.azrm_client_id
  secret          = var.azrm_client_secret
  tenant          = var.azrm_tenant_id
}


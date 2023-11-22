
data "awx_organization" "infra" {
  name = "Infrastructure"
}

resource "awx_credential_azure_resource_manager" "azrm_infra" {
  name            = "AzureRM Infrastructure"
  subscription    = var.azrm_subscription_id
  tenant          = var.azrm_tenant_id
  organization_id = data.awx_organization.infra.id
  client          = var.azrm_client_id
  secret          = var.azrm_client_secret
}


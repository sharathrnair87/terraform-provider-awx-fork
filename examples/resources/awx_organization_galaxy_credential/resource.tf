
resource "awx_organization_galaxy_credential" "baseconfig" {
  organization_id = awx_organization.baseconfig.id
  credential_id   = awx_credential_machine.pi_connection.id
}


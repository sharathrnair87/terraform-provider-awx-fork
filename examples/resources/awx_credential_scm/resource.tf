
data "awx_organization" "devops" {
  name = "DevOps"
}

resource "awx_credential_scm" "scm_cred" {
  name            = "SCM Cred"
  organization_id = data.awx_organization.devops.id
  username        = "scmuser@example.com"
  password        = "securepasswordscm"
}


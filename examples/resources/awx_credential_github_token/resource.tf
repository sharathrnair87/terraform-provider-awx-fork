
data "awx_organization" "devops" {
  name = "DevOpsOrg"
}

resource "awx_credential_github_token" "gh_pat" {
  name            = "devops_gh_pat"
  token           = "..."
  organization_id = data.awx_organization.devops.id
}


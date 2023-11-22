
data "awx_credential_github_token" "my_gh_token" {
  credential_id = var.my_gh_token_id
}

output "my_gh_token" {
  value     = data.awx_credential_github_token.my_gh_token
  sensitive = true
}


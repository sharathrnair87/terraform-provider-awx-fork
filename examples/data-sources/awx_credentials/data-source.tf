// Query all available credentials

data "awx_credentials" "all_creds" {}

output "all_creds" {
  value = length(data.awx_credentials.all_creds.credentials)
}

// Query a specific credential by type

output "ssh" {
  value = toset([for each in data.awx_credentials.creds.credentials : each if each.kind == "ssh"])
}

// Query all credentials in a given Organization

data "awx_organization" "org" {
  name = "My Org"
}

output "creds_in_my_org" {
  value = toset([for each in data.awx_credentials.creds.credentials :
    each if each.organization_id == data.awx_organization.org.id
  ])
}


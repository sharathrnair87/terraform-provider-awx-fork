
data "awx_organization" "devops" {
  name = "DevOps"
}

resource "awx_execution_environment" "default" {
  name         = "acc-test"
  image        = "us-docker.pkg.dev/cloudrun/container/hello"
  pull         = "never"
  organization = data.awx_organization.devops.id
}


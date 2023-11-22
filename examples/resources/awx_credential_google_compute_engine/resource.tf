
data "awx_organization" "gcp_iaas" {
  name = "Google Infra"
}

resource "awx_credential_google_compute_engine" "gce" {
  name            = "GCE Credential"
  organization_id = data.awx_organization.gcp_iaas.id
  username        = "svc_acccount@gcp-prj.iam.gserviceaccount.com"
  ssh_key_data    = <<-EOT
	    -----BEGIN RSA PRIVATE KEY-----
	    ...
	    EOT
}


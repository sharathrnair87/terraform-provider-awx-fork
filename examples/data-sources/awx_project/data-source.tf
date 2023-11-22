
// By Name
data "awx_project" "default" {
  name = "Default"
}

// By ID
data "awx_project" "sharedServices" {
  id = var.shared_services_prj_id
}


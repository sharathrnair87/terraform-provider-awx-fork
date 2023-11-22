
resource "awx_project" "myproj" {
  name = "My AWX Project"
  // Truncated //
}

data "awx_project_role" "proj_admins" {
  name       = "Admin"
  project_id = resource.awx_project.myproj.id
}


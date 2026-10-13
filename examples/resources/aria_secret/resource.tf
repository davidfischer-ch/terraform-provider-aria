variable "database_password" {
  type      = string
  sensitive = true
}

resource "aria_secret" "example" {
  name        = "DATABASE_PASSWORD"
  description = "Password of the application's database user"
  org_scoped  = true

  # Write-only, Terraform never stores it. Changing it updates the secret.
  value = var.database_password
}

resource "aria_secret" "project_example" {
  name        = "PROJECT_TOKEN"
  org_scoped  = false
  project_ids = ["56ff588e-31ca-4878-883f-b74f28fa27b8"]
  value       = var.database_password
}

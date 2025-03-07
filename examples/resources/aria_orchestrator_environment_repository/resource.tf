# main.tf

resource "aria_orchestrator_environment_repository" "internal_python" {
  name     = "Internal_Python"
  runtime  = "python:3.10"
  location = "https://your-registry.your-company.net/repository/pypi-all/simple"
}

# main.tf

resource "aria_orchestrator_environment_repository" "internal_python" {
  name     = "Internal_Python"
  runtime  = "python:3.10"
  location = "https://registry.devops.etat-ge.ch/ctinexus/repository/pypi-all/simple"
}

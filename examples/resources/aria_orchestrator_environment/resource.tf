# main.tf

resource "aria_orchestrator_environment_repository" "python_ocsin" {
  name     = "Python_OCSIN"
  runtime  = "python:3.10"
  location = "https://your-registry.your-company.net/repository/pypi-all/simple"
}

resource "aria_orchestrator_environment" "python_for_tools" {
  name                 = "Python_For_Tools"
  description          = "Python runtime for our tools (packaged with common dependencies)."
  version              = "1.0.0"
  runtime              = aria_orchestrator_environment_repository.python_ocsin.runtime
  runtime_memory_limit = 256 * 1024 * 1024 # 256 MB
  runtime_timeout      = 180               # seconds

  dependencies = {
    build-tools = "== 3.20.2"
    pydantic    = "== 2.10.6"
    requests    = "== 2.32.3"
  }

  repositories = {
    build-tools = aria_orchestrator_environment_repository.python_ocsin.id
    pydantic    = aria_orchestrator_environment_repository.python_ocsin.id
    requests    = aria_orchestrator_environment_repository.python_ocsin.id
  }

  variables = {
    HTTPS_PROXY = "https://some-proxy.com"
    NO_PROXY    = "..."

    TERRAFORM_PATH = "/usr/bin/terraform1.10"
  }
}

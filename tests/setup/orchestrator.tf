# Orchestrator environment referenced by the action tests through two variables:
#
# * TF_VAR_test_environment_id
# * TF_VAR_test_environment_name
#
# Building an environment is slow and its runtime must be one the appliance offers, hence a single
# shared fixture: the tests that merely need an environment to point at reuse this one instead of
# building their own.

resource "aria_orchestrator_environment" "test" {
  name        = "${local.prefix}_ENVIRONMENT"
  description = "Environment referenced by Aria provider's acceptance tests."
  version     = "1.0.0"

  runtime              = local.test_runtime
  runtime_memory_limit = 0
  runtime_timeout      = 0

  dependencies = {}
  repositories = {}
  variables    = {}

  lifecycle {
    postcondition {
      condition = self.validation_message == ""
      error_message = join(" ", [
        "Runtime ${local.test_runtime} is not usable on this appliance",
        "(${self.validation_message}).",
        "Set TF_VAR_test_runtime to a runtime it offers.",
      ])
    }
  }
}

# Workflow referenced by the resource action tests through two variables:
#
# * TF_VAR_test_workflow_id
# * TF_VAR_test_workflow_name
#
# A resource action backed by a workflow needs the service broker to know it, and the import is
# asynchronous. Waiting for it once here is what lets every run consuming this workflow start from
# an imported one, instead of each test creating a workflow and waiting fifteen minutes for it.

resource "aria_orchestrator_category" "test" {
  name      = "${local.prefix}_WORKFLOWS"
  type      = "WorkflowCategory"
  parent_id = ""
}

resource "aria_orchestrator_workflow" "test" {
  name        = "${local.prefix}_WORKFLOW"
  description = "Workflow referenced by Aria provider's acceptance tests."
  category_id = aria_orchestrator_category.test.id
  version     = "0.1.0"

  position = { x = 100, y = 50 }

  restart_mode            = 1 # resume
  resume_from_failed_mode = 0 # default

  attrib        = jsonencode([])
  presentation  = jsonencode({})
  workflow_item = jsonencode([])

  input_parameters  = []
  output_parameters = []

  input_forms = jsonencode([
    {
      layout = {
        pages = []
      }
      schema = {}
    }
  ])

  wait_imported = true
}

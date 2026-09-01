# Orchestrator environment referenced by the action tests (TF_VAR_test_environment_id,
# TF_VAR_test_environment_name).
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

# main.tf

# Using javascript runtime
resource "aria_orchestrator_action" "get_deployment_by_id" {

  lifecycle {
    postcondition {
      condition     = length(self.validation_message) == 0
      error_message = "Validation of action ${self.fqn} failed with message ${self.validation_message}."
    }
  }

  name                 = "getDeploymentById" # You have to manage this boilerplate for some time
  module               = "com.company.core"  # Having name, module and fqn = module/name...
  fqn                  = "com.company.core/getDeploymentById"
  description          = "Return the deployment object matching given ID."
  version              = "1.0.0"
  runtime              = "" # Javascript, could be python:3.10 for Python ...
  runtime_memory_limit = 0
  runtime_timeout      = 0
  output_type          = "Any"

  input_parameters = [
    {
      name        = "vraHost"
      type        = "VRA:Host"
      description = ""
    },
    {
      name        = "deploymentId"
      type        = "string"
      description = ""
    }
  ]

  script = <<EOT
if (vraHost == null || deploymentId == null) return null;

var url = "/deployment/api/deployments/" + deploymentId;
var deployment = System.getModule("com.vmware.vra.extensibility.plugin.rest").getObjectFromUrl(vraHost, url);

return deployment;

EOT

}

# Using a custom execution environment
resource "aria_orchestrator_action" "foo_generator" {
  lifecycle {
    postcondition {
      condition     = length(self.validation_message) == 0
      error_message = "Validation of action ${self.fqn} failed with message ${self.validation_message}."
    }
  }

  name                 = "fooGenerator"
  module               = "com.company.core"
  fqn                  = "com.company.core/fooGenerator"
  description          = "Print 'hello' and return 'foo'."
  version              = "1.0.0"
  environment_id       = aria_orchestrator_environment.python_for_tools.id
  runtime_memory_limit = 0
  runtime_timeout      = 0
  output_type          = "string"

  input_parameters = []

  script = <<EOT
def handler(context, inputs):
    print('Hello')
    return 'foo'
EOT

}

# main.tf

resource "aria_orchestrator_action" "get_deployment_by_id" {
  name                 = "getDeploymentById" # You have to manage this boilerplate for some time
  module               = "com.company.core"  # Having name, module and fqn = module/name...
  fqn                  = "com.company.core/getDeploymentById"
  description          = "Return the deployment object matching given ID."
  version              = "1.0.0"
  runtime              = "" # javascript
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

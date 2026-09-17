provider "aria" {
  host = "https://my.aria-instance.net"
}

# Orchestrator is a standalone appliance declared as an integration in Aria: name it, and the
# provider resolves its endpoint among the integrations of the organization. Without this the
# aria_orchestrator_* resources are served by Aria's embedded Orchestrator.
provider "aria" {
  alias = "vro"

  host                 = "https://my.aria-instance.net"
  vro_integration_name = "External Orchestrator"
}

# Spelling the URI out works too, vro_host and vro_integration_name being mutually exclusive. The
# aria_integration data source is the way to read it from the platform, at the price of a second
# provider configuration: the data source is itself read through a configured provider, which a
# block feeding itself could not be.
data "aria_integration" "vro" {
  type_id = "com.vmw.vro.workflow"
  name    = "External Orchestrator"
}

provider "aria" {
  alias = "vro_by_uri"

  host     = "https://my.aria-instance.net"
  vro_host = data.aria_integration.vro.endpoint_uri
}

resource "aria_orchestrator_category" "example" {
  provider = aria.vro

  name      = "MyModule"
  type      = "WorkflowCategory"
  parent_id = ""
}

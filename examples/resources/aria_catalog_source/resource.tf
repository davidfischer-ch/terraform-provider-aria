# data.tf

data "aria_integration" "workflows" {
  type_id = "com.vmw.vro.workflow"
}

# main.tf

resource "aria_orchestrator_category" "my_company" {
  name      = "MyCompany"
  type      = "WorkflowCategory"
  parent_id = ""
}

# Example is dummy and contains no code
resource "aria_orchestrator_workflow" "dummy" {
  name        = "Dummy Workflow for Catalog Source"
  description = "Workflows doing nothing particular."
  category_id = aria_orchestrator_category.my_company.id
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
}

resource "aria_catalog_source" "dummy" {
  name    = "Dummy Workflow Catalog Source"
  type_id = data.aria_integration.workflows.type_id

  config = {
    workflows = [
      {
        id          = aria_orchestrator_workflow.dummy.id
        name        = aria_orchestrator_workflow.dummy.name
        description = aria_orchestrator_workflow.dummy.description
        version     = aria_orchestrator_workflow.dummy.version
        integration = {
          name                        = data.aria_integration.workflows.name
          endpoint_configuration_link = data.aria_integration.workflows.endpoint_configuration_link
          endpoint_uri                = data.aria_integration.workflows.endpoint_uri
        }
      }
    ]
  }
}

# Publish the Cloud Templates of a Project Libray using a Catalog Source ---------------------------

# variables.tf

variable "library_project_id" {
  description = "Identifier of the project containing Cloud templates to publish."
  type        = string
}

# main.tf

resource "aria_catalog_source" "library_project_cloud_templates" {
  name        = "Cloud Templates Catalog Source"
  description = "Publish some Cloud templates from a library project."
  project_id  = var.library_project_id
  type_id     = "com.vmw.abx.actions"

  config = {
    source_project_id = var.library_project_id
  }
}

# Create a Workflow and make it available using a Catalog Source -----------------------------------

# Method 1
#
# Using only the catalog source's waiting mechanism.
#
# Workflows's integration attribute will be null (cannot be guarantee).
# The aria integration data source is used to retrieve the integration endpoint.

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

  wait_imported = false
}

resource "aria_catalog_source" "dummy" {
  name        = "Dummy Workflow Catalog Source"
  description = "Publish the dummy workflow."
  type_id     = data.aria_integration.workflows.type_id

  config = {
    workflows = [
      {
        id          = aria_orchestrator_workflow.dummy.id
        name        = aria_orchestrator_workflow.dummy.name
        description = aria_orchestrator_workflow.dummy.description
        version     = aria_orchestrator_workflow.dummy.version
        integration = data.aria_integration.workflows
      }
    ]
  }

  # Refresh the catalog source every time the workflow is changed
  import_trigger = aria_orchestrator_workflow.dummy.version_id
}

# Create a Workflow and make it available using a Catalog Source -----------------------------------

# Method 2
#
# Using both workflow's andd catalog source's waiting mechanism.
#
# Workflows's integration attribute will be set.

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
  name        = "Dummy Workflow Catalog Source"
  description = "Publish the dummy workflow."
  type_id     = "com.vmw.vro.workflow"

  config = {
    workflows = [
      {
        id          = aria_orchestrator_workflow.dummy.id
        name        = aria_orchestrator_workflow.dummy.name
        description = aria_orchestrator_workflow.dummy.description
        version     = aria_orchestrator_workflow.dummy.version
        integration = aria_orchestrator_workflow.dummy.integration
      }
    ]
  }

  # Refresh the catalog source every time the workflow is changed
  import_trigger = aria_orchestrator_workflow.dummy.version_id
}

# main.tf

resource "aria_orchestrator_category" "root" {
  name      = "MyCompany"
  type      = "WorkflowCategory"
  parent_id = ""
}

resource "aria_orchestrator_category" "core" {
  name      = "Core"
  type      = "WorkflowCategory"
  parent_id = resource.aria_orchestrator_category.root.id
}

resource "aria_orchestrator_category" "mail" {
  name      = "Mail"
  type      = "WorkflowCategory"
  parent_id = resource.aria_orchestrator_category.core.id
}

resource "aria_orchestrator_category" "helpers" {
  name      = "Helpers"
  type      = "WorkflowCategory"
  parent_id = resource.aria_orchestrator_category.core.id
}

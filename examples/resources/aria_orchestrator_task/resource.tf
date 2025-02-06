# main.tf

resource "aria_orchestrator_category" "my_company" {
  name      = "MyCompany"
  type      = "WorkflowCategory"
  parent_id = ""
}

# Example is dummy and contains no code
resource "aria_orchestrator_workflow" "dummy" {
  name        = "Dummy Workflow for Task"
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

# Schedule the workflow we manage
resource "aria_orchestrator_task" "dummy_monthly" {
  name        = "Monthly execution of ${aria_orchestrator_workflow.dummy.name}"
  description = "Task doing nothing particular, on a monthly basis the 1st and 12th at midnight."

  recurrence_cycle      = "every-months"
  recurrence_pattern    = "(Europe/Zurich) 01 00:00:00,12 00:00:00,"
  recurrence_start_date = "1985-01-06T05:02:00Z"
  recurrence_end_date   = "2085-01-06T05:02:00Z"
  start_mode            = "normal"
  state                 = "pending"

  # input_parameters = [] # not yet implemented
  # Caveats:
  # * Terraform cannot detect modifications (drifts)
  # * Update method will always set it to []

  workflow = {
    id   = aria_orchestrator_workflow.dummy.id
    name = aria_orchestrator_workflow.dummy.name
  }
}

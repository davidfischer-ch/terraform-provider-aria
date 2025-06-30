# Lookup by ID -------------------------------------------------------------------------------------

data "aria_catalog_item" "my_item" {
  id = "f9f38e49-f463-3610-a799-54abacf1ed3c"
}

output "my_item" {
  value = data.aria_catalog_item.my_item
}

# Changes to Outputs:
#  + my_item = {
#      + created_at      = "2025-05-16T09:05:08.754866Z"
#      + created_by      = "SOMEONE"
#      + description     = "Action conernant PGSQL sur POWER/RHEL."
#      + external_id     = "63aff87f-64ab-4403-8e2f-a3d7bc4c49ef"
#      + form_id         = "c592fc7c-2ba8-44c1-be2a-dffb66d1b3b0"
#      + icon_id         = "4e9822c9-4053-3c1d-bb44-380a2fa1e028"
#      + id              = "f9f38e49-f463-3610-a799-54abacf1ed3c"
#      + last_updated_at = "2025-06-30T11:57:01.908410Z"
#      + last_updated_by = "system-user"
#      + name            = "VM Power with PostgreSQL update filesystem v1.0"
#      + schema          = jsonencode(
#            {
#              + properties = {
#                  + ...
#            }
#        )
#      + source_id       = ""
#      + source_name     = ""
#      + type_id         = "com.vmw.vro.workflow"
#    }

# Lookup by External ID and serch criteria ---------------------------------------------------------

data "aria_catalog_item" "my_item" {
  name        = aria_orchestrator_workflow.my_workflow.name # Optional but optimize query
  external_id = aria_orchestrator_workflow.my_workflow.id
  type_id     = "com.vmw.vro.workflow" # Optional but optimize query
}

output "my_item" {
  value = data.aria_catalog_item.my_item
}

# Changes to Outputs:
#  + my_item = {
#      + created_at      = "2025-05-16T09:05:08.754866Z"
#      + created_by      = "SOMEONE"
#      + description     = "Action conernant PGSQL sur POWER/RHEL."
#      + external_id     = "63aff87f-64ab-4403-8e2f-a3d7bc4c49ef"
#      + form_id         = "c592fc7c-2ba8-44c1-be2a-dffb66d1b3b0"
#      + icon_id         = "4e9822c9-4053-3c1d-bb44-380a2fa1e028"
#      + id              = "f9f38e49-f463-3610-a799-54abacf1ed3c"
#      + last_updated_at = "2025-06-30T11:57:01.908410Z"
#      + last_updated_by = "system-user"
#      + name            = "VM Power with PostgreSQL update filesystem v1.0"
#      + schema          = jsonencode(
#            {
#              + properties = {
#                  + ...
#            }
#        )
#      + source_id       = ""
#      + source_name     = ""
#      + type_id         = "com.vmw.vro.workflow"
#    }

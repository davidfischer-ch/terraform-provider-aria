data "aria_integration" "vro_workflows" {
  type_id = "com.vmw.vro.workflow"
}

output "vro_workflows_integration" {
  value = data.aria_integration.vro_workflows
}

# Changes to Outputs:
#   + vro_workflows_integration = {
#       + endpoint_configuration_link = "/resources/endpoints/8a430db3-924c-4d58-a29a-da811f9c992e"
#       + endpoint_uri                = "https://your-vra.your-company.net:443"
#       + name                        = "embedded-VRO"
#       + type_id                     = "com.vmw.vro.workflow"
#     }

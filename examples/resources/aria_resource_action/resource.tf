# variables.tf

variable "test_project_id" {
  type = string
}

# main.tf

resource "aria_abx_action" "hello_world" {
  name         = "Hello World"
  description  = "Say hello and display nice contextual data."
  runtime_name = "python"
  memory_in_mb = 128
  entrypoint   = "handler"
  dependencies = []
  constants    = []
  secrets      = []

  project_id = var.test_project_id

  source = <<EOT
from __future__ import annotations

import os


def handler(*args, **kwargs) -> None:
    print('Hello World!')
    print('Global symbols :', globals())
    print('Environment variables :', os.environ)
    print('Call Arguments: ', args, kwargs)
EOT

}

resource "aria_resource_action" "machine_hello_world" {
  name          = aria_abx_action.hello_world.name
  display_name  = aria_abx_action.hello_world.name
  description   = aria_abx_action.hello_world.description
  resource_type = "Cloud.vSphere.Machine"
  status        = "DRAFT"
  project_id    = var.test_project_id
  runnable_item = {
    id                = aria_abx_action.hello_world.id
    name              = aria_abx_action.hello_world.name
    project_id        = aria_abx_action.hello_world.project_id
    type              = "abx.action"
    input_parameters  = []
    output_parameters = []
  }

  // See https://docs.vmware.com/en/VMware-Aria-Automation/8.16/Using-Automation-Assembler/GUID-964816D8-DB67-406F-9224-DF597749397D.html
  criteria = jsonencode({
    matchExpression = [
      {
        and = [
          {
            key      = "$${properties.osType}"
            operator = "eq"
            value    = "WINDOWS"
          },
          {
            key      = "$${properties.totalMemoryMB}"
            operator = "greaterThan"
            value    = "1024"
          },
          {
            key      = "$${properties.tags}"
            operator = "hasAny"
            value = {
              matchExpression = [
                {
                  and = [
                    {
                      key      = "key"
                      operator = "eq"
                      value    = "env"
                    },
                    {
                      key      = "value"
                      operator = "eq"
                      value    = "REC"
                    }
                  ]
                }
              ]
            }
          }
        ]
      }
    ]
  })
}

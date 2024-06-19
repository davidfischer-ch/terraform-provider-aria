# variables.tf

variable "test_project_id" {
  type = string
}

# main.tf

// Not yet implemented
resource "aria_abx_action" "hello_world" {
  name         = "Hello World"
  description  = "Say hello and display nice contextual data."
  runtime_name = "python"
  memory_in_mb = 128
  entrypoint   = "handler"
  dependencies = []

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
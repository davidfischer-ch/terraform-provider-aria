# variables.tf

variable "project_id" {
  type = string
}

# main.tf

resource "aria_abx_constant" "hello_message" {
  name  = "HELLO_MESSAGE"
  value = "Hello World!"
}

resource "aria_abx_sensitive_constant" "some_secret" {
  name  = "SOME_SECRET"
  value = "sensitive stuff."
}

resource "aria_abx_action" "hello_world" {
  name         = "Hello World"
  description  = "Say hello and display nice contextual data."
  runtime_name = "python"
  memory_in_mb = 128
  entrypoint   = "handler"
  dependencies = []
  constants = [
    aria_abx_constant.hello_message.id,
    aria_abx_sensitive_constant.some_secret.id
  ]
  secrets = []

  project_id = var.project_id

  shared = true

  source = <<EOT
from __future__ import annotations

from typing import Any
import os


def handler(context, inputs: dict[str, Any]) -> None:
    print('Global symbols :', globals())
    print('Environment variables :', os.environ)
    print('Context: ', context)
    print('Inputs: , inputs)
    print(inputs['HELLO_MESSAGE'])
EOT

}

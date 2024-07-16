# variables.tf

variable "project_id" {
  type = string
}

# locals.tf

locals {
  source = <<EOT
import os

def handler(*args, **kwargs):
    print('Global symbols :', globals())
    print('Environment variables :', os.environ)
    print('Call Arguments: ', args, kwargs)
EOT
}

# main.tf

resource "aria_abx_action" "redis_create" {
  name            = "Custom.Redis.create"
  description     = "Provision an instance of a Redis server."
  runtime_name    = "python"
  memory_in_mb    = 128
  timeout_seconds = 60
  entrypoint      = "handler"
  dependencies    = []
  constants       = []
  secrets         = []
  source          = local.source
  shared          = true
  project_id      = var.project_id
}

resource "aria_abx_action" "redis_read" {
  name            = "Custom.Redis.read"
  description     = "Refresh properties by gathering the actual Redis instance attributes."
  runtime_name    = "python"
  memory_in_mb    = 128
  timeout_seconds = 60
  entrypoint      = "handler"
  dependencies    = []
  constants       = []
  secrets         = []
  source          = local.source
  shared          = true
  project_id      = var.project_id
}

resource "aria_abx_action" "redis_update" {
  name            = "Custom.Redis.update"
  description     = "Update Redis instance's attributes."
  runtime_name    = "python"
  memory_in_mb    = 128
  timeout_seconds = 60
  entrypoint      = "handler"
  dependencies    = []
  constants       = []
  secrets         = []
  source          = local.source
  shared          = true
  project_id      = var.project_id
}

resource "aria_abx_action" "redis_delete" {
  name            = "Custom.Redis.delete"
  description     = "Destroy the Redis instance."
  runtime_name    = "python"
  memory_in_mb    = 128
  timeout_seconds = 60
  entrypoint      = "handler"
  dependencies    = []
  constants       = []
  secrets         = []
  source          = local.source
  shared          = true
  project_id      = var.project_id
}

resource "aria_custom_resource" "redis" {
  display_name  = "Redis"
  description   = "Manage an instance of a Redis database."
  resource_type = "Custom.Redis"
  schema_type   = "ABX_USER_DEFINED"
  status        = "DRAFT"
  project_id    = var.project_id

  properties = [
    {
      name        = "version"
      title       = "Version"
      description = "Instance version."
      type        = "string"
      one_of = [
        { const = "7.4", title = "7.4", encrypted = false },
        { const = "8.0", title = "8.0", encrypted = false }
      ]
    },
    {
      name               = "storage_size"
      title              = "Storage Size"
      description        = "Storage size (MB)."
      type               = "integer"
      default            = tostring(10 * 1024)
      encrypted          = false
      read_only          = false
      recreate_on_update = false
      minimum            = 1 * 1024
      maximum            = 100 * 1024
      one_of             = []
    },
    {
      name               = "secret"
      title              = "Secret"
      description        = "Secret key."
      type               = "string"
      encrypted          = true
      read_only          = false
      recreate_on_update = false
      min_length         = 16
      max_length         = 64
      one_of             = []
    }
  ]

  create = {
    id                = aria_abx_action.redis_create.id
    project_id        = aria_abx_action.redis_create.project_id
    type              = "abx.action"
    input_parameters  = []
    output_parameters = []
  }

  read = {
    id                = aria_abx_action.redis_read.id
    project_id        = aria_abx_action.redis_read.project_id
    type              = "abx.action"
    input_parameters  = []
    output_parameters = []
  }

  update = {
    id                = aria_abx_action.redis_update.id
    project_id        = aria_abx_action.redis_update.project_id
    type              = "abx.action"
    input_parameters  = []
    output_parameters = []
  }

  delete = {
    id                = aria_abx_action.redis_delete.id
    project_id        = aria_abx_action.redis_delete.project_id
    type              = "abx.action"
    input_parameters  = []
    output_parameters = []
  }
}

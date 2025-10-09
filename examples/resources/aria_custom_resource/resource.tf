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

# constants.tf

resource "aria_abx_constant" "example" {
  name  = "THIS_IS_MY_CONSTANT"
  value = "42"
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
  inputs          = {}
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
  inputs          = {}
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
  inputs          = {}
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
  inputs          = {}
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
  status        = "RELEASED"
  project_id    = var.project_id

  properties = {
    version = {
      name               = "version"
      title              = "Version"
      description        = "Instance version."
      type               = "string"
      encrypted          = false
      read_only          = false
      recreate_on_update = false
      one_of = [
        { const = "7.4", title = "7.4", encrypted = false },
        { const = "8.0", title = "8.0", encrypted = false }
      ]
    }
    description = {
      name               = "description"
      title              = "Description"
      description        = "Some description here."
      type               = "string"
      default            = jsonencode("No description given.")
      encrypted          = false
      read_only          = false
      recreate_on_update = false
    }
    storage_size = {
      name               = "storage_size"
      title              = "Storage Size"
      description        = "Storage size (MB)."
      type               = "integer"
      default            = 10 * 1024
      encrypted          = false
      read_only          = false
      recreate_on_update = false
      minimum            = 1 * 1024
      maximum            = 100 * 1024
    }
    secret = {
      name               = "secret"
      title              = "Secret"
      description        = "Secret key."
      type               = "string"
      encrypted          = true
      read_only          = false
      recreate_on_update = false
      min_length         = 16
      max_length         = 64
    }
  }

  create = {
    id                = aria_abx_action.redis_create.id
    name              = aria_abx_action.redis_create.name
    project_id        = aria_abx_action.redis_create.project_id
    type              = "abx.action"
    input_parameters  = []
    output_parameters = []
  }

  read = {
    id                = aria_abx_action.redis_read.id
    name              = aria_abx_action.redis_read.name
    project_id        = aria_abx_action.redis_read.project_id
    type              = "abx.action"
    input_parameters  = []
    output_parameters = []
  }

  update = {
    id                = aria_abx_action.redis_update.id
    name              = aria_abx_action.redis_update.name
    project_id        = aria_abx_action.redis_update.project_id
    type              = "abx.action"
    input_parameters  = []
    output_parameters = []
  }

  delete = {
    id                = aria_abx_action.redis_delete.id
    name              = aria_abx_action.redis_delete.name
    project_id        = aria_abx_action.redis_delete.project_id
    type              = "abx.action"
    input_parameters  = []
    output_parameters = []
  }
}

# Additional actions (aka Day 2), managed using relational resources
# This design is intentional for Terraform to be able to succesfully apply any changes

resource "aria_abx_action" "redis_backup" {
  name            = "Custom.Redis.backup"
  description     = "Backup the Redis database (its data)."
  runtime_name    = "python"
  memory_in_mb    = 128
  timeout_seconds = 60
  entrypoint      = "handler"
  dependencies    = []
  constants       = [aria_abx_constant.example.id]
  inputs          = {}
  secrets         = []
  source          = local.source
  shared          = true
  project_id      = var.project_id
}

resource "aria_resource_action" "redis_backup" {
  name          = "backup"
  display_name  = "Backup data"
  description   = aria_abx_action.redis_backup.description
  status        = aria_custom_resource.redis.status
  resource_id   = aria_custom_resource.redis.id
  resource_type = aria_custom_resource.redis.resource_type
  project_id    = aria_custom_resource.redis.project_id
  runnable_item = {
    id                = aria_abx_action.redis_backup.id
    name              = aria_abx_action.redis_backup.name
    project_id        = aria_abx_action.redis_backup.project_id
    type              = "abx.action"
    input_parameters  = []
    output_parameters = []
  }
}

resource "aria_abx_action" "redis_restore" {
  name            = "Custom.Redis.restore"
  description     = "Restore the Redis database (its data)."
  runtime_name    = "python"
  memory_in_mb    = 128
  timeout_seconds = 60
  entrypoint      = "handler"
  dependencies    = []
  constants       = []
  inputs          = {}
  secrets         = []
  source          = local.source
  shared          = true
  project_id      = var.project_id
}

resource "aria_resource_action" "redis_restore" {
  name          = "restore"
  display_name  = "Restore"
  description   = aria_abx_action.redis_restore.description
  status        = aria_custom_resource.redis.status
  resource_id   = aria_custom_resource.redis.id
  resource_type = aria_custom_resource.redis.resource_type
  project_id    = aria_custom_resource.redis.project_id
  runnable_item = {
    id                = aria_abx_action.redis_restore.id
    name              = aria_abx_action.redis_restore.name
    project_id        = aria_abx_action.redis_restore.project_id
    type              = "abx.action"
    input_parameters  = []
    output_parameters = []
  }
}

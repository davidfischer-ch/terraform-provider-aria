// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomResourceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: `
variable "test_project_id" {
    description = "Project where to generate test resources."
  type        = string
}

locals {
    source = <<EOT
import os

def handler(*args, **kwargs):
    print('Global symbols :', globals())
    print('Environment variables :', os.environ)
    print('Call Arguments: ', args, kwargs)
EOT
}

resource "aria_abx_action" "create" {
  name            = "Custom.AriaProviderTest.create"
  description     = "Temporary action generated by Aria provider's acceptance tests."
  runtime_name    = "python"
  memory_in_mb    = 128
  timeout_seconds = 60
  entrypoint      = "handler"
  dependencies    = []
  constants       = []
  secrets         = []
  source          = local.source
  shared          = true
  project_id      = var.test_project_id
}

resource "aria_abx_action" "read" {
  name            = "Custom.AriaProviderTest.read"
  description     = "Temporary action generated by Aria provider's acceptance tests."
  runtime_name    = "python"
  memory_in_mb    = 128
  timeout_seconds = 60
  entrypoint      = "handler"
  dependencies    = []
  constants       = []
  secrets         = []
  source          = local.source
  shared          = true
  project_id      = var.test_project_id
}

resource "aria_abx_action" "update" {
  name            = "Custom.AriaProviderTest.update"
  description     = "Temporary action generated by Aria provider's acceptance tests."
  runtime_name    = "python"
  memory_in_mb    = 128
  timeout_seconds = 60
  entrypoint      = "handler"
  dependencies    = []
  constants       = []
  secrets         = []
  source          = local.source
  shared          = true
  project_id      = var.test_project_id
}

resource "aria_abx_action" "delete" {
  name            = "Custom.AriaProviderTest.delete"
  description     = "Temporary action generated by Aria provider's acceptance tests."
  runtime_name    = "python"
  memory_in_mb    = 128
  timeout_seconds = 60
  entrypoint      = "handler"
  dependencies    = []
  constants       = []
  secrets         = []
  source          = local.source
  shared          = true
  project_id      = var.test_project_id
}

resource "aria_custom_resource" "test" {
  display_name  = "Aria Provider Test Custom Resource"
  description   = "Temporary custom resource generated by Aria provider's acceptance tests."
  resource_type = "Custom.AriaProviderTest"
  schema_type   = "ABX_USER_DEFINED"
  status        = "DRAFT"
  project_id    = var.test_project_id

  properties = [
    {
      name        = "some_text"
      title       = "Some Text"
      description = "Some text, more text."
      type        = "string"
      one_of      = [
        { const = "a", title = "A", encrypted = false },
        { const = "b", title = "B", encrypted = false }
      ]
    },
    {
      name        = "number"
      title       = "Some Number"
      description = <<EOT
Some number.
It can be an integer or a float.
EOT
      type    = "number"
      default = "3.141592"
      minimum = 0
      maximum = 5
      one_of  = []
    },
    {
      name        = "super_secret"
      title       = "Super Secret"
      description = ""
      type        = "string"
      encrypted   = true
      min_length  = 16
      max_length  = 64
      one_of      = []
    }
  ]

  create = {
    id                = aria_abx_action.create.id
    project_id        = aria_abx_action.create.project_id
    type              = "abx.action"
    input_parameters  = []
    output_parameters = []
  }

  read = {
    id                = aria_abx_action.read.id
    project_id        = aria_abx_action.read.project_id
    type              = "abx.action"
    input_parameters  = []
    output_parameters = []
  }

  update = {
    id                = aria_abx_action.update.id
    project_id        = aria_abx_action.update.project_id
    type              = "abx.action"
    input_parameters  = []
    output_parameters = []
  }

  delete = {
    id                = aria_abx_action.delete.id
    project_id        = aria_abx_action.delete.project_id
    type              = "abx.action"
    input_parameters  = []
    output_parameters = []
  }
}
`,
				/*Check: resource.ComposeAggregateTestCheckFunc(
				    resource.TestCheckResourceAttrSet("aria_abx_action.test", "id"),
				    resource.TestCheckResourceAttr("aria_abx_action.test", "name", "ARIA_PROVIDER_TEST_ACTION"),
				    resource.TestCheckResourceAttr("aria_abx_action.test", "description", "Temporary action generated by Aria provider's acceptance tests."),
				    resource.TestCheckResourceAttr("aria_abx_action.test", "faas_provider", "auto"),
				    resource.TestCheckResourceAttr("aria_abx_action.test", "type", "SCRIPT"),
				    resource.TestCheckResourceAttr("aria_abx_action.test", "runtime_name", "python"),
				    // resource.TestCheckResourceAttrSet("aria_abx_action.test", "runtime_version"),
				    resource.TestCheckResourceAttr("aria_abx_action.test", "cpu_shares", "1024"),
				    resource.TestCheckResourceAttr("aria_abx_action.test", "memory_in_mb", "128"),
				    resource.TestCheckResourceAttr("aria_abx_action.test", "timeout_seconds", "60"),
				    resource.TestCheckResourceAttr("aria_abx_action.test", "deployment_timeout_seconds", "900"),
				    resource.TestCheckResourceAttr("aria_abx_action.test", "entrypoint", "handler"),
				    resource.TestCheckResourceAttr("aria_abx_action.test", "shared", "true"),
				    resource.TestCheckResourceAttr("aria_abx_action.test", "system", "false"),
				    resource.TestCheckResourceAttr("aria_abx_action.test", "async_deployed", "false"),
				    resource.TestCheckResourceAttrSet("aria_abx_action.test", "org_id"),
				),*/
			},
			// Update and "Un"-Scoping Test
			{
				Config: `
variable "test_project_id" {
    description = "Project where to generate test resources."
  type        = string
}

locals {
    source = <<EOT
import os

def handler(*args, **kwargs):
    print('Global symbols :', globals())
    print('Environment variables :', os.environ)
    print('Call Arguments: ', args, kwargs)
EOT
}

resource "aria_abx_action" "create" {
  name            = "Custom.AriaProviderTest.create"
  description     = "Temporary action generated by Aria provider's acceptance tests."
  runtime_name    = "python"
  memory_in_mb    = 128
  timeout_seconds = 60
  entrypoint      = "handler"
  dependencies    = []
  constants       = []
  secrets         = []
  source          = local.source
  shared          = true
  project_id      = var.test_project_id
}

resource "aria_abx_action" "read" {
  name            = "Custom.AriaProviderTest.read"
  description     = "Temporary action generated by Aria provider's acceptance tests."
  runtime_name    = "python"
  memory_in_mb    = 128
  timeout_seconds = 60
  entrypoint      = "handler"
  dependencies    = []
  constants       = []
  secrets         = []
  source          = local.source
  shared          = true
  project_id      = var.test_project_id
}

resource "aria_abx_action" "update" {
  name            = "Custom.AriaProviderTest.update"
  description     = "Temporary action generated by Aria provider's acceptance tests."
  runtime_name    = "python"
  memory_in_mb    = 128
  timeout_seconds = 60
  entrypoint      = "handler"
  dependencies    = []
  constants       = []
  secrets         = []
  source          = local.source
  shared          = true
  project_id      = var.test_project_id
}

resource "aria_abx_action" "delete" {
  name            = "Custom.AriaProviderTest.delete"
  description     = "Temporary action generated by Aria provider's acceptance tests."
  runtime_name    = "python"
  memory_in_mb    = 128
  timeout_seconds = 60
  entrypoint      = "handler"
  dependencies    = []
  constants       = []
  secrets         = []
  source          = local.source
  shared          = true
  project_id      = var.test_project_id
}

resource "aria_custom_resource" "test" {
  display_name  = "Aria Provider Test Custom Resource"
  description   = "Temporary custom resource generated by Aria provider's acceptance tests."
  resource_type = "Custom.AriaProviderTest"
  schema_type   = "ABX_USER_DEFINED"
  status        = "DRAFT"

  properties = [
    {
      name        = "some_text"
      title       = "Some Text"
      description = "Some text, more text."
      type        = "string"
      one_of      = [
        { const = "a", title = "A", encrypted = false },
        { const = "b", title = "B", encrypted = false }
      ]
    },
    {
      name        = "number"
      title       = "Some Number"
      description = <<EOT
Some number.
It can be an integer or a float.
EOT
      type    = "number"
      default = "3.141592"
      minimum = 0
      maximum = 5
      one_of  = []
    },
    {
      name        = "super_secret"
      title       = "Super Secret"
      description = ""
      type        = "string"
      encrypted   = true
      min_length  = 16
      max_length  = 64
      one_of      = []
    }
  ]

  create = {
    id                = aria_abx_action.create.id
    project_id        = aria_abx_action.create.project_id
    type              = "abx.action"
    input_parameters  = []
    output_parameters = []
  }

  read = {
    id                = aria_abx_action.read.id
    project_id        = aria_abx_action.read.project_id
    type              = "abx.action"
    input_parameters  = []
    output_parameters = []
  }

  update = {
    id                = aria_abx_action.update.id
    project_id        = aria_abx_action.update.project_id
    type              = "abx.action"
    input_parameters  = []
    output_parameters = []
  }

  delete = {
    id                = aria_abx_action.delete.id
    project_id        = aria_abx_action.delete.project_id
    type              = "abx.action"
    input_parameters  = []
    output_parameters = []
  }
}
`,
				/*Check: resource.ComposeAggregateTestCheckFunc(
				    resource.TestCheckResourceAttrSet("aria_abx_action.test", "id"),
				    resource.TestCheckResourceAttr("aria_abx_action.test", "name", "ARIA_PROVIDER_TEST_ACTION"),
				    resource.TestCheckResourceAttr("aria_abx_action.test", "description", "Temporary action generated by Aria provider's acceptance tests."),
				    resource.TestCheckResourceAttr("aria_abx_action.test", "faas_provider", "auto"),
				    resource.TestCheckResourceAttr("aria_abx_action.test", "type", "SCRIPT"),
				    resource.TestCheckResourceAttr("aria_abx_action.test", "runtime_name", "python"),
				    // resource.TestCheckResourceAttrSet("aria_abx_action.test", "runtime_version"),
				    resource.TestCheckResourceAttr("aria_abx_action.test", "cpu_shares", "1024"),
				    resource.TestCheckResourceAttr("aria_abx_action.test", "memory_in_mb", "128"),
				    resource.TestCheckResourceAttr("aria_abx_action.test", "timeout_seconds", "60"),
				    resource.TestCheckResourceAttr("aria_abx_action.test", "deployment_timeout_seconds", "900"),
				    resource.TestCheckResourceAttr("aria_abx_action.test", "entrypoint", "handler"),
				    resource.TestCheckResourceAttr("aria_abx_action.test", "shared", "true"),
				    resource.TestCheckResourceAttr("aria_abx_action.test", "system", "false"),
				    resource.TestCheckResourceAttr("aria_abx_action.test", "async_deployed", "false"),
				    resource.TestCheckResourceAttrSet("aria_abx_action.test", "org_id"),
				),*/
			},
			// Delete testing automatically occurs in TestCase
			// TODO Check https://developer.hashicorp.com/terraform/plugin/sdkv2/testing/acceptance-tests/testcase#checkdestroy
		},
	})
}

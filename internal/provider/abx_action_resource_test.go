// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccABXActionResource(t *testing.T) {
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

resource "aria_abx_constant" "test" {
	count = 3
  name  = "ARIA_PROVIDER_TEST_ACTION_CONSTANT_${count.index}"
  value = "Some value."
}

resource "aria_abx_sensitive_constant" "test" {
	count = 2
  name  = "ARIA_PROVIDER_TEST_ACTION_SECRET_${count.index}"
  value = "sensitive stuff."
}

locals {
	constants = concat(
		[for constant in aria_abx_constant.test: constant.id],
		[for constant in aria_abx_sensitive_constant.test: constant.id])

	source = <<EOT
import os

def handler(*args, **kwargs):
		print('Global symbols :', globals())
		print('Environment variables :', os.environ)
		print('Call Arguments: ', args, kwargs)
EOT
}

resource "aria_abx_action" "test" {
  name         = "ARIA_PROVIDER_TEST_ACTION"
  description  = "Temporary action generated by Aria provider's acceptance tests."
  runtime_name = "python"
  memory_in_mb = 128
  entrypoint   = "handler"
  dependencies = []
  constants    = local.constants
  secrets      = [] # TODO Test this once secret is available

	project_id = var.test_project_id

  source = local.source

  lifecycle {
    postcondition {
      condition     = length(self.dependencies) == 0
      error_message = "Dependencies must be empty, actual [${join(", ", self.dependencies)}]"
    }
    postcondition {
      condition     = self.constants == toset(local.constants)
      error_message = "Constants must be [${join(", ", local.constants)}], actual [${join(", ", self.constants)}]"
    }
    postcondition {
      condition     = length(self.secrets) == 0
      error_message = "Secrets must be empty, actual [${join(", ", self.secrets)}]"
    }
    postcondition {
      condition     = self.project_id == var.test_project_id
      error_message = "Project ID must be ${var.test_project_id}, actual ${self.project_id}"
    }
    postcondition {
      condition     = self.source == local.source
      error_message = "Source must be ${local.source}, actual ${self.source}"
    }
  }
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_abx_action.test", "id"),
					resource.TestCheckResourceAttrSet("aria_abx_action.test", "org_id"),
					resource.TestCheckResourceAttr("aria_abx_action.test", "name", "ARIA_PROVIDER_TEST_ACTION"),
					resource.TestCheckResourceAttr("aria_abx_action.test", "description", "Temporary action generated by Aria provider's acceptance tests."),
					resource.TestCheckResourceAttr("aria_abx_action.test", "runtime_name", "python"),
					resource.TestCheckResourceAttr("aria_abx_action.test", "memory_in_mb", "128"),
					resource.TestCheckResourceAttr("aria_abx_action.test", "entrypoint", "handler"),
				),
			},
			// ImportState testing
			/*{
				ResourceName:      "aria_abx_action.test",
				ImportState:       true,
				ImportStateVerify: true,
			},*/
			// Update and Read testing
			{
				Config: `
variable "test_project_id" {
	description = "Project where to generate test resources."
  type        = string
}

resource "aria_abx_constant" "test" {
	count = 3
  name  = "ARIA_PROVIDER_TEST_ACTION_CONSTANT_${count.index}"
  value = "Some value."
}

resource "aria_abx_sensitive_constant" "test" {
	count = 2
  name  = "ARIA_PROVIDER_TEST_ACTION_SECRET_${count.index}"
  value = "sensitive stuff."
}

locals {
	dependencies = ["requests", "pytoolbox==14.8.2"]

	source = <<EOT
from __future__ import annotations

import os

import requests


def handler(*args, **kwargs) -> None:
		print('Global symbols :', globals())
		print('Environment variables :', os.environ)
		print('Call Arguments: ', args, kwargs)
		print('Requests module: ', requests)
EOT
}

resource "aria_abx_action" "test" {
  name          = "ARIA_PROVIDER_TEST_ACTION_RENAMED"
  description   = "Temporary action generated by Aria provider's acceptance tests (changed)."
  faas_provider = "on-prem"
  runtime_name  = "python"
  memory_in_mb  = 64
  entrypoint    = "handler"
  dependencies  = local.dependencies
  constants     = []
  secrets       = []

	project_id = var.test_project_id

  source = local.source

  lifecycle {
    postcondition {
      condition     = self.dependencies == tolist(local.dependencies)
      error_message = "Dependencies must be [${join(", ", local.dependencies)}], actual [${join(", ", self.dependencies)}]"
    }
    postcondition {
      condition     = length(self.constants) == 0
      error_message = "Constants must be empty, actual [${join(", ", self.constants)}]"
    }
    postcondition {
      condition     = length(self.secrets) == 0
      error_message = "Secrets must be empty, actual [${join(", ", self.secrets)}]"
    }
    postcondition {
      condition     = self.project_id == var.test_project_id
      error_message = "Project ID must be ${var.test_project_id}, actual ${self.project_id}"
    }
    postcondition {
      condition     = self.source == local.source
      error_message = "Source must be ${local.source}, actual ${self.source}"
    }
  }
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_abx_action.test", "id"),
					resource.TestCheckResourceAttrSet("aria_abx_action.test", "org_id"),
					resource.TestCheckResourceAttr("aria_abx_action.test", "name", "ARIA_PROVIDER_TEST_ACTION_RENAMED"),
					resource.TestCheckResourceAttr("aria_abx_action.test", "description", "Temporary action generated by Aria provider's acceptance tests (changed)."),
					resource.TestCheckResourceAttr("aria_abx_action.test", "runtime_name", "python"),
					resource.TestCheckResourceAttr("aria_abx_action.test", "memory_in_mb", "64"),
					resource.TestCheckResourceAttr("aria_abx_action.test", "entrypoint", "handler"),
				),
			},
			// Delete testing automatically occurs in TestCase
			// TODO Check https://developer.hashicorp.com/terraform/plugin/sdkv2/testing/acceptance-tests/testcase#checkdestroy
		},
	})
}

// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSubscriptionResource(t *testing.T) {
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

variable "test_project_ids" {
  description = "Scoping some resources to given projects (for testing purposes)."
  type        = string
}

variable "test_abx_action_id" {
  description = "ABX action to use for testing subscriptions."
  type        = string
}

locals {
  project_ids   = split(",", var.test_project_ids)
  subscriber_id = data.aria_catalog_type.abx_actions.created_by
}

data "aria_catalog_type" "abx_actions" {
  id = "com.vmw.abx.actions"
}

resource "aria_subscription" "hello_world" {
  name           = "ARIA_PROVIDER_TEST_SUBSCRIPTION"
  description    = "Say hello when some machine is provisionned"
  type           = "RUNNABLE"
  runnable_type  = "extensibility.abx"
  runnable_id    = var.test_abx_action_id
  event_topic_id = "compute.provision.post"
  project_ids    = []
  blocking       = true
  contextual     = false
  disabled       = true # Its safer
  timeout        = 0
  priority       = 10

  lifecycle {
    postcondition {
      condition     = can(regex("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}", self.id))
      error_message = "Identifier must be a valid UUID string, actual ${self.id}"
    }
    postcondition {
      condition     = self.runnable_id == var.test_abx_action_id
      error_message = "Runnable ID must be ${var.test_abx_action_id}, actual ${self.runnable_id}"
    }
    postcondition {
      condition     = self.subscriber_id == local.subscriber_id
      error_message = "Subscriber ID must be ${local.subscriber_id}, actual ${self.subscriber_id}"
    }
    postcondition {
      condition     = length(self.project_ids) == 0
      error_message = "Project IDs must be [], actual [${join(", ", self.project_ids)}]"
    }
  }
}

resource "aria_subscription" "hello_world_scoped" {
  name           = "ARIA_PROVIDER_TEST_SUBSCRIPTION_SCOPED"
  description    = "Say hello when any machine is provisionned"
  type           = "RUNNABLE"
  runnable_type  = "extensibility.abx"
  runnable_id    = var.test_abx_action_id
  event_topic_id = "compute.provision.post"
  project_ids    = local.project_ids
  blocking       = true
  contextual     = false
  disabled       = true # Its safer
  timeout        = 0
  priority       = 10

  lifecycle {
    postcondition {
      condition     = can(regex("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}", self.id))
      error_message = "Identifier must be a valid UUID string, actual ${self.id}"
    }
    postcondition {
      condition     = self.runnable_id == var.test_abx_action_id
      error_message = "Runnable ID must be ${var.test_abx_action_id}, actual ${self.runnable_id}"
    }
    postcondition {
      condition     = self.subscriber_id == local.subscriber_id
      error_message = "Subscriber ID must be ${local.subscriber_id}, actual ${self.subscriber_id}"
    }
    postcondition {
      condition     = self.project_ids == toset(local.project_ids)
      error_message = "Project IDs must be [${join(", ", local.project_ids)}], actual [${join(", ", self.project_ids)}]"
    }
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_subscription.hello_world", "id"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "name", "ARIA_PROVIDER_TEST_SUBSCRIPTION"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "description", "Say hello when some machine is provisionned"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "type", "RUNNABLE"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "runnable_type", "extensibility.abx"),
					//resource.TestCheckResourceAttr("aria_subscription.hello_world", "recover_runnable_type", ""),
					//resource.TestCheckResourceAttr("aria_subscription.hello_world", "recover_runnable_id", ""),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "event_topic_id", "compute.provision.post"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "blocking", "true"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "broadcast", "false"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "contextual", "false"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "criteria", ""),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "disabled", "true"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "priority", "10"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "system", "false"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "timeout", "0"),
					resource.TestCheckResourceAttrSet("aria_subscription.hello_world", "org_id"),
					resource.TestCheckResourceAttrSet("aria_subscription.hello_world", "owner_id"),
					resource.TestCheckResourceAttrSet("aria_subscription.hello_world", "subscriber_id"),

					resource.TestCheckResourceAttrSet("aria_subscription.hello_world_scoped", "id"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world_scoped", "name", "ARIA_PROVIDER_TEST_SUBSCRIPTION_SCOPED"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world_scoped", "description", "Say hello when any machine is provisionned"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world_scoped", "type", "RUNNABLE"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world_scoped", "runnable_type", "extensibility.abx"),
					//resource.TestCheckResourceAttr("aria_subscription.hello_world_scoped", "recover_runnable_type", ""),
					//resource.TestCheckResourceAttr("aria_subscription.hello_world_scoped", "recover_runnable_id", ""),
					resource.TestCheckResourceAttr("aria_subscription.hello_world_scoped", "event_topic_id", "compute.provision.post"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world_scoped", "blocking", "true"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world_scoped", "broadcast", "false"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world_scoped", "contextual", "false"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world_scoped", "criteria", ""),
					resource.TestCheckResourceAttr("aria_subscription.hello_world_scoped", "disabled", "true"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world_scoped", "priority", "10"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world_scoped", "system", "false"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world_scoped", "timeout", "0"),
					resource.TestCheckResourceAttrSet("aria_subscription.hello_world_scoped", "org_id"),
					resource.TestCheckResourceAttrSet("aria_subscription.hello_world_scoped", "owner_id"),
					resource.TestCheckResourceAttrSet("aria_subscription.hello_world_scoped", "subscriber_id"),
				),
			},
			// Update and Read testing
			{
				Config: `
variable "test_project_id" {
  description = "Project where to generate test resources."
  type        = string
}

variable "test_abx_action_id" {
  description = "ABX action to use for testing subscriptions."
  type        = string
}

locals {
  project_ids   = [var.test_project_id]
  subscriber_id = data.aria_catalog_type.abx_actions.created_by
}

data "aria_catalog_type" "abx_actions" {
  id = "com.vmw.abx.actions"
}

resource "aria_subscription" "hello_world" {
  name           = "ARIA_PROVIDER_TEST_SUBSCRIPTION"
  description    = "Say hello when a machine is provisionned"
  type           = "RUNNABLE"
  runnable_type  = "extensibility.abx"
  runnable_id    = "8a7480d38e535332018e857e0d4f3437"
  event_topic_id = "compute.provision.post"
  project_ids    = [] # All projects
  blocking       = false
  contextual     = false
  disabled       = true # Its safer
  timeout        = 60
  priority       = 10

  lifecycle {
    postcondition {
      condition     = can(regex("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}", self.id))
      error_message = "Identifier must be a valid UUID string, actual ${self.id}"
    }
    postcondition {
      condition     = self.runnable_id == var.test_abx_action_id
      error_message = "Runnable ID must be ${var.test_abx_action_id}, actual ${self.runnable_id}"
    }
    postcondition {
      condition     = self.subscriber_id == local.subscriber_id
      error_message = "Subscriber ID must be ${local.subscriber_id}, actual ${self.subscriber_id}"
    }
    postcondition {
      condition     = length(self.project_ids) == 0
      error_message = "Project IDs must be [], actual [${join(", ", self.project_ids)}]"
    }
  }
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_subscription.hello_world", "id"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "name", "ARIA_PROVIDER_TEST_SUBSCRIPTION"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "description", "Say hello when a machine is provisionned"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "type", "RUNNABLE"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "runnable_type", "extensibility.abx"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "event_topic_id", "compute.provision.post"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "blocking", "false"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "broadcast", "false"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "contextual", "false"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "criteria", ""),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "disabled", "true"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "priority", "10"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "system", "false"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "timeout", "60"),
					resource.TestCheckResourceAttrSet("aria_subscription.hello_world", "org_id"),
					resource.TestCheckResourceAttrSet("aria_subscription.hello_world", "owner_id"),
					resource.TestCheckResourceAttrSet("aria_subscription.hello_world", "subscriber_id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aria_subscription.hello_world",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
			// TODO Check https://developer.hashicorp.com/terraform/plugin/sdkv2/testing/acceptance-tests/testcase#checkdestroy
		},
	})
}

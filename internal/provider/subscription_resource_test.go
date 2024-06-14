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
data "aria_catalog_type" "abx_actions" {
  id = "com.vmw.abx.actions"
}

locals {
	subscriber_id = data.aria_catalog_type.abx_actions.created_by
}

resource "aria_subscription" "hello_world" {
  name           = "ARIA_PROVIDER_TEST_SUBSCRIPTION"
  description    = "Say hello when some machine is provisionned"
  type           = "RUNNABLE"
  runnable_type  = "extensibility.abx"
  runnable_id    = "8a7480d38e535332018e857e0d4f3437"
  event_topic_id = "compute.provision.post"
  blocking       = true
  contextual     = false
  disabled       = true # Its safer
  timeout        = 0
  priority       = 10
  # constraints

  lifecycle {
  	postcondition {
  		condition     = self.subscriber_id == local.subscriber_id
  		error_message = "Expected ${local.subscriber_id}, actual ${self.subscriber_id}"
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
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "runnable_id", "8a7480d38e535332018e857e0d4f3437"),
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
				),
			},
			// ImportState testing
			{
				ResourceName:      "aria_subscription.hello_world",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: `
resource "aria_subscription" "hello_world" {
  name           = "ARIA_PROVIDER_TEST_SUBSCRIPTION"
  description    = "Say hello when a machine is provisionned"
  type           = "RUNNABLE"
  runnable_type  = "extensibility.abx"
  runnable_id    = "8a7480d38e535332018e857e0d4f3437"
  event_topic_id = "compute.provision.post"
  blocking       = false
  contextual     = false
  disabled       = true # Its safer
  timeout        = 60
  priority       = 10
  # constraints
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_subscription.hello_world", "id"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "name", "ARIA_PROVIDER_TEST_SUBSCRIPTION"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "description", "Say hello when a machine is provisionned"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "type", "RUNNABLE"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "runnable_type", "extensibility.abx"),
					resource.TestCheckResourceAttr("aria_subscription.hello_world", "runnable_id", "8a7480d38e535332018e857e0d4f3437"),
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
			// Delete testing automatically occurs in TestCase
			// TODO Check https://developer.hashicorp.com/terraform/plugin/sdkv2/testing/acceptance-tests/testcase#checkdestroy
		},
	})
}

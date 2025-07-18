// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrchestratorWorkflowBrokerResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: `
resource "aria_orchestrator_category" "root" {
  name      = "TEST_ARIA_PROVIDER"
  type      = "WorkflowCategory"
  parent_id = ""
}

locals {
  input_forms = [
    {
      layout = {
        pages = []
      }
      schema = {}
    }
  ]
}

resource "aria_orchestrator_workflow" "test" {
  name        = "Test Workflow Sync"
  description = "Workflow generated by the acceptance tests of Aria provider."
  category_id = aria_orchestrator_category.root.id
  version     = "0.1.0"

  position = { x = 100, y = 50 }

  restart_mode            = 1 # resume
  resume_from_failed_mode = 0 # default

  attrib        = jsonencode([])
  presentation  = jsonencode({})
  workflow_item = jsonencode([])

  input_parameters  = []
  output_parameters = []

  input_forms = jsonencode(local.input_forms)

  force_delete = true
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_orchestrator_workflow.test", "id"),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_workflow.test", "name", "Test Workflow Sync",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_workflow.test", "description",
						"Workflow generated by the acceptance tests of Aria provider.",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_workflow.test", "wait_imported", "true",
					),
					resource.TestCheckResourceAttrSet(
						"aria_orchestrator_workflow.test",
						"integration.name",
					),
					resource.TestCheckResourceAttrSet(
						"aria_orchestrator_workflow.test",
						"integration.endpoint_configuration_link",
					),
					resource.TestCheckResourceAttrSet(
						"aria_orchestrator_workflow.test",
						"integration.endpoint_uri",
					),
				),
			},
			// ImportState testing
			// FIXME https://github.com/davidfischer-ch/terraform-provider-aria/issues/122
			/*{
				ResourceName:      "aria_orchestrator_workflow.test",
				ImportState:       true,
				ImportStateVerify: true,

				// Prevent diff on force_delete field
				ImportStateVerifyIgnore: []string{"force_delete"},
			},*/
			// Delete testing automatically occurs in TestCase
		},
	})
}

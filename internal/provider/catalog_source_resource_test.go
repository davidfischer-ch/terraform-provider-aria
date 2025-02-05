// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCatalogSourceResource(t *testing.T) {
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

resource "aria_orchestrator_workflow" "test" {
  name        = "Test Workflow for Catalog Source"
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

  input_forms = jsonencode([
    {
      layout = {
        pages = []
      }
      schema = {}
    }
  ])
}

data "aria_integration" "workflows" {
  type_id = "com.vmw.vro.workflow"
}

# Publish the workflow we manage
resource "aria_catalog_source" "test" {
  name    = "ARIA_PROVIDER_TEST_CATALOG_SOURCE"
  type_id = data.aria_integration.workflows.type_id

  config = {
    workflows = [
      {
        id          = aria_orchestrator_workflow.test.id
        name        = aria_orchestrator_workflow.test.name
        description = aria_orchestrator_workflow.test.description
        version     = aria_orchestrator_workflow.test.version
        integration = {
          name                        = data.aria_integration.workflows.name
          endpoint_configuration_link = data.aria_integration.workflows.endpoint_configuration_link
          endpoint_uri                = data.aria_integration.workflows.endpoint_uri
        }
      }
    ]
  }

  lifecycle {
    postcondition {
      condition     = self.config.workflows[0].id == aria_orchestrator_workflow.test.id
      error_message = "Oups workflow is missing or not the good one..."
    }
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_catalog_source.test", "id"),
					resource.TestCheckResourceAttr("aria_catalog_source.test", "name", "ARIA_PROVIDER_TEST_CATALOG_SOURCE"),
					resource.TestCheckResourceAttr("aria_catalog_source.test", "type_id", "com.vmw.vro.workflow"),
					resource.TestCheckResourceAttr("aria_catalog_source.test", "global", "true"),
					resource.TestCheckResourceAttrSet("aria_catalog_source.test", "created_at"),
					resource.TestCheckResourceAttrSet("aria_catalog_source.test", "created_by"),
					resource.TestCheckResourceAttrSet("aria_catalog_source.test", "last_updated_at"),
					resource.TestCheckResourceAttrSet("aria_catalog_source.test", "last_updated_by"),
					resource.TestCheckResourceAttrSet("aria_catalog_source.test", "last_import_started_at"),
					resource.TestCheckResourceAttrSet("aria_catalog_source.test", "last_import_completed_at"),
					resource.TestCheckResourceAttr("aria_catalog_source.test", "last_import_errors.#", "0"),
					resource.TestCheckResourceAttr("aria_catalog_source.test", "items_found", "1"),
					resource.TestCheckResourceAttr("aria_catalog_source.test", "items_imported", "1"),
					resource.TestCheckResourceAttr("aria_catalog_source.test", "wait_imported", "true"),
					resource.TestCheckResourceAttr("aria_catalog_source.test", "config.source_project_id", ""),
				),
			},

			// ImportState testing
			/* TODO https://github.com/davidfischer-ch/terraform-provider-aria/issues/111
			   {
			     ResourceName:      "aria_catalog_source.test",
			     ImportState:       true,
			     ImportStateVerify: true,
			   }, */
			// Delete testing automatically occurs in TestCase
			// TODO Check https://developer.hashicorp.com/terraform/plugin/sdkv2/testing/acceptance-tests/testcase#checkdestroy
		},
	})
}

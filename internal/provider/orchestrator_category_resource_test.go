// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrchestratorCategoryResource(t *testing.T) {
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

resource "aria_orchestrator_category" "a" {
  name      = "A"
  type      = "WorkflowCategory"
  parent_id = resource.aria_orchestrator_category.root.id
}

resource "aria_orchestrator_category" "b" {
  name      = "B"
  type      = "WorkflowCategory"
  parent_id = resource.aria_orchestrator_category.a.id
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_orchestrator_category.root", "id"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.root", "name", "TEST_ARIA_PROVIDER"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.root", "path", "TEST_ARIA_PROVIDER"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.root", "type", "WorkflowCategory"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.root", "parent_id", ""),

					resource.TestCheckResourceAttrSet("aria_orchestrator_category.a", "id"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.a", "name", "A"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.a", "path", "TEST_ARIA_PROVIDER/A"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.a", "type", "WorkflowCategory"),
					resource.TestCheckResourceAttrPair(
						"aria_orchestrator_category.a", "parent_id",
						"aria_orchestrator_category.root", "id",
					),

					resource.TestCheckResourceAttrSet("aria_orchestrator_category.b", "id"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.b", "name", "B"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.b", "path", "TEST_ARIA_PROVIDER/A/B"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.b", "type", "WorkflowCategory"),
					resource.TestCheckResourceAttrPair(
						"aria_orchestrator_category.b", "parent_id",
						"aria_orchestrator_category.a", "id",
					),
				),
			},
			// Update (change name) and Read testing
			{
				Config: `
resource "aria_orchestrator_category" "root" {
  name      = "TEST_ARIA_PROVIDER"
  type      = "WorkflowCategory"
  parent_id = ""
}

resource "aria_orchestrator_category" "a" {
  name      = "A"
  type      = "WorkflowCategory"
  parent_id = resource.aria_orchestrator_category.root.id
}

resource "aria_orchestrator_category" "b" {
  name      = "Bb"
  type      = "WorkflowCategory"
  parent_id = resource.aria_orchestrator_category.a.id
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_orchestrator_category.root", "id"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.root", "name", "TEST_ARIA_PROVIDER"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.root", "path", "TEST_ARIA_PROVIDER"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.root", "type", "WorkflowCategory"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.root", "parent_id", ""),

					resource.TestCheckResourceAttrSet("aria_orchestrator_category.a", "id"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.a", "name", "A"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.a", "path", "TEST_ARIA_PROVIDER/A"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.a", "type", "WorkflowCategory"),
					resource.TestCheckResourceAttrPair(
						"aria_orchestrator_category.a", "parent_id",
						"aria_orchestrator_category.root", "id",
					),

					resource.TestCheckResourceAttrSet("aria_orchestrator_category.b", "id"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.b", "name", "Bb"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.b", "path", "TEST_ARIA_PROVIDER/A/Bb"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.b", "type", "WorkflowCategory"),
					resource.TestCheckResourceAttrPair(
						"aria_orchestrator_category.b", "parent_id",
						"aria_orchestrator_category.a", "id",
					),
				),
			},
			// Update (change parent) and Read testing
			{
				Config: `
resource "aria_orchestrator_category" "root" {
  name      = "TEST_ARIA_PROVIDER"
  type      = "WorkflowCategory"
  parent_id = ""
}

resource "aria_orchestrator_category" "a" {
  name      = "A"
  type      = "WorkflowCategory"
  parent_id = resource.aria_orchestrator_category.root.id
}

resource "aria_orchestrator_category" "b" {
  name      = "Bb"
  type      = "WorkflowCategory"
  parent_id = resource.aria_orchestrator_category.root.id

  # Help's the cleanup phase to destroy resource in appropriate order in case of error
  depends_on = [resource.aria_orchestrator_category.a]
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_orchestrator_category.root", "id"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.root", "name", "TEST_ARIA_PROVIDER"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.root", "path", "TEST_ARIA_PROVIDER"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.root", "type", "WorkflowCategory"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.root", "parent_id", ""),

					resource.TestCheckResourceAttrSet("aria_orchestrator_category.a", "id"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.a", "name", "A"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.a", "path", "TEST_ARIA_PROVIDER/A"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.a", "type", "WorkflowCategory"),
					resource.TestCheckResourceAttrPair(
						"aria_orchestrator_category.a", "parent_id",
						"aria_orchestrator_category.root", "id",
					),

					resource.TestCheckResourceAttrSet("aria_orchestrator_category.b", "id"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.b", "name", "Bb"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.b", "path", "TEST_ARIA_PROVIDER/Bb"),
					resource.TestCheckResourceAttr("aria_orchestrator_category.b", "type", "WorkflowCategory"),
					resource.TestCheckResourceAttrPair(
						"aria_orchestrator_category.b", "parent_id",
						"aria_orchestrator_category.root", "id",
					),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aria_orchestrator_category.root",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "aria_orchestrator_category.a",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "aria_orchestrator_category.b",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

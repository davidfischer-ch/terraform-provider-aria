// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccCatalogItemIconResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: `
variable "test_catalog_item_id" {
	description = "Catalog item which icon will be manipulated."
  type        = string
}

resource "aria_icon" "test" {
  path = "../../tests/icon.png"
}

resource "aria_icon" "test_2" {
  path = "../../tests/icon.png"
}

resource "aria_catalog_item_icon" "test" {
	item_id = var.test_catalog_item_id
	icon_id = aria_icon.test.id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_catalog_item_icon.test", "item_id"),
					resource.TestCheckResourceAttrPair(
						"aria_catalog_item_icon.test", "icon_id",
						"aria_icon.test", "id",
					),
				),
			},
			// Delete duplicate Icon and Read testing
			{
				Config: `
variable "test_catalog_item_id" {
	description = "Catalog item which icon will be manipulated."
  type        = string
}

resource "aria_icon" "test" {
  path = "../../tests/icon.png"
}

resource "aria_catalog_item_icon" "test" {
	item_id = var.test_catalog_item_id
	icon_id = aria_icon.test.id
}
`,
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{ // Terraform don't think it has to be recreated (before refresh)
						plancheck.ExpectResourceAction("aria_icon.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{ // But icon has to be created again (after refresh)
						plancheck.ExpectResourceAction("aria_icon.test", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("aria_catalog_item_icon.test", plancheck.ResourceActionUpdate),
					},
				},
			},
			// Delete testing automatically occurs in TestCase
			// Check https://developer.hashicorp.com/terraform/plugin/sdkv2/testing/acceptance-tests/testcase#checkdestroy
		},
	})
}

func TestAccCatalogItemIconKeepOnDeleteResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: `
variable "test_catalog_item_id" {
	description = "Catalog item which icon will be manipulated."
  type        = string
}

resource "aria_icon" "test" {
  path = "../../tests/icon.png"
}

resource "aria_icon" "test_2" {
  path            = "../../tests/icon.png"
  keep_on_destroy = true
}

resource "aria_catalog_item_icon" "test" {
	item_id = var.test_catalog_item_id
	icon_id = aria_icon.test.id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_catalog_item_icon.test", "item_id"),
					resource.TestCheckResourceAttrPair(
						"aria_catalog_item_icon.test", "icon_id",
						"aria_icon.test", "id",
					),
				),
			},
			// "Soft" Delete duplicate Icon and Read testing
			{
				Config: `
variable "test_catalog_item_id" {
	description = "Catalog item which icon will be manipulated."
  type        = string
}

resource "aria_icon" "test" {
  path = "../../tests/icon.png"
}

resource "aria_catalog_item_icon" "test" {
	item_id = var.test_catalog_item_id
	icon_id = aria_icon.test.id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_catalog_item_icon.test", "item_id"),
					resource.TestCheckResourceAttrPair("aria_catalog_item_icon.test", "icon_id", "aria_icon.test", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
			// Check https://developer.hashicorp.com/terraform/plugin/sdkv2/testing/acceptance-tests/testcase#checkdestroy
		},
	})
}

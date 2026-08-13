// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomFormResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: `
variable "test_catalog_item_id" {
  description = "Catalog item which form will be manipulated."
  type        = string
}

variable "test_catalog_item_type" {
  description = "Catalog item which form will be manipulated."
  type        = string
}

resource "aria_custom_form" "test" {
  name        = "ARIA_PROVIDER_TEST_FORM"
  type        = "requestForm"
  form        = jsonencode({})
  styles      = ""
  source_id   = var.test_catalog_item_id
  source_type = var.test_catalog_item_type
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_custom_form.test", "id"),
					resource.TestCheckResourceAttr(
						"aria_custom_form.test", "name", "ARIA_PROVIDER_TEST_FORM",
					),
					resource.TestCheckResourceAttr("aria_custom_form.test", "type", "requestForm"),
				),
			},

			// Update and Read testing
			{
				Config: `
variable "test_catalog_item_id" {
  description = "Catalog item which form will be manipulated."
  type        = string
}

variable "test_catalog_item_type" {
  description = "Catalog item which form will be manipulated."
  type        = string
}

resource "aria_custom_form" "test" {
  name        = "ARIA_PROVIDER_TEST_FORM_RENAMED"
  type        = "requestForm"
  form        = jsonencode({})
  styles      = ""
  source_id   = var.test_catalog_item_id
  source_type = var.test_catalog_item_type
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_custom_form.test", "id"),
					resource.TestCheckResourceAttr(
						"aria_custom_form.test", "name", "ARIA_PROVIDER_TEST_FORM_RENAMED",
					),
					resource.TestCheckResourceAttr("aria_custom_form.test", "type", "requestForm"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aria_custom_form.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

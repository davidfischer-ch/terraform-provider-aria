// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCatalogItemDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read (by ID) testing
			{
				Config: `
variable "test_catalog_item_id" {
  description = "Catalog item which data will be retrieved."
  type        = string
}

data "aria_catalog_item" "test_a" {
  id = var.test_catalog_item_id

  lifecycle {
  	postcondition {
  		condition     = self.id == var.test_catalog_item_id
  		error_message = "Identifier must be ${var.test_catalog_item_id}, actual ${self.id}"
  	}
  }
}

// With a lot of criteria to optimize query
data "aria_catalog_item" "test_b" {
	name        = data.aria_catalog_item.test_a.name
	external_id = data.aria_catalog_item.test_a.external_id
	type_id     = data.aria_catalog_item.test_a.type_id
}

// No filtering criteria... All items are candidates
data "aria_catalog_item" "test_c" {
	external_id = data.aria_catalog_item.test_a.external_id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.aria_catalog_item.test_a", "name"),
					resource.TestCheckResourceAttrSet("data.aria_catalog_item.test_a", "description"),
					resource.TestCheckResourceAttrSet("data.aria_catalog_item.test_a", "schema"),
					resource.TestCheckResourceAttrSet("data.aria_catalog_item.test_a", "external_id"),
					//resource.TestCheckResourceAttrSet("data.aria_catalog_item.test_a", "form_id"),
					resource.TestCheckResourceAttrSet("data.aria_catalog_item.test_a", "icon_id"),
					resource.TestCheckResourceAttrSet("data.aria_catalog_item.test_a", "type_id"),
					//resource.TestCheckResourceAttrSet("data.aria_catalog_item.test_a", "source_id"),
					//resource.TestCheckResourceAttrSet("data.aria_catalog_item.test_a", "source_name"),
					resource.TestCheckResourceAttrSet("data.aria_catalog_item.test_a", "created_at"),
					resource.TestCheckResourceAttrSet("data.aria_catalog_item.test_a", "created_by"),
					resource.TestCheckResourceAttrSet("data.aria_catalog_item.test_a", "last_updated_at"),
					resource.TestCheckResourceAttrSet("data.aria_catalog_item.test_a", "last_updated_by"),

					// By External ID
					resource.TestCheckResourceAttrPair(
						"data.aria_catalog_item.test_a", "id",
						"data.aria_catalog_item.test_b", "id",
					),

					resource.TestCheckResourceAttrPair(
						"data.aria_catalog_item.test_a", "name",
						"data.aria_catalog_item.test_b", "name",
					),

					resource.TestCheckResourceAttrPair(
						"data.aria_catalog_item.test_a", "description",
						"data.aria_catalog_item.test_b", "description",
					),

					resource.TestCheckResourceAttrPair(
						"data.aria_catalog_item.test_a", "external_id",
						"data.aria_catalog_item.test_b", "external_id",
					),

					// By External ID
					resource.TestCheckResourceAttrPair(
						"data.aria_catalog_item.test_a", "id",
						"data.aria_catalog_item.test_c", "id",
					),

					resource.TestCheckResourceAttrPair(
						"data.aria_catalog_item.test_a", "name",
						"data.aria_catalog_item.test_c", "name",
					),

					resource.TestCheckResourceAttrPair(
						"data.aria_catalog_item.test_a", "description",
						"data.aria_catalog_item.test_c", "description",
					),

					resource.TestCheckResourceAttrPair(
						"data.aria_catalog_item.test_a", "external_id",
						"data.aria_catalog_item.test_c", "external_id",
					),
				),
			},
		},
	})
}

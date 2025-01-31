// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccTagResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: `
resource "aria_tag" "test" {
  key = "ARIA_PROVIDER_TEST_TAG"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_tag.test", "id"),
					resource.TestCheckResourceAttr("aria_tag.test", "key", "ARIA_PROVIDER_TEST_TAG"),
					resource.TestCheckResourceAttr("aria_tag.test", "value", ""),
					resource.TestCheckResourceAttr("aria_tag.test", "force_delete", "false"),
					resource.TestCheckResourceAttr("aria_tag.test", "keep_on_destroy", "false"),
				),
			},
			// Update (recreate) and Read testing
			{
				Config: `
resource "aria_tag" "test" {
	key   = "ARIA_PROVIDER_TEST_TAG"
	value = "newvalue"
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("aria_tag.test", plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_tag.test", "id"),
					resource.TestCheckResourceAttr("aria_tag.test", "key", "ARIA_PROVIDER_TEST_TAG"),
					resource.TestCheckResourceAttr("aria_tag.test", "value", "newvalue"),
					resource.TestCheckResourceAttr("aria_tag.test", "force_delete", "false"),
					resource.TestCheckResourceAttr("aria_tag.test", "keep_on_destroy", "false"),
				),
			},
			// Update (change flags) and Read testing
			{
				Config: `
resource "aria_tag" "test" {
	key   = "ARIA_PROVIDER_TEST_TAG"
	value = "newvalue"

	force_delete = true
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_tag.test", "id"),
					resource.TestCheckResourceAttr("aria_tag.test", "key", "ARIA_PROVIDER_TEST_TAG"),
					resource.TestCheckResourceAttr("aria_tag.test", "value", "newvalue"),
					resource.TestCheckResourceAttr("aria_tag.test", "force_delete", "true"),
					resource.TestCheckResourceAttr("aria_tag.test", "keep_on_destroy", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aria_tag.test",
				ImportState:       true,
				ImportStateVerify: true,

				// Prevent diff on force_delete field
				ImportStateVerifyIgnore: []string{"force_delete"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccIconResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: `
resource "aria_icon" "test" {
  path = "../../tests/icon.svg"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_icon.test", "id"),
					resource.TestCheckResourceAttr("aria_icon.test", "path", "../../tests/icon.svg"),
					resource.TestCheckResourceAttr("aria_icon.test", "hash", "9eb36dc3af8fe94b1814dd419bb5bc6405cac9cbd13e42af1bcc545dc8b69a0c"),
				),
			},
			// Update (recreate) and Read testing
			{
				Config: `
resource "aria_icon" "test" {
  path = "../../tests/icon.png"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_icon.test", "id"),
					resource.TestCheckResourceAttr("aria_icon.test", "path", "../../tests/icon.png"),
					resource.TestCheckResourceAttr("aria_icon.test", "hash", "0e6822039f0795d0d02f2660c25e68a5fd31446083541922b8a9336ccbc75943"),
				),
			},
			// No-op and Read testing
			{
				Config: `
resource "aria_icon" "test" {
  path = "../../tests/icon.png"
  hash = "0e6822039f0795d0d02f2660c25e68a5fd31446083541922b8a9336ccbc75943"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_icon.test", "id"),
					resource.TestCheckResourceAttr("aria_icon.test", "path", "../../tests/icon.png"),
					resource.TestCheckResourceAttr("aria_icon.test", "hash", "0e6822039f0795d0d02f2660c25e68a5fd31446083541922b8a9336ccbc75943"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Delete testing automatically occurs in TestCase
			// Check https://developer.hashicorp.com/terraform/plugin/sdkv2/testing/acceptance-tests/testcase#checkdestroy
		},
	})
}

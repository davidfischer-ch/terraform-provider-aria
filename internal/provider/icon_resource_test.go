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
					resource.TestCheckResourceAttr("aria_icon.test", "hash", "1ef3f922f7072f1b2326f538f411cbe4d121a0ed50a308716f1f229628ae7d60"),
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

func TestAccIconConcurrentSameContentResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create/delete same icon multiple times in "parallel" (protected by mutex) and Read testing
			{
				Config: `
resource "aria_icon" "test" {
  path = "../../tests/icon.svg"
}

resource "aria_icon" "test_others" {
	count = 5
  path  = "../../tests/icon.svg"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_icon.test", "id"),
					resource.TestCheckResourceAttr("aria_icon.test", "path", "../../tests/icon.svg"),
					resource.TestCheckResourceAttr("aria_icon.test", "hash", "1ef3f922f7072f1b2326f538f411cbe4d121a0ed50a308716f1f229628ae7d60"),
				),
			},
			// Delete duplicated copies but keep one and Read testing
			{
				Config: `
resource "aria_icon" "test" {
  path = "../../tests/icon.svg"
}
`,
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{ // Terraform don't think it has to be recreated (before refresh)
						plancheck.ExpectResourceAction("aria_icon.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{ // But icon has to be created again (after refresh)
						plancheck.ExpectResourceAction("aria_icon.test", plancheck.ResourceActionCreate),
					},
				},
			},
			// Delete testing automatically occurs in TestCase
			// Check https://developer.hashicorp.com/terraform/plugin/sdkv2/testing/acceptance-tests/testcase#checkdestroy
		},
	})
}

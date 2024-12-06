// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"regexp"
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
					resource.TestMatchResourceAttr("aria_icon.test", "hash", regexp.MustCompile("[0-9a-f]{64}")),
					resource.TestCheckResourceAttr("aria_icon.test", "keep_on_destroy", "false"),
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
					resource.TestCheckResourceAttr("aria_icon.test", "hash", "724d45fec592788dcaca7526cfbb68e0867adb48ed1a9f8d6f5a6fde094bcf7d"),
					resource.TestCheckResourceAttr("aria_icon.test", "keep_on_destroy", "false"),
				),
			},
			// No-op and Read testing
			{
				Config: `
resource "aria_icon" "test" {
  path = "../../tests/icon.png"
  hash = "724d45fec592788dcaca7526cfbb68e0867adb48ed1a9f8d6f5a6fde094bcf7d"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_icon.test", "id"),
					resource.TestCheckResourceAttr("aria_icon.test", "path", "../../tests/icon.png"),
					resource.TestCheckResourceAttr("aria_icon.test", "hash", "724d45fec592788dcaca7526cfbb68e0867adb48ed1a9f8d6f5a6fde094bcf7d"),
					resource.TestCheckResourceAttr("aria_icon.test", "keep_on_destroy", "false"),
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
					resource.TestMatchResourceAttr("aria_icon.test", "hash", regexp.MustCompile("[0-9a-f]{64}")),
					resource.TestCheckResourceAttr("aria_icon.test", "keep_on_destroy", "false"),
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

func TestAccIconKeepOnDestroyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create/delete same icon multiple times in "parallel" (protected by mutex) and Read testing
			{
				Config: `
resource "aria_icon" "test" {
  path = "../../tests/icon.png"
}

resource "aria_icon" "test_other" {
  path            = "../../tests/icon.png"
  keep_on_destroy = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_icon.test", "id"),
					resource.TestCheckResourceAttr("aria_icon.test", "path", "../../tests/icon.png"),
					resource.TestMatchResourceAttr("aria_icon.test", "hash", regexp.MustCompile("[0-9a-f]{64}")),
					resource.TestCheckResourceAttr("aria_icon.test", "keep_on_destroy", "false"),
					resource.TestCheckResourceAttrPair("aria_icon.test", "id", "aria_icon.test_other", "id"),
					resource.TestCheckResourceAttr("aria_icon.test_other", "keep_on_destroy", "true"),
				),
			},
			// "Soft" Delete duplicated copy but keep one and Read testing
			{
				Config: `
resource "aria_icon" "test" {
  path = "../../tests/icon.png"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_icon.test", "id"),
					resource.TestCheckResourceAttr("aria_icon.test", "path", "../../tests/icon.png"),
					resource.TestMatchResourceAttr("aria_icon.test", "hash", regexp.MustCompile("[0-9a-f]{64}")),
					resource.TestCheckResourceAttr("aria_icon.test", "keep_on_destroy", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
			// Check https://developer.hashicorp.com/terraform/plugin/sdkv2/testing/acceptance-tests/testcase#checkdestroy
		},
	})
}

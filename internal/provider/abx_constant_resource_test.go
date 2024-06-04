// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccABXConstantResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: `
resource "aria_abx_constant" "test" {
  name  = "ARIA_PROVIDER_TEST_CONSTANT"
  value = "42"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_abx_constant.test", "id"),
					resource.TestCheckResourceAttrSet("aria_abx_constant.test", "org_id"),
					resource.TestCheckResourceAttr("aria_abx_constant.test", "name", "ARIA_PROVIDER_TEST_CONSTANT"),
					resource.TestCheckResourceAttr("aria_abx_constant.test", "value", "42"),
					resource.TestCheckResourceAttr("aria_abx_constant.test", "encrypted", "false"),
				),
			},
            // ImportState testing
            {
                ResourceName:      "aria_abx_constant.test",
                ImportState:       true,
                ImportStateVerify: true,
            },
			// Update and Read testing
			{
				Config: `
resource "aria_abx_constant" "test" {
  name  = "ARIA_PROVIDER_TEST_CONSTANT"
  value = "newvalue"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_abx_constant.test", "id"),
					resource.TestCheckResourceAttrSet("aria_abx_constant.test", "org_id"),
					resource.TestCheckResourceAttr("aria_abx_constant.test", "name", "ARIA_PROVIDER_TEST_CONSTANT"),
					resource.TestCheckResourceAttr("aria_abx_constant.test", "value", "newvalue"),
					resource.TestCheckResourceAttr("aria_abx_constant.test", "encrypted", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
			// TODO Check https://developer.hashicorp.com/terraform/plugin/sdkv2/testing/acceptance-tests/testcase#checkdestroy
		},
	})
}

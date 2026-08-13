// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccABXSensitiveConstantResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: `
resource "aria_abx_sensitive_constant" "test" {
  name  = "ARIA_PROVIDER_TEST_SENSITIVE_CONSTANT"
  value = "pass1234"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_abx_sensitive_constant.test", "id"),
					resource.TestCheckResourceAttrSet("aria_abx_sensitive_constant.test", "org_id"),
					resource.TestCheckResourceAttr("aria_abx_sensitive_constant.test", "name", "ARIA_PROVIDER_TEST_SENSITIVE_CONSTANT"),
					resource.TestCheckResourceAttr("aria_abx_sensitive_constant.test", "value", "pass1234"),
					resource.TestCheckResourceAttr("aria_abx_sensitive_constant.test", "encrypted", "true"),
				),
			},
			// Update and Read testing
			{
				Config: `
resource "aria_abx_sensitive_constant" "test" {
  name  = "ARIA_PROVIDER_TEST_SENSITIVE_CONSTANT"
  value = "newvalue"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_abx_sensitive_constant.test", "id"),
					resource.TestCheckResourceAttrSet("aria_abx_sensitive_constant.test", "org_id"),
					resource.TestCheckResourceAttr("aria_abx_sensitive_constant.test", "name", "ARIA_PROVIDER_TEST_SENSITIVE_CONSTANT"),
					resource.TestCheckResourceAttr("aria_abx_sensitive_constant.test", "value", "newvalue"),
					resource.TestCheckResourceAttr("aria_abx_sensitive_constant.test", "encrypted", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccABXSecretResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: `
resource "aria_abx_secret" "test" {
  name  = "ARIA_PROVIDER_TEST_SECRET"
  value = "pass1234"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_abx_secret.test", "id"),
					resource.TestCheckResourceAttrSet("aria_abx_secret.test", "org_id"),
					resource.TestCheckResourceAttr("aria_abx_secret.test", "name", "ARIA_PROVIDER_TEST_SECRET"),
					resource.TestCheckResourceAttr("aria_abx_secret.test", "value", "pass1234"),
					resource.TestCheckResourceAttr("aria_abx_secret.test", "encrypted", "true"),
				),
			},
			// Update and Read testing
			{
				Config: `
resource "aria_abx_secret" "test" {
  name  = "ARIA_PROVIDER_TEST_SECRET"
  value = "newvalue"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_abx_secret.test", "id"),
					resource.TestCheckResourceAttrSet("aria_abx_secret.test", "org_id"),
					resource.TestCheckResourceAttr("aria_abx_secret.test", "name", "ARIA_PROVIDER_TEST_SECRET"),
					resource.TestCheckResourceAttr("aria_abx_secret.test", "value", "newvalue"),
					resource.TestCheckResourceAttr("aria_abx_secret.test", "encrypted", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
			// TODO Check https://developer.hashicorp.com/terraform/plugin/sdkv2/testing/acceptance-tests/testcase#checkdestroy
		},
	})
}

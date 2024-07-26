// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: `
resource "aria_project" "test" {
  name              = "ARIA_PROVIDER_TEST_PROJECT"
  operation_timeout = 0
  shared_resources  = true

  constraints = {}

  properties = {
    toto = "tata"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_project.test", "id"),
					resource.TestCheckResourceAttr("aria_project.test", "name", "ARIA_PROVIDER_TEST_PROJECT"),
					resource.TestCheckResourceAttr("aria_project.test", "operation_timeout", "0"),
					resource.TestCheckResourceAttr("aria_project.test", "shared_resources", "true"),
					resource.TestCheckResourceAttr("aria_project.test", "properties.toto", "tata"),
					resource.TestCheckResourceAttrSet("aria_project.test", "org_id"),
				),
			},
			// Update testing
			// TODO
			// ImportState testing
			/*{
				ResourceName:      "aria_project.test",
				ImportState:       true,
				ImportStateVerify: true,
			},*/
			// Delete testing automatically occurs in TestCase
			// TODO Check https://developer.hashicorp.com/terraform/plugin/sdkv2/testing/acceptance-tests/testcase#checkdestroy
		},
	})
}

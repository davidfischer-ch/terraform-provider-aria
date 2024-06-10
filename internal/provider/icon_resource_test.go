// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIconResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIconResourceConfig(svgIcon),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_icon.test", "id"),
					resource.TestCheckResourceAttr("aria_icon.test", "content", svgIcon),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aria_icon.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			// TODO Implement this test
			/* {
				Config: testAccIconResourceConfig(bSvgIcon),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aria_icon.test", "content", bSvgIcon),
				),
			}, */
			// Delete testing automatically occurs in TestCase
			// TODO Check https://developer.hashicorp.com/terraform/plugin/sdkv2/testing/acceptance-tests/testcase#checkdestroy
		},
	})
}

func testAccIconResourceConfig(content string) string {
	return fmt.Sprintf(`
resource "aria_icon" "test" {
  content = %[1]q
}
`, content)
}

const svgIcon = `<svg width="24" height="24" xmlns="http://www.w3.org/2000/svg" fill-rule="evenodd" clip-rule="evenodd"><path d="M12 0c6.623 0 12 5.377 12 12s-5.377 12-12 12-12-5.377-12-12 5.377-12 12-12zm0 1c6.071 0 11 4.929 11 11s-4.929 11-11 11-11-4.929-11-11 4.929-11 11-11z"/></svg>`

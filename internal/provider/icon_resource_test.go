// Copyright (c) HashiCorp, Inc.
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
				Config: testAccIconResourceConfig(aSvgIcon),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_icon.test", "id"),
					resource.TestCheckResourceAttr("aria_icon.test", "content", aSvgIcon),
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

func testAccIconResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "aria_icon" "test" {
  content = %[1]q
}
`, configurableAttribute)
}

const aSvgIcon = `<svg width="24" height="24" xmlns="http://www.w3.org/2000/svg" fill-rule="evenodd" clip-rule="evenodd"><path d="M12 0c6.623 0 12 5.377 12 12s-5.377 12-12 12-12-5.377-12-12 5.377-12 12-12zm0 1c6.071 0 11 4.929 11 11s-4.929 11-11 11-11-4.929-11-11 4.929-11 11-11z"/></svg>`
const bSvgIcon = `
<?xml version="1.0" encoding="iso-8859-1"?>
<!-- Uploaded to: SVG Repo, www.svgrepo.com, Generator: SVG Repo Mixer Tools -->
<svg fill="#000000" height="800px" width="800px" version="1.1" id="Layer_1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"
	 viewBox="0 0 330 330" xml:space="preserve">
<path id="XMLID_523_" d="M315,0H15C6.716,0,0,6.716,0,15v300c0,8.284,6.716,15,15,15h300c8.284,0,15-6.716,15-15V15
	C330,6.716,323.285,0,315,0z M300,300H30V30h270V300z"/>
</svg>
`

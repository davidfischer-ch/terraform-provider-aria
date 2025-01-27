// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIntegrationDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: `
data "aria_integration" "test" {
  type_id = "com.vmw.vro.workflow"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.aria_integration.test", "type_id",
						"com.vmw.vro.workflow",
					),
					resource.TestCheckResourceAttr(
						"data.aria_integration.test", "name",
						"embedded-VRO",
					),
					resource.TestMatchResourceAttr(
						"data.aria_integration.test", "endpoint_configuration_link",
						regexp.MustCompile(`^/resources/endpoints/[0-9a-f]{8}-([0-9a-f]{4}-){3}[0-9a-f]{12}$`),
					),
					resource.TestMatchResourceAttr(
						"data.aria_integration.test", "endpoint_uri",
						regexp.MustCompile(`^https://\S+$`),
					),
				),
			},
		},
	})
}

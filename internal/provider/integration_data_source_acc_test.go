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

// The name is what disambiguates a tenant exposing several integrations of the type
data "aria_integration" "test_by_name" {
  type_id = "com.vmw.vro.workflow"
  name    = data.aria_integration.test.name
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.aria_integration.test", "type_id",
						"com.vmw.vro.workflow",
					),
					resource.TestCheckResourceAttrSet("data.aria_integration.test", "name"),
					resource.TestCheckResourceAttrPair(
						"data.aria_integration.test", "name",
						"data.aria_integration.test_by_name", "name",
					),
					resource.TestCheckResourceAttrPair(
						"data.aria_integration.test", "endpoint_configuration_link",
						"data.aria_integration.test_by_name", "endpoint_configuration_link",
					),
					resource.TestCheckResourceAttrPair(
						"data.aria_integration.test", "endpoint_uri",
						"data.aria_integration.test_by_name", "endpoint_uri",
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

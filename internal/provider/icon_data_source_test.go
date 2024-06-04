// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIconDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccIconDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.aria_icon.test", "id", "72a9a2c7-494e-31d7-afe8-cd27479c407e"),
				),
			},
		},
	})
}

const testAccIconDataSourceConfig = `
data "aria_icon" "test" {
  id = "72a9a2c7-494e-31d7-afe8-cd27479c407e"
}
`

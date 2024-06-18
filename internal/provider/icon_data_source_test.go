// Copyright (c) State of Geneva (Switzerland)
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
				Config: `
variable "test_icon_id" {
	description = "ABX action to use for testing subscriptions."
  type        = string
}

data "aria_icon" "test" {
  id = var.test_icon_id

  lifecycle {
  	postcondition {
  		condition     = self.id == var.test_icon_id
  		error_message = "Identifier must be ${var.test_icon_id}, actual ${self.id}"
  	}
  }
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.aria_icon.test", "content"),
				),
			},
		},
	})
}

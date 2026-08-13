// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSecretDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: `
variable "test_secret_id" {
  description = "Secret to use for testing the data source."
  type        = string
}

data "aria_secret" "secret" {
  id = var.test_secret_id
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.aria_secret.secret", "id"),
					resource.TestCheckResourceAttrSet("data.aria_secret.secret", "name"),
					resource.TestCheckResourceAttrSet("data.aria_secret.secret", "description"),
					resource.TestCheckResourceAttrSet("data.aria_secret.secret", "org_id"),
					resource.TestCheckResourceAttrSet("data.aria_secret.secret", "org_scoped"),
					// resource.TestCheckResourceAttrSet("data.aria_secret.secret", "project_ids"),
					resource.TestCheckResourceAttrSet("data.aria_secret.secret", "created_at"),
					resource.TestCheckResourceAttrSet("data.aria_secret.secret", "created_by"),
					resource.TestCheckResourceAttrSet("data.aria_secret.secret", "updated_at"),
					resource.TestCheckResourceAttrSet("data.aria_secret.secret", "updated_by"),
				),
			},
		},
	})
}

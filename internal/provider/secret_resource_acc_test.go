// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccSecretResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		// Write-only attributes need Terraform 1.11 or later
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_11_0),
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: `
resource "aria_secret" "test" {
  name       = "ARIA_PROVIDER_TEST_SECRET"
  value      = "a"
  org_scoped = true
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_secret.test", "id"),
					resource.TestCheckResourceAttrSet("aria_secret.test", "org_id"),
					resource.TestCheckResourceAttrSet("aria_secret.test", "updated_at"),
					resource.TestCheckResourceAttr("aria_secret.test", "name", "ARIA_PROVIDER_TEST_SECRET"),
					resource.TestCheckResourceAttr("aria_secret.test", "description", ""),
					resource.TestCheckResourceAttr("aria_secret.test", "org_scoped", "true"),
					resource.TestCheckResourceAttr("aria_secret.test", "project_ids.#", "0"),
					// Write-only: never stored, only its checksum is
					resource.TestCheckNoResourceAttr("aria_secret.test", "value"),
					resource.TestCheckResourceAttr(
						"aria_secret.test", "value_checksum", SecretValueChecksum("a")),
				),
			},
			// No-op testing (neither the platform's masked value nor its date-time can drift us)
			{
				Config: `
resource "aria_secret" "test" {
  name       = "ARIA_PROVIDER_TEST_SECRET"
  value      = "a"
  org_scoped = true
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Update the value only, the plan must propose it
			{
				Config: `
resource "aria_secret" "test" {
  name       = "ARIA_PROVIDER_TEST_SECRET"
  value      = "b"
  org_scoped = true
}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("aria_secret.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aria_secret.test", "org_scoped", "true"),
					resource.TestCheckResourceAttr(
						"aria_secret.test", "value_checksum", SecretValueChecksum("b")),
				),
			},
			// ImportState testing, the value (and its checksum) cannot be known
			{
				ResourceName:            "aria_secret.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"value", "value_checksum"},
			},
			// Update the description, the secret stays organization scoped
			{
				Config: `
resource "aria_secret" "test" {
  name        = "ARIA_PROVIDER_TEST_SECRET"
  description = "Secret managed by Aria provider's acceptance tests."
  value       = "b"
  org_scoped  = true
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"aria_secret.test", "description", "Secret managed by Aria provider's acceptance tests."),
					resource.TestCheckResourceAttr("aria_secret.test", "org_scoped", "true"),
					resource.TestCheckResourceAttr(
						"aria_secret.test", "value_checksum", SecretValueChecksum("b")),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

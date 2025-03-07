// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrchestratorEnvironmentRepositoryResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: `
resource "aria_orchestrator_environment_repository" "test" {
	name     = "TEST_ARIA_PROVIDER"
	runtime  = "python:3.10"
	location = "https://your-registry.your-company.net/repository/pypi-all/simple"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"aria_orchestrator_environment_repository.test", "id",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_environment_repository.test", "name",
						"TEST_ARIA_PROVIDER",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_environment_repository.test", "runtime", "python:3.10",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_environment_repository.test", "location",
						"https://your-registry.your-company.net/repository/pypi-all/simple",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_environment_repository.test", "basic_auth", "false",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_environment_repository.test", "system_user", "",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_environment_repository.test", "system_credentials", "",
					),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aria_orchestrator_environment_repository.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: `
resource "aria_orchestrator_environment_repository" "test" {
	name               = "TEST_ARIA_PROVIDER"
	runtime            = "python:3.10"
	location           = "https://your-registry.your-company.net/repository/pypi-all/simple"
	system_user        = "toto"
	system_credentials = "tata"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"aria_orchestrator_environment_repository.test", "id",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_environment_repository.test", "name",
						"TEST_ARIA_PROVIDER",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_environment_repository.test", "runtime", "python:3.10",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_environment_repository.test", "location",
						"https://your-registry.your-company.net/repository/pypi-all/simple",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_environment_repository.test", "basic_auth", "true",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_environment_repository.test", "system_user", "toto",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_environment_repository.test", "system_credentials",
						"tata",
					),
				),
			},
			// Update and Read testing
			{
				Config: `
resource "aria_orchestrator_environment_repository" "test" {
	name     = "TEST_ARIA_PROVIDER_RENAMED"
	runtime  = "python:3.10"
	location = "https://your-registry.your-company.net/repository/pypi-all/other"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"aria_orchestrator_environment_repository.test", "id",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_environment_repository.test", "name",
						"TEST_ARIA_PROVIDER_RENAMED",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_environment_repository.test", "runtime", "python:3.10",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_environment_repository.test", "location",
						"https://your-registry.your-company.net/repository/pypi-all/other",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_environment_repository.test", "basic_auth", "false",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_environment_repository.test", "system_user", "",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_environment_repository.test", "system_credentials", "",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

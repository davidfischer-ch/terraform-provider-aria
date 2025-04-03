// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrchestratorActionDeleteConvergeResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create interlinked resources that Terraform doesn't know they are, three levels deep
			{
				Config: `
# The way the configuration is written Terraform is not aware that test_b depends on test_a, and
# test_c depends on test_b. Destroy should normally do not work properly if deletion wasn't retried.
# Automatically in case of conflits.

resource "aria_orchestrator_action" "test_a" {
  name                 = "actionA"
  module               = "aria_provider_tests"
  fqn                  = "aria_provider_tests/actionA"
  description          = "An action used by actionB."
  version              = "1.0.0"
  runtime              = "" # javascript
  runtime_memory_limit = 0
  runtime_timeout      = 0
  script               = ""
  input_parameters     = []
  output_type          = "Any"
}

resource "aria_orchestrator_action" "test_b" {
  name                 = "actionB"
  module               = "aria_provider_tests"
  fqn                  = "aria_provider_tests/actionB"
  description          = "An action using actionA."
  version              = "1.0.0"
  runtime              = "" # javascript
  runtime_memory_limit = 0
  runtime_timeout      = 0
  script               = <<EOT
 var actionA = System.getModule("aria_provider_tests").actionA();
 EOT
  input_parameters     = []
  output_type          = "Any"
}

resource "aria_orchestrator_action" "test_c" {
  name                 = "actionC"
  module               = "aria_provider_tests"
  fqn                  = "aria_provider_tests/actionC"
  description          = "An action using actionB."
  version              = "1.0.0"
  runtime              = "" # javascript
  runtime_memory_limit = 0
  runtime_timeout      = 0
  script               = <<EOT
 var actionA = System.getModule("aria_provider_tests").actionB();
 EOT
  input_parameters     = []
  output_type          = "Any"
}
`,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccOrchestratorActionForceDeleteResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create interlinked resources that Terraform doesn't know they are
			{
				Config: `
# The way the configuration is written Terraform is not aware that test_e depends on test_d.
# Destroy of test_d should not be possible unless test_e is destroyed too.
# In this case test_d's force_delete is true so it will be forced.

resource "aria_orchestrator_action" "test_d" {
  name                 = "actionD"
  module               = "aria_provider_tests"
  fqn                  = "aria_provider_tests/actionD"
  description          = "An action used by actionE."
  version              = "1.0.0"
  runtime              = "" # javascript
  runtime_memory_limit = 0
  runtime_timeout      = 0
  script               = ""
  input_parameters     = []
  output_type          = "Any"
  force_delete         = true
}

resource "aria_orchestrator_action" "test_e" {
  name                 = "actionE"
  module               = "aria_provider_tests"
  fqn                  = "aria_provider_tests/actionE"
  description          = "An action using actionD."
  version              = "1.0.0"
  runtime              = "" # javascript
  runtime_memory_limit = 0
  runtime_timeout      = 0
  script               = <<EOT
 var actionA = System.getModule("aria_provider_tests").actionD();
 EOT
  input_parameters     = []
  output_type          = "Any"
}
`,
			},
			// Destroy test_d shouldn't be possible if not forced. Here it is...
			{
				Config: `
resource "aria_orchestrator_action" "test_e" {
  name                 = "actionE"
  module               = "aria_provider_tests"
  fqn                  = "aria_provider_tests/actionE"
  description          = "An action using actionD."
  version              = "1.0.0"
  runtime              = "" # javascript
  runtime_memory_limit = 0
  runtime_timeout      = 0
  script               = <<EOT
 var actionA = System.getModule("aria_provider_tests").actionD();
 EOT
  input_parameters     = []
  output_type          = "Any"
}
`,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

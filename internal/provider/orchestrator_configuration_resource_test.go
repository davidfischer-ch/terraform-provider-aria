// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrchestratorConfigurationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: `
resource "aria_orchestrator_category" "root" {
  name      = "TEST_ARIA_PROVIDER"
  type      = "ConfigurationElementCategory"
  parent_id = ""
}

resource "aria_orchestrator_configuration" "test" {
  name        = "Test Config"
  description = "Config generated by the acceptance tests of Aria provider."
  category_id = aria_orchestrator_category.root.id
  version     = "0.0.0"

  attributes = []

  force_delete = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_orchestrator_configuration.test", "id"),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "name",
						"Test Config",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "description",
						"Config generated by the acceptance tests of Aria provider.",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "version",
						"0.0.0",
					),
					resource.TestMatchResourceAttr(
						"aria_orchestrator_configuration.test", "version_id",
						regexp.MustCompile("[0-9a-f]{40}"),
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.#",
						"0",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "force_delete",
						"true",
					),
					resource.TestCheckResourceAttrPair(
						"aria_orchestrator_configuration.test", "category_id",
						"aria_orchestrator_category.root", "id",
					),
				),
			},
			// Update (attributes, description, version) and Read testing
			{
				Config: `
resource "aria_orchestrator_category" "root" {
  name      = "TEST_ARIA_PROVIDER"
  type      = "ConfigurationElementCategory"
  parent_id = ""
}

resource "aria_orchestrator_configuration" "test" {
  name        = "Test Config"
  description = "Config generated by the acceptance tests of Aria provider (updated)."
  category_id = aria_orchestrator_category.root.id
  version     = "0.1.0"

  attributes = [
    {
      name        = "someString"
      description = "Some string value"
      type        = "string"
      value = {
        string = {
          value = "some value"
        }
      }
    },
    {
      name        = "someBoolean"
      description = "Some bool value"
      type        = "boolean"
      value = {
        boolean = {
          value = true
        }
      }
    },
    {
      name        = "someInteger"
      description = "Some integer value"
      type        = "number"
      value = {
        number = {
          value = 42
        }
      }
    },
    {
      name        = "someFloat"
      description = "Some float value"
      type        = "number"
      value = {
        number = {
          value = 3.141592
        }
      }
    },
    {
      name        = "restServer"
      description = "Some REST Host"
      type        = "REST:RESTHost"
      value = {
        sdk_object = {
          id   = "08bb4b24-2f8e-4d4a-ba6f-07c8aa7b3c2d"
          type = "REST:RESTHost"
        }
      }
    },
    {
      name        = "someArrayOfString"
      description = "Some array of string"
      type        = "Array/string"
      value = {
        array = {
          elements = [
            {
              string = {
                value = "foo"
              }
            },
            {
              string = {
                value = "bar"
              }
            }
          ]
        }
      }
    },
    {
      name        = "someEmptyArray"
      description = "Some array with no elements."
      type        = "Array/REST:RESTHost"
      value = {
        array = {
          elements = []
        }
      }
    }
    /*
    This is not yet well handled (mutated by the platform), we have to find a pattern for this.
    We cannot (?) expose 1-1 the API to make it declarative, we have to tackle this challenge.
    {
      name        = "someSecureString"
      description = "Some secure string value"
      type        = "SecureString"
      value = {
        secure_string = {
          value         = "test"
          is_plain_text = true
        }
      }
    }
    */
  ]

  force_delete = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_orchestrator_configuration.test", "id"),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "name",
						"Test Config",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "description",
						"Config generated by the acceptance tests of Aria provider (updated).",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "version",
						"0.1.0",
					),
					resource.TestMatchResourceAttr(
						"aria_orchestrator_configuration.test", "version_id",
						regexp.MustCompile("[0-9a-f]{40}"),
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.#",
						"7",
					),

					// String
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.0.name",
						"someString",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.0.description",
						"Some string value",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.0.type",
						"string",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.0.value.string.value",
						"some value",
					),

					// Boolean
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.1.name",
						"someBoolean",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.1.description",
						"Some bool value",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.1.type",
						"boolean",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.1.value.boolean.value",
						"true",
					),

					// Number (Integer)
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.2.name",
						"someInteger",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.2.description",
						"Some integer value",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.2.type",
						"number",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.2.value.number.value",
						"42",
					),

					// Number (Float)
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.3.name",
						"someFloat",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.3.description",
						"Some float value",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.3.type",
						"number",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.3.value.number.value",
						"3.141592",
					),

					// SDK Object
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.4.name",
						"restServer",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.4.description",
						"Some REST Host",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.4.type",
						"REST:RESTHost",
					),
					resource.TestCheckResourceAttrSet(
						"aria_orchestrator_configuration.test", "attributes.4.value.sdk_object.id",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.4.value.sdk_object.type",
						"REST:RESTHost",
					),

					// Array of Strings
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.5.name",
						"someArrayOfString",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.5.description",
						"Some array of string",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.5.type",
						"Array/string",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.5.value.array.elements.#",
						"2",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test",
						"attributes.5.value.array.elements.0.string.value",
						"foo",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test",
						"attributes.5.value.array.elements.1.string.value",
						"bar",
					),

					// Array (Empty)
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.6.value.array.elements.#",
						"0",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "force_delete",
						"true",
					),
					resource.TestCheckResourceAttrPair(
						"aria_orchestrator_configuration.test", "category_id",
						"aria_orchestrator_category.root", "id",
					),
				),
			},
			// Update (name & attributes) and Read testing
			{
				Config: `
resource "aria_orchestrator_category" "root" {
  name      = "TEST_ARIA_PROVIDER"
  type      = "ConfigurationElementCategory"
  parent_id = ""
}

resource "aria_orchestrator_configuration" "test" {
  name        = "Test Config Renamed"
  description = "Config generated by the acceptance tests of Aria provider (updated)."
  category_id = aria_orchestrator_category.root.id
  version     = "0.1.0"

  attributes = [
    {
      name        = "someString"
      description = "Some string value"
      type        = "string"
      value = {
        string = {
          value = "some value"
        }
      }
    },
    {
      name        = "someBoolean"
      description = "Some boolean value"
      type        = "boolean"
      value = {
        boolean = {
          value = true
        }
      }
    }
  ]

  force_delete = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_orchestrator_configuration.test", "id"),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "name",
						"Test Config Renamed",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "description",
						"Config generated by the acceptance tests of Aria provider (updated).",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "version",
						"0.1.0",
					),
					resource.TestMatchResourceAttr(
						"aria_orchestrator_configuration.test", "version_id",
						regexp.MustCompile("[0-9a-f]{40}"),
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.#",
						"2",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.0.name",
						"someString",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.0.description",
						"Some string value",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.0.type",
						"string",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.0.value.string.value",
						"some value",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.1.name",
						"someBoolean",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.1.description",
						"Some boolean value",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.1.type",
						"boolean",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.1.value.boolean.value",
						"true",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "force_delete",
						"true",
					),
					resource.TestCheckResourceAttrPair(
						"aria_orchestrator_configuration.test", "category_id",
						"aria_orchestrator_category.root", "id",
					),
				),
			},
			// Update (category) and Read testing
			{
				Config: `
resource "aria_orchestrator_category" "root" {
  name      = "TEST_ARIA_PROVIDER"
  type      = "ConfigurationElementCategory"
  parent_id = ""
}

resource "aria_orchestrator_category" "misc" {
  name      = "Misc"
  type      = "ConfigurationElementCategory"
  parent_id = aria_orchestrator_category.root.id
}

resource "aria_orchestrator_configuration" "test" {
  name        = "Test Config Renamed"
  description = "Config generated by the acceptance tests of Aria provider (updated)."
  category_id = aria_orchestrator_category.misc.id
  version     = "0.1.0"

  attributes = [
    {
      name        = "someString"
      description = "Some string value"
      type        = "string"
      value = {
        string = {
          value = "some value"
        }
      }
    },
    {
      name        = "someBoolean"
      description = "Some boolean value"
      type        = "boolean"
      value = {
        boolean = {
          value = true
        }
      }
    }
  ]

  force_delete = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("aria_orchestrator_configuration.test", "id"),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "name",
						"Test Config Renamed",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "description",
						"Config generated by the acceptance tests of Aria provider (updated).",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "version",
						"0.1.0",
					),
					resource.TestMatchResourceAttr(
						"aria_orchestrator_configuration.test", "version_id",
						regexp.MustCompile("[0-9a-f]{40}"),
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "attributes.#",
						"2",
					),
					resource.TestCheckResourceAttr(
						"aria_orchestrator_configuration.test", "force_delete",
						"true",
					),
					resource.TestCheckResourceAttrPair(
						"aria_orchestrator_configuration.test", "category_id",
						"aria_orchestrator_category.misc", "id",
					),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aria_orchestrator_configuration.test",
				ImportState:       true,
				ImportStateVerify: true,

				// Prevent diff on force_delete field
				ImportStateVerifyIgnore: []string{"force_delete"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

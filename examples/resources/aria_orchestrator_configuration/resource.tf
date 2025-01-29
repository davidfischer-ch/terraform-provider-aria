# main.tf

resource "aria_orchestrator_category" "my_company" {
  name      = "MyCompany"
  type      = "ConfigurationElementCategory"
  parent_id = ""
}

resource "aria_orchestrator_configuration" "Dummy" {
  name        = "Dummy Config"
  description = "Example configuration showing all (manageable) attribute types."
  category_id = aria_orchestrator_category.my_company.id
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

  # If required...
  force_delete = true
}

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
    /*
    This is not yet well handled (mutated by the platform), we have to find a pattern for this.
    We cannot (?) expose 1-1 the API to make it declarative, we have to tackle this challenge.
    {
      name        = "someSecureString"
      description = "Some secure string value"
      type        = "SecureString"
      value = {
        secure_string = {
          value         = "test" -> "A.....Z", changing everytime ...
          is_plain_text = true -> false
        }
      }
    }*/
  ]

  # If required...
  force_delete = true
}

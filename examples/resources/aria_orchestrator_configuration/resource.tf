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

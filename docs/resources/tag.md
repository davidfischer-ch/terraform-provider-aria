---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "aria_tag Resource - aria"
subcategory: ""
description: |-
  Tag resource
---

# aria_tag (Resource)

Tag resource

## Example Usage

```terraform
resource "aria_tag" "basic_example" {
  key             = "MyTag"
  value           = "My Value"
  force_delete    = true  # Force deletion even if in use (use with caution)
  keep_on_destroy = false # Let Terraform delete icon for real on destroy (the default)
}

resource "aria_tag" "shared_example" {
  key             = "AnotherTag"
  keep_on_destroy = true # Do not delete, preventing issues if same tag is declared multiple time
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `key` (String) Key (force recreation on change)

### Optional

- `force_delete` (Boolean) Force destroying the tag (bypass references check).
- `keep_on_destroy` (Boolean) Keep the tag on destroy?
This can help preventing issues if this tag should never be destroyed for good reasons.
Default value is false.
- `value` (String) Value (force recreation on change)

### Read-Only

- `id` (String) Identifier

## Import

Import is supported using the following syntax:

```shell
# Tag can be imported by specifying the instance's unique identifier.
terraform import aria_tag.example 9ea6205b-e0e1-4188-b275-b17299efe49a
```

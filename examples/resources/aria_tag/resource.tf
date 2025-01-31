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

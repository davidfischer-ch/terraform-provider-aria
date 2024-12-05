resource "aria_icon" "basic_example" {
  path            = "icon.svg"
  hash            = filesha256("icon.svg") # Allow tracking content change
  keep_on_destroy = false                  # Let Terraform delete icon for real on destroy
}

resource "aria_icon" "shared_example" {
  path            = "common.png"
  hash            = filesha256("common.png") # Allow tracking content change
  keep_on_destroy = true                     # Do not delete, preventing issues if same icon is declared multiple time
}

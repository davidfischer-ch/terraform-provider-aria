resource "aria_icon" "example" {
  path = "icon.svg"
  hash = filesha256("icon.svg") # Allow tracking content change
}

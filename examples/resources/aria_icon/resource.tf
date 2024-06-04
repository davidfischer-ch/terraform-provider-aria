resource "aria_icon" "example" {
  content = file("icon.svg")
}

output "example_icon" {
  value = {
    id   = aria_icon.example.id
    hash = sha256(aria_icon.example.content)
  }
}

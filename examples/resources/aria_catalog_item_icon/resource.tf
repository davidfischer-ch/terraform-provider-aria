resource "aria_icon" "example" {
  path = "icon.svg"
  hash = filesha256("icon.svg") # Allow tracking content change
}

resource "aria_catalog_item_icon" "example" {
  item_id = "746f8902-c2d2-4f91-9a8d-1d4ff27033fd"
  icon_id = aria_icon.example.id
}

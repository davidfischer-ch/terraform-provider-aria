# Icon read by the icon data source tests (TF_VAR_test_icon_id).

resource "aria_icon" "test" {
  path            = "${path.module}/icon.svg"
  hash            = filesha256("${path.module}/icon.svg")
  keep_on_destroy = true
}

# Icon read by the icon data source tests (TF_VAR_test_icon_id).
#
# Aria derives the icon identifier from its content, this file must stay different from
# tests/icon.svg and tests/icon.png, otherwise the icon tests would destroy this fixture.

resource "aria_icon" "test" {
  path            = "${path.module}/icon.svg"
  hash            = filesha256("${path.module}/icon.svg")
  keep_on_destroy = false
}

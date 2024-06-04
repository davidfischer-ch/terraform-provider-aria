resource "aria_abx_constant" "example" {
  name  = "THIS_IS_MY_CONSTANT"
  value = "42"
}

output "example_constant" {
  value = "My constant ${aria_abx_constant.example.name} ID is ${aria_abx_constant.example.id}"
}

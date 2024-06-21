resource "aria_abx_sensitive_constant" "example" {
  name  = "THIS_IS_MY_SENSITIVE_CONSTANT"
  value = "1234pass"
}

output "example_sensitive_constant" {
  value = "My sensitive constant ${aria_abx_sensitive_constant.example.name} ID is ${aria_abx_sensitive_constant.example.id}"
}

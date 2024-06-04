resource "aria_abx_secret" "example" {
  name  = "THIS_IS_MY_SECRET"
  value = "1234pass"
}

output "example_secret" {
  value = "My secret ${aria_abx_secret.example.name} ID is ${aria_abx_secret.example.id}"
}

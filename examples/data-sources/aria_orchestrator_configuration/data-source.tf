data "aria_orchestrator_configuration" "some" {
  id = "eb261901-4dc3-aa22-4957-19e7a5b6b4d"
}

output "configuration" {
  value     = data.aria_orchestrator_configuration.some
  sensitive = true # Terraform is picky (for a good reason)
}

# Changes to Outputs:
#  + configuration = (sensitive value)
# ...sorry

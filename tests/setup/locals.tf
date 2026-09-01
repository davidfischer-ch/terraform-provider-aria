locals {
  # The cleanup binary (cmd/cleanup) deletes everything named ARIA_PROVIDER_TEST*, fixtures are
  # named differently to survive it.
  prefix = "ARIA_PROVIDER_FIXTURE"

  is_vcf9 = var.aria_tenant != null && var.aria_tenant != ""

  # Approvers are principals stripped of their domain, e.g. USER:ELIOTT.
  approver_name = "USER:${split("@", coalesce(
    try(local.me.acct, null),
    try(local.me.username, null),
    try(local.me.email, null),
  ))[0]}"

  me = data.restful_resource.me.output

  test_runtime = coalesce(var.test_runtime, local.is_vcf9 ? "python:3.11" : "python:3.10")
}

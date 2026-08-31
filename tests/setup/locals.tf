locals {
  # The cleanup binary (cmd/cleanup) deletes everything named ARIA_PROVIDER_TEST*, fixtures are
  # named differently to survive it.
  prefix = "ARIA_PROVIDER_FIXTURE"

  me = data.restful_resource.me.output

  # Approvers are principals stripped of their domain, e.g. USER:ELIOTT.
  approver_name = "USER:${split("@", coalesce(
    try(local.me.acct, null),
    try(local.me.username, null),
    try(local.me.email, null),
  ))[0]}"
}

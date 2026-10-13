# Secret read by the secret data source tests (TF_VAR_test_secret_id).
#
# Organization scoped, hence independent of the fixture projects: nothing but this secret refers to
# it, and a leftover can not block the deletion of a project.
#
# The value is a placeholder, no test reads it back, only the metadata of the secret is asserted.
# Needs an aria provider shipping aria_secret, run with DEV=1 until it is released.

resource "aria_secret" "secret" {
  name        = "${local.prefix}_SECRET"
  description = "Secret read by Aria provider's acceptance tests."
  org_scoped  = true
  value       = "${local.prefix}_SECRET_PLACEHOLDER_VALUE"
}

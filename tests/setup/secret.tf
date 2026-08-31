# Secret read by the secret data source tests (TF_VAR_test_secret_id).
#
# Neither the aria nor the vra provider exposes a secret resource, the platform API is called
# directly. Replacing this by an aria_secret resource would drop the restful provider and the three
# variables it needs.
#
# The value is a placeholder, no test reads it back, only the metadata of the secret is asserted.

resource "restful_resource" "secret" {
  path         = "/platform/api/secrets"
  read_path    = "$(path)/$(body.id)"
  output_attrs = ["id"]

  body = {
    name        = "${local.prefix}_SECRET"
    description = "Secret read by Aria provider's acceptance tests."
    projectId   = vra_project.test[0].id
  }

  # Merge-patched into body on write. The API never returns the value back, keeping it here avoids
  # a permanent drift.
  ephemeral_body = {
    value = "${local.prefix}_SECRET_PLACEHOLDER_VALUE"
  }
}

# Aria Automation 8.x exchange of the refresh token for an access token, feeding the restful
# provider below. VCF 9 needs none of this, its exchange is a plain OAuth2 refresh_token grant the
# provider performs on its own, hence the count.

data "restful_resource" "login" {
  count = local.is_vcf9 ? 0 : 1

  provider = restful.anonymous
  id       = "/iaas/api/login"
  method   = "POST"

  body = {
    refreshToken = var.aria_refresh_token
  }

  use_sensitive_output = true
}

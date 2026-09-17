# The aria and vra providers exchange the refresh token themselves. Only restful needs an access
# token handed to it, and it gets one without asking: through the OAuth2 grant below on VCF 9,
# through the login data source on Aria Automation 8.x. Nothing here needs ARIA_ACCESS_TOKEN.
#
# The aria provider is left implicit, it reads ARIA_* on its own.

# Unauthenticated, used solely by the Aria Automation 8.x login data source.
provider "restful" {
  alias    = "anonymous"
  base_url = var.aria_host

  client = {
    tls_insecure_skip_verify = var.aria_insecure
  }
}

provider "restful" {
  base_url = var.aria_host

  client = {
    tls_insecure_skip_verify = var.aria_insecure
  }

  security = {
    # Aria Automation 8.x
    http = local.is_vcf9 ? null : {
      token = {
        token = try(one(data.restful_resource.login).sensitive_output.token, null)
      }
    }
    # VCF 9
    oauth2 = local.is_vcf9 ? {
      refresh_token = {
        token_url     = "${var.aria_host}/tm/oauth/tenant/${var.aria_tenant}/token"
        refresh_token = var.aria_refresh_token
        token_type    = "Bearer"
      }
    } : null
  }
}

provider "vra" {
  url           = var.aria_host
  organization  = local.is_vcf9 ? var.aria_tenant : null
  refresh_token = var.aria_refresh_token
  insecure      = var.aria_insecure
}

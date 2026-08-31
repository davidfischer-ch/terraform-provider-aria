# Only the restful provider is configured here, aria and vra read their own environment variables.

provider "restful" {
  base_url = var.aria_host

  client = {
    tls_insecure_skip_verify = var.aria_insecure
  }

  security = {
    http = {
      token = {
        token = var.aria_access_token
      }
    }
  }
}

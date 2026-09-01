# The aria provider reads ARIA_* itself, every attribute of its schema being optional. The vra and
# restful providers need an explicit configuration, which these variables feed.

variable "aria_host" {
  type        = string
  description = "The URI to Aria, mirror of ARIA_HOST."

  validation {
    condition     = startswith(var.aria_host, "https://")
    error_message = "Argument `aria_host` must be an URI, is ARIA_HOST exported?"
  }
}

variable "aria_tenant" {
  type        = string
  description = "The VCF 9 organization (tenant) name, null on Aria Automation 8.x, mirror of ARIA_TENANT."
  default     = null
}

variable "aria_refresh_token" {
  type        = string
  description = "The refresh token to use for making API requests, mirror of ARIA_REFRESH_TOKEN."
  sensitive   = true

  validation {
    condition     = length(var.aria_refresh_token) > 0
    error_message = "Argument `aria_refresh_token` must not be empty, is ARIA_REFRESH_TOKEN exported?"
  }
}

variable "aria_insecure" {
  type        = bool
  description = "Whether server should be accessed without verifying the TLS certificate, mirror of ARIA_INSECURE."
  default     = false
}

variable "test_runtime" {
  type        = string
  description = <<-EOT
    The Orchestrator runtime the environment, action and repository tests are built on.
    Defaults to python:3.11 on VCF 9 and python:3.10 on Aria Automation 8.x.
  EOT
  default     = null

  validation {
    condition     = var.test_runtime == null ? true : length(var.test_runtime) > 0
    error_message = "Argument `test_runtime` must not be empty."
  }
}

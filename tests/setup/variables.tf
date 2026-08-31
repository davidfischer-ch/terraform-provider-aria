# The aria and vra providers read ARIA_* and VRA_* themselves, these three are only consumed by the
# restful provider, which has no environment variable defaults.

variable "aria_host" {
  type        = string
  description = "The URI to Aria, mirror of ARIA_HOST."

  validation {
    condition     = startswith(var.aria_host, "https://")
    error_message = "Argument `aria_host` must be an URI, is ARIA_HOST exported?"
  }
}

variable "aria_access_token" {
  type        = string
  description = "The access token to use for making API requests, mirror of ARIA_ACCESS_TOKEN."
  sensitive   = true

  validation {
    condition     = length(var.aria_access_token) > 0
    error_message = "Argument `aria_access_token` must not be empty, is ARIA_ACCESS_TOKEN exported?"
  }
}

variable "aria_insecure" {
  type        = bool
  description = "Whether server should be accessed without verifying the TLS certificate, mirror of ARIA_INSECURE."
  default     = false
}

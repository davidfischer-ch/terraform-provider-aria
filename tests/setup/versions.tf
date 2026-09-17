terraform {
  required_version = ">= 1.11"

  required_providers {
    aria = {
      source  = "davidfischer-ch/aria"
      version = ">= 0.7.4"
    }
    local = {
      source  = "hashicorp/local"
      version = ">= 2.9.0"
    }
    restful = {
      source  = "magodo/restful"
      version = ">= 0.25.2"
    }
    time = {
      source  = "hashicorp/time"
      version = ">= 0.14.1"
    }
    vra = {
      source  = "vmware/vra"
      version = ">= 0.17.3"
    }
  }
}

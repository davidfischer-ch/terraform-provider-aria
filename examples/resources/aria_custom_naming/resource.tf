# variables.tf

variable "org_id" {
  type = string
}

# locals.tf

locals {
  pattern = join("-", [
    "$${VM_Common.os_family}",
    "$${VM_Common.vlan}",
    "$${VM_Common.srv_role}",
    "$${#####}" # Counter
  ])
}

resource "aria_custom_naming" "machines" {
  name        = "MachineNaming"
  description = "Standardized naming convention for our machines."

  projects = [
    {
      active       = true
      org_default  = true
      org_id       = var.org_id
      project_id   = "*"
      project_name = "*"
    }
  ]

  templates = {
    "COMPUTE.Machine > Default" = {
      name               = ""
      resource_type      = "COMPUTE"
      resource_type_name = "Machine"
      unique_name        = true
      pattern            = local.pattern
      static_pattern     = ""
      start_counter      = 1
      incrment_step      = 1
    }
    "COMPUTE.Machine > RH-LAB-CACHE" = {
      name               = ""
      resource_type      = "COMPUTE"
      resource_type_name = "Machine"
      unique_name        = true
      pattern            = local.pattern
      static_pattern     = "RH-LAB-CACHE"
      start_counter      = 1
      incrment_step      = 1
    }
    "COMPUTE.Machine > WS-LAB-CACHE" = {
      name               = ""
      resource_type      = "COMPUTE"
      resource_type_name = "Machine"
      unique_name        = true
      pattern            = local.pattern
      static_pattern     = "WS-LAB-CACHE"
      start_counter      = 1
      incrment_step      = 1
    }
  }
}

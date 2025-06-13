# main.tf

resource "aria_policy" "snapshot_revert_approval" {
  name             = "Snapshot Revert Approval Policy"
  description      = "Require some user to approve a snapshot's revert operation."
  enforcement_type = "HARD"
  type_id          = "com.vmware.policy.approval"

  scope_criteria = jsonencode({
    matchExpression = [
      {
        or = [
          {
            key      = "project.name"
            operator = "contains"
            value    = "REC"
          },
          {
            key      = "project.name"
            operator = "contains"
            value    = "PROD"
          }
        ]
      }
    ]
  })

  definition = jsonencode({
    level                = 1
    actions              = ["Cloud.vSphere.Machine.Snapshot.Revert"]
    approvers            = ["USER:SOMEUSER"]
    approvalMode         = "ANY_OF"
    approverType         = "USER"
    autoApprovalExpiry   = 7
    autoApprovalDecision = "REJECT"
  })
}

locals {
  redis_custom_resource = {
    type = "Redis_v1.0"
  }

  redis_template = {
    catalog_item = {
      id = "c7e95ef0-7608-4b09-b11b-87a80b9347aa"
    }
  }
}

resource "aria_policy" "redis_day2" {
  name             = "${local.redis_custom_resource.type} - Day 2 Policy"
  description      = "Restrict access to Redis Day 2 actions to a restricted set of users"
  enforcement_type = "HARD"
  type_id          = "com.vmware.policy.deployment.action"

  definition = jsonencode({
    allowedActions = [
      {
        actions = [
          "Deployment.Update",
          "Deployment.Delete",
          "${local.redis_custom_resource.type}.custom.snapshot",
          "${local.redis_custom_resource.type}.custom.restore",
          "${local.redis_custom_resource.type}.custom.getresourcestatus",
        ]
        authorities = [
          "GROUP:<redacted>.ADMINISTRATOR@<redacted>@<redacted>",
          "GROUP:<redacted>.DEV@<redacted>@<redacted>",
          "GROUP:<redacted>.OPS@<redacted>@<redacted>"
        ]
      }
    ]
  })

  scope_criteria = jsonencode({
    matchExpression = [
      {
        or = [
          for project_name in ["LAB", "DEV", "REC", "PROD"] :
          {
            key      = "project.name"
            operator = "eq"
            value    = project_name
          }
        ]
      }
    ]
  })

  criteria = jsonencode({
    matchExpression = [
      {
        key      = "catalogItemId"
        operator = "eq"
        value    = local.redis_template.catalog_item.id
      }
    ]
  })
}

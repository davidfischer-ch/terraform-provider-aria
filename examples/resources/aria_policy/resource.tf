# main.tf

resource "aria_policy" "test" {
  name             = "Snapshot Revert Approval Policy"
  description      = "Require some user to approve a snapshot's revert operation."
  enforcement_type = "HARD"
  type_id          = "com.vmware.policy.approval"

  scope_criteria = jsonencode({
    matchExpression = [
      {
        or = [
          {
            or = [
              {
                key      = "project.name"
                operator = "contains"
                value    = "PROD"
              }
            ]
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

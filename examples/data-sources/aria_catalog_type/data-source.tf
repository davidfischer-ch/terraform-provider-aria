data "aria_catalog_type" "abx_actions" {
  id = "com.vmw.abx.actions"
}

output "abx_actions_catalog_type" {
  value = data.aria_catalog_type.abx_actions
}

# Changes to Outputs:
#   + abx_actions_catalog_type = {
#       + base_uri   = "http://abx-service.prelude.svc.cluster.local/abx/api/catalog"
#       + created_at = "2023-05-05T00:37:14.955948Z"
#       + created_by = "abx-KfXBHw4DSir1X516"
#       + icon_id    = "c2c83e12-a908-30a6-b138-2b134b3b29bb"
#       + id         = "com.vmw.abx.actions"
#       + name       = "Extensibility actions"
#    }

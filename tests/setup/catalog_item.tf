# Catalog item whose icon and custom form are modified by the acceptance tests
# (TF_VAR_test_catalog_item_id and TF_VAR_test_catalog_item_type).
#
# A catalog item is a released cloud template published through a catalog source.

resource "vra_blueprint" "test" {
  name        = "${local.prefix}_CLOUD_TEMPLATE"
  description = "Cloud template published to give Aria provider's acceptance tests a catalog item."
  project_id  = vra_project.test[0].id

  content = <<-EOT
    formatVersion: 1
    inputs:
      message:
        type: string
        title: Message
        default: Hello World!
    resources:
      Machine:
        type: Cloud.Machine
        properties:
          image: ubuntu
          flavor: small
  EOT
}

resource "vra_blueprint_version" "test" {
  blueprint_id = vra_blueprint.test.id
  description  = "Version published to the catalog."
  change_log   = "Initial version."
  version      = "1"
  release      = true
}

resource "vra_catalog_source_blueprint" "test" {
  name        = "${local.prefix}_CATALOG_SOURCE"
  description = "Catalog source publishing the cloud template of the fixture project."
  project_id  = vra_project.test[0].id

  depends_on = [vra_blueprint_version.test]
}

# The import of a catalog source is asynchronous, the item isn't queryable right away.
resource "time_sleep" "catalog_import" {
  create_duration = "30s"
  triggers = {
    catalog_source_id = vra_catalog_source_blueprint.test.id
  }
}

# The externalId of a catalog item is not the identifier of the cloud template it publishes, the
# aria data source cannot find it from the blueprint. Select the only item of the catalog source
# instead, then read it by identifier.
data "restful_resource" "catalog_item" {
  id       = "/catalog/api/admin/items"
  selector = "content.0"

  query = {
    sourceIds = [vra_catalog_source_blueprint.test.id]
  }

  depends_on = [time_sleep.catalog_import]
}

data "aria_catalog_item" "test" {
  id = data.restful_resource.catalog_item.output.id
}

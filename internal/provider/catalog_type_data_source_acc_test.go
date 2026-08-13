// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCatalogTypeDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: `
data "aria_catalog_type" "abx_actions" {
  id = "com.vmw.abx.actions"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.aria_catalog_type.abx_actions", "id", "com.vmw.abx.actions"),
					resource.TestCheckResourceAttr("data.aria_catalog_type.abx_actions", "name", "Extensibility actions"),
					resource.TestCheckResourceAttr("data.aria_catalog_type.abx_actions", "base_uri", "http://abx-service.prelude.svc.cluster.local/abx/api/catalog"),
					resource.TestCheckResourceAttrSet("data.aria_catalog_type.abx_actions", "created_at"),
					resource.TestCheckResourceAttrSet("data.aria_catalog_type.abx_actions", "created_by"),
					resource.TestCheckResourceAttrSet("data.aria_catalog_type.abx_actions", "icon_id"),
				),
			},
		},
	})
}

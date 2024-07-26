// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func CatalogTypeDataSourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Catalog type data source",
		Attributes: map[string]schema.Attribute{
			"id": RequiredIdentifierSchema(),
			"name": schema.StringAttribute{
				MarkdownDescription: "Type name",
				Computed:            true,
			},
			"base_uri": schema.StringAttribute{
				MarkdownDescription: "Base URI",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Creation date",
				Computed:            true,
			},
			"created_by": schema.StringAttribute{
				MarkdownDescription: "Ask VMware",
				Computed:            true,
			},
			"icon_id": schema.StringAttribute{
				MarkdownDescription: "Icon identifier",
				Computed:            true,
			},
		},
	}
}

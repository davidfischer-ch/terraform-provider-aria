// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func CatalogItemDataSourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Catalog Item's data source",
		Attributes: map[string]schema.Attribute{
			"id": OptionalIdentifierSchema("Identifier"),
			"name": schema.StringAttribute{
				MarkdownDescription: "Name",
				Computed:            true,
				Optional:            true,
			},
			"description": ComputedDescriptionSchema(),
			"schema": schema.StringAttribute{
				MarkdownDescription: "Schema" + JSON_INSTEAD_OF_DYNAMIC_DISCLAIMER,
				CustomType:          jsontypes.NormalizedType{},
				Optional:            true,
			},
			"external_id": OptionalIdentifierSchema("External identifier"),
			"icon_id":     ComputedIdentifierSchema("Icon identifier"),
			"form_id":     ComputedIdentifierSchema("Form identifier"),
			"type_id":     OptionalIdentifierSchema("Catalog type identifier"),
			"source_id":   OptionalIdentifierSchema("Catalog source identifier"),
			"source_name": schema.StringAttribute{
				MarkdownDescription: "Catalog source name",
				Computed:            true,
				Optional:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Creation timestamp (RFC3339)",
				CustomType:          timetypes.RFC3339Type{},
				Computed:            true,
			},
			"created_by": schema.StringAttribute{
				MarkdownDescription: "User who created the resource",
				Computed:            true,
			},
			"last_updated_at": schema.StringAttribute{
				MarkdownDescription: "Last update timestamp (RFC3339)",
				CustomType:          timetypes.RFC3339Type{},
				Computed:            true,
			},
			"last_updated_by": schema.StringAttribute{
				MarkdownDescription: "Last user who updated the resource",
				Computed:            true,
			},
		},
	}
}

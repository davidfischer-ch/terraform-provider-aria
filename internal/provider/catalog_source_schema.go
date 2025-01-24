// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func CatalogSourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Catalog source resource",
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "Source name (e.g. getVRAHost)",
				Required:            true,
			},
			"type_id": schema.StringAttribute{
				MarkdownDescription: "Source type (e.g. com.vmw.vro.workflow)",
				Required:            true,
			},
			"global": schema.BoolAttribute{
				MarkdownDescription: "Is it globally shared?",
				Computed:            true,
			},
			"config": CatalogSourceConfigSchema(),
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
			"last_import_started_at": schema.StringAttribute{
				MarkdownDescription: "Last import start timestamp (RFC3339)",
				CustomType:          timetypes.RFC3339Type{},
				Computed:            true,
			},
			"last_import_completed_at": schema.StringAttribute{
				MarkdownDescription: "Last import end timestamp (RFC3339)",
				CustomType:          timetypes.RFC3339Type{},
				Computed:            true,
			},
			/*"last_import_errors": schema.ListNestedAttribute{
				MarkdownDescription: "Action input parameters",
				Required:            true,
				NestedObject:        ParameterSchema(),
			},*/
			"items_imported": schema.Int32Attribute{
				MarkdownDescription: "Number of imported items",
				Computed:            true,
			},
			"items_found": schema.Int32Attribute{
				MarkdownDescription: "Number of existing items",
				Computed:            true,
			},
		},
	}
}

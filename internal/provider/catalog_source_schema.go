// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func CatalogSourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Catalog source resource. Please be aware of this: " +
			"https://github.com/davidfischer-ch/terraform-provider-aria/issues/114",
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "Source name (e.g. getVRAHost)",
				Required:            true,
			},
			"description": RequiredDescriptionSchema(),
			"type_id": schema.StringAttribute{
				MarkdownDescription: "Source type (e.g. `com.vmw.vro.workflow`)",
				Required:            true,
			},
			"global": schema.BoolAttribute{
				MarkdownDescription: "Is it globally shared?",
				Computed:            true,
			},
			"project_id": OptionalImmutableProjectIdSchema(),
			"config":     CatalogSourceConfigSchema(),
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Creation timestamp (RFC3339)",
				CustomType:          timetypes.RFC3339Type{},
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_by": schema.StringAttribute{
				MarkdownDescription: "User who created the resource",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
			"last_import_errors": schema.ListAttribute{
				MarkdownDescription: "Action input parameters",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"items_found": schema.Int32Attribute{
				MarkdownDescription: "Number of existing items",
				Computed:            true,
			},
			"items_imported": schema.Int32Attribute{
				MarkdownDescription: "Number of imported items",
				Computed:            true,
			},
			"wait_imported": schema.BoolAttribute{
				MarkdownDescription: "Wait for import to be completed " +
					"(up to 15 minutes, checked every 30 seconds)",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
		},
	}
}

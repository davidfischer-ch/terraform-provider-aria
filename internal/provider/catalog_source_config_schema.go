// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
)

// The Configuration embedded inside a CatalogSourceSchema.
func CatalogSourceConfigSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "Configuration",
		Required:            true,
		Attributes: map[string]schema.Attribute{
			"source_project_id": schema.StringAttribute{
				MarkdownDescription: "Project to make available " +
					"(required for Cloud Templates or ABX Actions catalog sources)",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"workflows": schema.ListNestedAttribute{
				MarkdownDescription: "Workflows to make available " +
					"(required for Orchestrator Worflows catalog sources)",
				Optional:     true,
				NestedObject: CatalogSourceWorkflowSchema(),
			},
		},
	}
}

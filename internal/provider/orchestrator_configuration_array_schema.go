// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// A String embedded inside an Orchestrator Configuration Value.
func OrchestratorConfigurationArraySchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "Array",
		Optional:            true,
		Attributes: map[string]schema.Attribute{
			"elements": schema.ListNestedAttribute{
				MarkdownDescription: "Elements",
				Required:            true,
				NestedObject:        OrchestratorConfigurationArrayElementSchema(),
			},
		},
	}
}

// A String embedded inside a Computed Orchestrator Configuration Value.
func ComputedOrchestratorConfigurationArraySchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "Array",
		Computed:            true,
		Attributes: map[string]schema.Attribute{
			"elements": schema.ListNestedAttribute{
				MarkdownDescription: "Elements",
				Computed:            true,
				NestedObject:        ComputedOrchestratorConfigurationArrayElementSchema(),
			},
		},
	}
}

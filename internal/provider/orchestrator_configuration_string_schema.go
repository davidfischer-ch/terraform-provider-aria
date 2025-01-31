// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// A String embedded inside an Orchestrator Configuration Value.
func OrchestratorConfigurationStringSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "String",
		Optional:            true,
		Attributes: map[string]schema.Attribute{
			"value": schema.StringAttribute{
				MarkdownDescription: "Value",
				Required:            true,
			},
		},
	}
}

// A String embedded inside a Computed Orchestrator Configuration Value.
func ComputedOrchestratorConfigurationStringSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "String",
		Computed:            true,
		Attributes: map[string]schema.Attribute{
			"value": schema.StringAttribute{
				MarkdownDescription: "Value",
				Computed:            true,
			},
		},
	}
}

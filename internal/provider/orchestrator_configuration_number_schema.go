// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// A Number embedded inside an Orchestrator Configuration Value.
func OrchestratorConfigurationNumberSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "Number",
		Optional:            true,
		Attributes: map[string]schema.Attribute{
			"value": schema.Float64Attribute{
				MarkdownDescription: "Value",
				Required:            true,
			},
		},
	}
}

// A Number embedded inside a Computed Orchestrator Configuration Value.
func ComputedOrchestratorConfigurationNumberSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "Number",
		Computed:            true,
		Attributes: map[string]schema.Attribute{
			"value": schema.Float64Attribute{
				MarkdownDescription: "Value",
				Computed:            true,
			},
		},
	}
}

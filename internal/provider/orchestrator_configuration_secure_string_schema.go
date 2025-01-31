// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// A Secure String embedded inside an Orchestrator Configuration Value.
func OrchestratorConfigurationSecureStringSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "Secure String",
		Optional:            true,
		Attributes: map[string]schema.Attribute{
			"value": schema.StringAttribute{
				MarkdownDescription: "Value",
				Required:            true,
				Sensitive:           true,
			},
			"is_plain_text": schema.BoolAttribute{
				MarkdownDescription: "Plain text?",
				Required:            true,
			},
		},
	}
}

// A Secure String embedded inside a Computed Orchestrator Configuration Value.
func ComputedOrchestratorConfigurationSecureStringSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "Secure String",
		Computed:            true,
		Attributes: map[string]schema.Attribute{
			"value": schema.StringAttribute{
				MarkdownDescription: "Value",
				Computed:            true,
				Sensitive:           true,
			},
			"is_plain_text": schema.BoolAttribute{
				MarkdownDescription: "Plain text?",
				Computed:            true,
			},
		},
	}
}

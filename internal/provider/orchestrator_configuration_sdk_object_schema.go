// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// A SDK Object embedded inside an Orchestrator Configuration Value.
func OrchestratorConfigurationSDKObjectSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "SDK Object",
		Optional:            true,
		Attributes: map[string]schema.Attribute{
			"id": RequiredIdentifierSchema(""),
			"type": schema.StringAttribute{
				MarkdownDescription: "Type",
				Required:            true,
			},
		},
	}
}

// A SDK Object embedded inside a Computed Orchestrator Configuration Value.
func ComputedOrchestratorConfigurationSDKObjectSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "SDK Object",
		Computed:            true,
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"type": schema.StringAttribute{
				MarkdownDescription: "Type",
				Computed:            true,
			},
		},
	}
}

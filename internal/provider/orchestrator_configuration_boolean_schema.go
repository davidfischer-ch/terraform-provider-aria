// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// A Boolean embedded inside an Orchestrator Configuration Value.
func OrchestratorConfigurationBooleanSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "Boolean",
		Optional:            true,
		Attributes: map[string]schema.Attribute{
			"value": schema.BoolAttribute{
				MarkdownDescription: "Value",
				Required:            true,
			},
		},
	}
}

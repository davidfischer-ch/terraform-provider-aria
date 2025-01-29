// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// The Value embedded inside an Orchestrator Configuration Attribute.
func OrchestratorConfigurationValueSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "Value",
		Required:            true,
		Attributes: map[string]schema.Attribute{
			"boolean": OrchestratorConfigurationBooleanSchema(),
			"string":  OrchestratorConfigurationStringSchema(),
		},
	}
}

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
			"boolean":       OrchestratorConfigurationBooleanSchema(),
			"string":        OrchestratorConfigurationStringSchema(),
			"secure_string": OrchestratorConfigurationSecureStringSchema(),
			"sdk_object":    OrchestratorConfigurationSDKObjectSchema(),
		},
	}
}

// The Value embedded inside a Computed Orchestrator Configuration Attribute.
func ComputedOrchestratorConfigurationValueSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "Value",
		Computed:            true,
		Attributes: map[string]schema.Attribute{
			"boolean":       ComputedOrchestratorConfigurationBooleanSchema(),
			"string":        ComputedOrchestratorConfigurationStringSchema(),
			"secure_string": ComputedOrchestratorConfigurationSecureStringSchema(),
			"sdk_object":    ComputedOrchestratorConfigurationSDKObjectSchema(),
		},
	}
}

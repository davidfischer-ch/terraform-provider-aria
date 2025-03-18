// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// The Values embedded inside an Orchestrator Configuration Attribute.
func OrchestratorConfigurationArrayElementSchema() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"boolean":       OrchestratorConfigurationBooleanSchema(),
			"string":        OrchestratorConfigurationStringSchema(),
			"secure_string": OrchestratorConfigurationSecureStringSchema(),
			"sdk_object":    OrchestratorConfigurationSDKObjectSchema(),
		},
	}
}

// The Values embedded inside a Computed Orchestrator Configuration Attribute.
func ComputedOrchestratorConfigurationArrayElementSchema() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"boolean":       ComputedOrchestratorConfigurationBooleanSchema(),
			"string":        ComputedOrchestratorConfigurationStringSchema(),
			"secure_string": ComputedOrchestratorConfigurationSecureStringSchema(),
			"sdk_object":    ComputedOrchestratorConfigurationSDKObjectSchema(),
		},
	}
}

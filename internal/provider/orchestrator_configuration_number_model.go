// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// OrchestratorConfigurationNumberModel describes the resource data model.
type OrchestratorConfigurationNumberModel struct {
	Value types.Float64 `tfsdk:"value"`
}

// OrchestratorConfigurationNumberAPIModel describes the resource API model.
type OrchestratorConfigurationNumberAPIModel struct {
	Value float64 `json:"value"`
}

func (self *OrchestratorConfigurationNumberModel) FromAPI(
	raw OrchestratorConfigurationNumberAPIModel,
) {
	self.Value = types.Float64Value(raw.Value)
}

func (self OrchestratorConfigurationNumberModel) ToAPI() OrchestratorConfigurationNumberAPIModel {
	return OrchestratorConfigurationNumberAPIModel{
		Value: self.Value.ValueFloat64(),
	}
}

// Utils -------------------------------------------------------------------------------------------

// Used to convert structure to a types.Object.
func (self OrchestratorConfigurationNumberModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"value": types.Float64Type,
	}
}

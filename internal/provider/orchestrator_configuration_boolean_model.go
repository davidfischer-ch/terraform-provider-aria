// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// OrchestratorConfigurationBooleanModel describes the resource data model.
type OrchestratorConfigurationBooleanModel struct {
	Value types.Bool `tfsdk:"value"`
}

// OrchestratorConfigurationBooleanAPIModel describes the resource API model.
type OrchestratorConfigurationBooleanAPIModel struct {
	Value bool `json:"value"`
}

func (self *OrchestratorConfigurationBooleanModel) FromAPI(
	raw OrchestratorConfigurationBooleanAPIModel,
) {
	self.Value = types.BoolValue(raw.Value)
}

func (self OrchestratorConfigurationBooleanModel) ToAPI() OrchestratorConfigurationBooleanAPIModel {
	return OrchestratorConfigurationBooleanAPIModel{
		Value: self.Value.ValueBool(),
	}
}

// Utils -------------------------------------------------------------------------------------------

// Used to convert structure to a types.Object.
func (self OrchestratorConfigurationBooleanModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"value": types.BoolType,
	}
}

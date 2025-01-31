// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// OrchestratorConfigurationStringModel describes the resource data model.
type OrchestratorConfigurationStringModel struct {
	Value types.String `tfsdk:"value"`
}

// OrchestratorConfigurationStringAPIModel describes the resource API model.
type OrchestratorConfigurationStringAPIModel struct {
	Value string `json:"value"`
}

func (self *OrchestratorConfigurationStringModel) FromAPI(
	raw OrchestratorConfigurationStringAPIModel,
) {
	self.Value = types.StringValue(raw.Value)
}

func (self OrchestratorConfigurationStringModel) ToAPI() OrchestratorConfigurationStringAPIModel {
	return OrchestratorConfigurationStringAPIModel{
		Value: self.Value.ValueString(),
	}
}

// Utils -------------------------------------------------------------------------------------------

// Used to convert structure to a types.Object.
func (self OrchestratorConfigurationStringModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"value": types.StringType,
	}
}

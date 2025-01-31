// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// OrchestratorConfigurationSecureStringModel describes the resource data model.
type OrchestratorConfigurationSecureStringModel struct {
	Value       types.String `tfsdk:"value"`
	IsPlainText types.Bool   `tfsdk:"is_plain_text"`
}

// OrchestratorConfigurationSecureStringAPIModel describes the resource API model.
type OrchestratorConfigurationSecureStringAPIModel struct {
	Value       string `json:"value"`
	IsPlainText bool   `json:"isPlainText"`
}

func (self *OrchestratorConfigurationSecureStringModel) FromAPI(
	raw OrchestratorConfigurationSecureStringAPIModel,
) {
	self.Value = types.StringValue(raw.Value)
	self.IsPlainText = types.BoolValue(raw.IsPlainText)
}

func (self OrchestratorConfigurationSecureStringModel) ToAPI() OrchestratorConfigurationSecureStringAPIModel {
	return OrchestratorConfigurationSecureStringAPIModel{
		Value:       self.Value.ValueString(),
		IsPlainText: self.IsPlainText.ValueBool(),
	}
}

// Utils -------------------------------------------------------------------------------------------

// Used to convert structure to a types.Object.
func (self OrchestratorConfigurationSecureStringModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"value":         types.StringType,
		"is_plain_text": types.BoolType,
	}
}

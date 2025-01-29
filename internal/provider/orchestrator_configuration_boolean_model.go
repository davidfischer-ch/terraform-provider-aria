// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	ctx context.Context,
	raw OrchestratorConfigurationBooleanAPIModel,
) diag.Diagnostics {
	self.Value = types.BoolValue(raw.Value)
	return diag.Diagnostics{}
}

func (self OrchestratorConfigurationBooleanModel) ToAPI(
	ctx context.Context,
) (OrchestratorConfigurationBooleanAPIModel, diag.Diagnostics) {
	return OrchestratorConfigurationBooleanAPIModel{
		Value: self.Value.ValueBool(),
	}, diag.Diagnostics{}
}

// Utils -------------------------------------------------------------------------------------------

// Used to convert structure to a types.Object.
func (self OrchestratorConfigurationBooleanModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"value": types.BoolType,
	}
}

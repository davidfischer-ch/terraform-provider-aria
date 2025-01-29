// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	ctx context.Context,
	raw OrchestratorConfigurationStringAPIModel,
) diag.Diagnostics {
	self.Value = types.StringValue(raw.Value)
	return diag.Diagnostics{}
}

func (self OrchestratorConfigurationStringModel) ToAPI(
	ctx context.Context,
) (OrchestratorConfigurationStringAPIModel, diag.Diagnostics) {
	return OrchestratorConfigurationStringAPIModel{
		Value: self.Value.ValueString(),
	}, diag.Diagnostics{}
}

// Utils -------------------------------------------------------------------------------------------

// Used to convert structure to a types.Object.
func (self OrchestratorConfigurationStringModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"value": types.StringType,
	}
}

// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// OrchestratorActionInputParameterModel describes the resource data model.
type OrchestratorActionInputParameterModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
}

// OrchestratorActionInputParameterAPIModel describes the resource API model.
type OrchestratorActionInputParameterAPIModel struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

// Used to convert structure to a types.Object.
func OrchestratorActionInputParameterAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        types.StringType,
		"description": types.StringType,
		"type":        types.StringType,
	}
}

func (self OrchestratorActionInputParameterModel) String() string {
	return fmt.Sprintf(
		"Orchestrator Action Input Parameter %s (%s)",
		self.Name.ValueString(),
		self.Type.ValueString())
}

func (self *OrchestratorActionInputParameterModel) FromAPI(
	ctx context.Context,
	raw OrchestratorActionInputParameterAPIModel,
) diag.Diagnostics {
	self.Name = types.StringValue(raw.Name)
	self.Description = types.StringValue(raw.Description)
	self.Type = types.StringValue(raw.Type)
	return diag.Diagnostics{}
}

func (self OrchestratorActionInputParameterModel) ToAPI(
	ctx context.Context,
) (OrchestratorActionInputParameterAPIModel, diag.Diagnostics) {
	return OrchestratorActionInputParameterAPIModel{
		Name: self.Name.ValueString(),
		Description: CleanString(self.Description.ValueString()),
		Type: self.Type.ValueString(),
	}, diag.Diagnostics{}
}

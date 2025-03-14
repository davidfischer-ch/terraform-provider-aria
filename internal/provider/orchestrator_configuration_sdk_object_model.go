// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// OrchestratorConfigurationSDKObjectModel describes the resource data model.
type OrchestratorConfigurationSDKObjectModel struct {
	Id   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

// OrchestratorConfigurationSDKObjectAPIModel describes the resource API model.
type OrchestratorConfigurationSDKObjectAPIModel struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

func (self *OrchestratorConfigurationSDKObjectModel) FromAPI(
	raw OrchestratorConfigurationSDKObjectAPIModel,
) {
	self.Id = types.StringValue(raw.Id)
	self.Type = types.StringValue(raw.Type)
}

func (
	self OrchestratorConfigurationSDKObjectModel,
) ToAPI() OrchestratorConfigurationSDKObjectAPIModel {
	return OrchestratorConfigurationSDKObjectAPIModel{
		Id:   self.Id.ValueString(),
		Type: self.Type.ValueString(),
	}
}

// Utils -------------------------------------------------------------------------------------------

// Used to convert structure to a types.Object.
func (self OrchestratorConfigurationSDKObjectModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"type": types.StringType,
	}
}

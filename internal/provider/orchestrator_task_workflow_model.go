// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// OrchestratorTaskWorkflowModel describes the resource data model.
type OrchestratorTaskWorkflowModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// OrchestratorTaskWorkflowAPIModel describes the resource API model.
type OrchestratorTaskWorkflowAPIModel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (self *OrchestratorTaskWorkflowModel) FromAPI(raw OrchestratorTaskWorkflowAPIModel) {
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
}

func (self OrchestratorTaskWorkflowModel) ToAPI() OrchestratorTaskWorkflowAPIModel {
	return OrchestratorTaskWorkflowAPIModel{
		Id:   self.Id.ValueString(),
		Name: self.Name.ValueString(),
	}
}

// Utils -------------------------------------------------------------------------------------------

// Used to convert structure to a types.Object.
func (self OrchestratorTaskWorkflowModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	}
}

// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CustomNamingTemplateCounterModel describes the resource data model.
type CustomNamingTemplateCounterModel struct {
	Id             types.String `tfsdk:"id"`
	ResourceType   types.String `tfsdk:"resource_type"`
	CurrentCounter types.Int32  `tfsdk:"current_counter"`
	ProjectId      types.String `tfsdk:"project_id"`
}

// CustomNamingTemplateCounterAPIModel describes the resource API model.
type CustomNamingTemplateCounterAPIModel struct {
	Id             string `json:"id"`
	ResourceType   string `json:"cnResourceType"`
	CurrentCounter int32  `json:"currentCounter"`
	ProjectId      string `json:"projectId"`
}

func (self *CustomNamingTemplateCounterModel) FromAPI(raw CustomNamingTemplateCounterAPIModel) {
	self.Id = types.StringValue(raw.Id)
	self.ResourceType = types.StringValue(raw.ResourceType)
	self.CurrentCounter = types.Int32Value(raw.CurrentCounter)
	self.ProjectId = types.StringValue(raw.ProjectId)
}

func (self CustomNamingTemplateCounterModel) ToAPI() CustomNamingTemplateCounterAPIModel {
	return CustomNamingTemplateCounterAPIModel{
		Id:             self.Id.ValueString(),
		ResourceType:   self.ResourceType.ValueString(),
		CurrentCounter: self.CurrentCounter.ValueInt32(),
		ProjectId:      self.ProjectId.ValueString(),
	}
}

// Utils -------------------------------------------------------------------------------------------

// Used to convert structure to a types.Object.
func (self CustomNamingTemplateCounterModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":              types.StringType,
		"resource_type":   types.StringType,
		"current_counter": types.Int32Type,
		"project_id":      types.StringType,
	}
}

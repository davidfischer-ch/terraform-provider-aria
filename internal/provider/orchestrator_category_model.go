// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// OrchestratorCategoryModel describes the resource data model.
type OrchestratorCategoryModel struct {
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Path     types.String `tfsdk:"path"`
	Type     types.String `tfsdk:"type"`
	ParentId types.String `tfsdk:"parent_id"`
}

// OrchestratorCategoryAPIModel describes the resource API model.
type OrchestratorCategoryAPIModel struct {
	Id       string   `json:"id,omitempty"`
	Name     string   `json:"name"`
	Path     string   `json:"path,omitempty"`
	Type     string   `json:"type"`
	ParentId string   `json:"parent-category-id,omitempty"`
	PathIds  []string `json:"path-ids,omitempty"`
}

func (self OrchestratorCategoryModel) String() string {
	return fmt.Sprintf(
		"Orchestrator %s %s (%s)",
		self.Type.ValueString(),
		self.Id.ValueString(),
		self.Name.ValueString())
}

// Return an appropriate key that can be used for naming mutexes.
// Create: Identifier can be used to prevent concurrent creation of vRO actions.
// Read Update Delete: Identifier can be used to prevent concurrent modifications on the instance.
func (self OrchestratorCategoryModel) LockKey() string {
	return "orchestrator-category-" + self.Id.ValueString()
}

func (self OrchestratorCategoryModel) CreatePath() string {
	if len(self.ParentId.ValueString()) == 0 {
		return "vco/api/categories"
	}
	return "vco/api/categories/" + self.ParentId.ValueString()
}

func (self OrchestratorCategoryModel) ReadPath() string {
	return "vco/api/categories/" + self.Id.ValueString()
}

func (self OrchestratorCategoryModel) UpdatePath() string {
	return self.ReadPath()
}

func (self OrchestratorCategoryModel) DeletePath() string {
	return self.ReadPath()
}

func (self *OrchestratorCategoryModel) FromAPI(raw OrchestratorCategoryAPIModel) {
	// Retrieve parent ID from PathIds
	parentId := ""
	parentIndex := len(raw.PathIds) - 2
	if parentIndex >= 0 {
		parentId = raw.PathIds[parentIndex]
	}
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Path = types.StringValue(raw.Path)
	self.Type = types.StringValue(raw.Type)
	self.ParentId = types.StringValue(parentId)
}

func (self OrchestratorCategoryModel) ToAPI() OrchestratorCategoryAPIModel {
	return OrchestratorCategoryAPIModel{
		// ID and Path are computed by Aria
		Name:     self.Name.ValueString(),
		Type:     self.Type.ValueString(),
		ParentId: self.ParentId.ValueString(),
	}
}

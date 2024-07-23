// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CustomNamingTemplateModel describes the resource data model.
type CustomNamingTemplateModel struct {
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	ResourceType     types.String `tfsdk:"resource_type"`
	ResourceTypeName types.String `tfsdk:"resource_type_name"`
	ResourceDefault  types.Bool   `tfsdk:"resource_default"`
	Pattern          types.String `tfsdk:"pattern"`
	StaticPattern    types.String `tfsdk:"static_pattern"`
	UniqueName       types.Bool   `tfsdk:"unique_name"`
	StartCounter     types.Int32  `tfsdk:"start_counter"`
	IncrementStep    types.Int32  `tfsdk:"increment_step"`
}

// CustomNamingTemplateAPIModel describes the resource API model.
type CustomNamingTemplateAPIModel struct {
	Id               string `json:"id,omitempty"`
	Name             string `json:"name,omitempty"`
	ResourceType     string `json:"resourceType"`
	ResourceTypeName string `json:"resourceTypeName"`
	ResourceDefault  bool   `json:"resourceDefault"`
	UniqueName       bool   `json:"uniqueName"`
	Pattern          string `json:"pattern"`
	StaticPattern    string `json:"staticPattern"`
	StartCounter     int32  `json:"startCounter"`
	IncrementStep    int32  `json:"incrementStep"`
}

func (self CustomNamingTemplateModel) String() string {
	return fmt.Sprintf("Custom Naming Template %s (%s)", self.Id.ValueString(), self.Key())
}

func (self CustomNamingTemplateModel) Key() string {
	pattern := self.StaticPattern.ValueString()
	if len(pattern) == 0 {
		pattern = "Default"
	}
	return fmt.Sprintf(
		"%s.%s > %s",
		self.ResourceType.ValueString(),
		self.ResourceTypeName.ValueString(),
		pattern)
}

func (self *CustomNamingTemplateModel) FromAPI(
	ctx context.Context,
	raw CustomNamingTemplateAPIModel,
) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.ResourceType = types.StringValue(raw.ResourceType)
	self.ResourceTypeName = types.StringValue(raw.ResourceTypeName)
	self.ResourceDefault = types.BoolValue(raw.ResourceDefault)
	self.UniqueName = types.BoolValue(raw.UniqueName)
	self.Pattern = types.StringValue(raw.Pattern)
	self.StaticPattern = types.StringValue(raw.StaticPattern)
	self.StartCounter = types.Int32Value(raw.StartCounter)
	self.IncrementStep = types.Int32Value(raw.IncrementStep)
	return diag.Diagnostics{}
}

func (self CustomNamingTemplateModel) toAPI() CustomNamingTemplateAPIModel {
	return CustomNamingTemplateAPIModel{
		Id:               self.Id.ValueString(),
		Name:             self.Name.ValueString(),
		ResourceType:     self.ResourceType.ValueString(),
		ResourceTypeName: self.ResourceTypeName.ValueString(),
		ResourceDefault:  len(self.StaticPattern.ValueString()) == 0,
		UniqueName:       self.UniqueName.ValueBool(),
		Pattern:          self.Pattern.ValueString(),
		StaticPattern:    self.StaticPattern.ValueString(),
		StartCounter:     self.StartCounter.ValueInt32(),
		IncrementStep:    self.IncrementStep.ValueInt32(),
	}
}

func (self CustomNamingTemplateModel) ToAPI(
	ctx context.Context,
	state CustomNamingTemplateModel,
) (CustomNamingTemplateAPIModel, diag.Diagnostics) {
	raw := self.toAPI()
	// If the identifier is set, means its an UPDATE
	if len(raw.Id) > 0 {
		stateRaw := state.toAPI()
		// Attributes are writable once, any changes requires a replacement
		// In that case, the identifier is wiped to trigger the replacement (by Aria)
		if raw.Name != stateRaw.Name ||
			raw.ResourceType != stateRaw.ResourceType ||
			raw.ResourceTypeName != stateRaw.ResourceTypeName ||
			raw.ResourceDefault != stateRaw.ResourceDefault ||
			raw.UniqueName != stateRaw.UniqueName ||
			raw.Pattern != stateRaw.Pattern ||
			raw.StaticPattern != stateRaw.StaticPattern ||
			raw.StartCounter != stateRaw.StartCounter ||
			raw.IncrementStep != stateRaw.IncrementStep {
			raw.Id = ""
		}
	}
	return raw, diag.Diagnostics{}
}

// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

	// Of type CustomNamingTemplateCounterModel
	Counters types.List `tfsdk:"counters"`
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

	Counters []CustomNamingTemplateCounterAPIModel `json:"counters,omitempty"`
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

	diags := diag.Diagnostics{}

	// Convert counters from raw to list
	attrs := types.ObjectType{AttrTypes: CustomNamingTemplateCounterModel{}.AttributeTypes()}
	if raw.Counters == nil {
		self.Counters = types.ListNull(attrs)
	} else {
		counters := []CustomNamingTemplateCounterModel{}
		for _, counterRaw := range raw.Counters {
			counter := CustomNamingTemplateCounterModel{}
			counter.FromAPI(counterRaw)
			counters = append(counters, counter)
		}

		var someDiags diag.Diagnostics
		self.Counters, someDiags = types.ListValueFrom(ctx, attrs, counters)
		diags.Append(someDiags...)
	}

	return diags
}

func (self CustomNamingTemplateModel) toAPI(
	ctx context.Context,
) (CustomNamingTemplateAPIModel, diag.Diagnostics) {

	diags := diag.Diagnostics{}

	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/list
	var countersRaw []CustomNamingTemplateCounterAPIModel
	if self.Counters.IsUnknown() || self.Counters.IsNull() {
		countersRaw = nil
	} else {
		// Extract counters from list value and then convert to raw
		counters := make([]CustomNamingTemplateCounterModel, 0, len(self.Counters.Elements()))
		diags.Append(self.Counters.ElementsAs(ctx, &counters, false)...)
		if !diags.HasError() {
			for _, counter := range counters {
				countersRaw = append(countersRaw, counter.ToAPI())
			}
		}
	}

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
		Counters:         countersRaw,
	}, diags
}

func (self CustomNamingTemplateModel) ToAPI(
	ctx context.Context,
	state CustomNamingTemplateModel,
) (CustomNamingTemplateAPIModel, diag.Diagnostics) {
	raw, diags := self.toAPI(ctx)

	stateRaw, someDiags := state.toAPI(ctx)
	diags.Append(someDiags...)
	if diags.HasError() {
		return raw, diags
	}

	// Attributes are writable once, any changes requires a replacement
	// In that case, the identifier is wiped to trigger the replacement (by Aria)
	if raw.Name == stateRaw.Name &&
		raw.ResourceType == stateRaw.ResourceType &&
		raw.ResourceTypeName == stateRaw.ResourceTypeName &&
		raw.ResourceDefault == stateRaw.ResourceDefault &&
		raw.UniqueName == stateRaw.UniqueName &&
		raw.Pattern == stateRaw.Pattern &&
		raw.StaticPattern == stateRaw.StaticPattern &&
		raw.StartCounter == stateRaw.StartCounter &&
		raw.IncrementStep == stateRaw.IncrementStep {
		// Keep last known identifier and counters
		raw.Id = stateRaw.Id
		raw.Counters = stateRaw.Counters
		tflog.Debug(ctx,fmt.Sprintf("Keep last known %s ID and counters", self.String()))
	} else {
		// Wipe identifier and counters
		raw.Id = ""
		raw.Counters = []CustomNamingTemplateCounterAPIModel{}
		tflog.Debug(ctx, fmt.Sprintf("Wipe %s ID and counters", self.String()))
	}
	return raw, diags
}

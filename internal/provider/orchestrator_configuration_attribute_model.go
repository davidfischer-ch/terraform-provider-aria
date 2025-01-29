// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// OrchestratorConfigurationAttributeModel describes the resource data model.
type OrchestratorConfigurationAttributeModel struct {
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`

	// Of type OrchestratorConfigurationValueModel
	Value types.Object `tfsdk:"value"`
}

// OrchestratorConfigurationAttributeAPIModel describes the resource API model.
type OrchestratorConfigurationAttributeAPIModel struct {
	Name string `json:"name"`
	Type string `json:"type"`

	Value OrchestratorConfigurationValueAPIModel `json:"value"`
}

func (self OrchestratorConfigurationAttributeModel) String() string {
	return fmt.Sprintf(
		"Orchestrator Configuration Attribute %s (%s)",
		self.Name.ValueString(),
		self.Type.ValueString())
}

func (self *OrchestratorConfigurationAttributeModel) FromAPI(
	ctx context.Context,
	raw OrchestratorConfigurationAttributeAPIModel,
) diag.Diagnostics {
	self.Name = types.StringValue(raw.Name)
	self.Type = types.StringValue(raw.Type)

	// Convert value from raw and then to object
	var someDiags diag.Diagnostics
	value := OrchestratorConfigurationValueModel{}
	diags := value.FromAPI(ctx, raw.Value)
	self.Value, someDiags = types.ObjectValueFrom(ctx, value.AttributeTypes(), value)
	diags.Append(someDiags...)

	return diags
}

func (self OrchestratorConfigurationAttributeModel) ToAPI(
	ctx context.Context,
) (OrchestratorConfigurationAttributeAPIModel, diag.Diagnostics) {

	diags := diag.Diagnostics{}
	valueRaw := OrchestratorConfigurationValueAPIModel{}

	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/object
	if self.Value.IsNull() || self.Value.IsUnknown() {
		diags.AddError(
			"Configuration error",
			fmt.Sprintf("Unable to manage %s, value is either null or unknown", self.String()))
	} else {
		// Convert value from object to raw
		var someDiags diag.Diagnostics
		value := OrchestratorConfigurationValueModel{}
		diags.Append(self.Value.As(ctx, &value, basetypes.ObjectAsOptions{})...)
		valueRaw, someDiags = value.ToAPI(ctx)
		diags.Append(someDiags...)
	}

	return OrchestratorConfigurationAttributeAPIModel{
		Name:  self.Name.ValueString(),
		Type:  self.Type.ValueString(),
		Value: valueRaw,
	}, diags
}

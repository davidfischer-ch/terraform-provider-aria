// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// OrchestratorConfigurationValueModel describes the resource data model.
type OrchestratorConfigurationValueModel struct {
	Boolean types.Object `tfsdk:"boolean"` // Of type OrchestratorConfigurationBooleanModel
	String  types.Object `tfsdk:"string"`  // Of type OrchestratorConfigurationStringModel
}

// OrchestratorConfigurationValueAPIModel describes the resource API model.
type OrchestratorConfigurationValueAPIModel struct {
	Boolean *OrchestratorConfigurationBooleanAPIModel `json:"boolean,omitempty"`
	String  *OrchestratorConfigurationStringAPIModel  `json:"string,omitempty"`
}

func (self *OrchestratorConfigurationValueModel) FromAPI(
	ctx context.Context,
	raw OrchestratorConfigurationValueAPIModel,
) diag.Diagnostics {

	diags := diag.Diagnostics{}

	// Convert boolean from raw and then to object
	boolean := OrchestratorConfigurationBooleanModel{}
	if raw.Boolean == nil {
		self.Boolean = types.ObjectNull(boolean.AttributeTypes())
	} else {
		var someDiags diag.Diagnostics
		diags := boolean.FromAPI(ctx, *raw.Boolean)
		self.Boolean, someDiags = types.ObjectValueFrom(ctx, boolean.AttributeTypes(), boolean)
		diags.Append(someDiags...)
	}

	// Convert string from raw and then to object
	string := OrchestratorConfigurationStringModel{}
	if raw.String == nil {
		self.String = types.ObjectNull(string.AttributeTypes())
	} else {
		var someDiags diag.Diagnostics
		diags := string.FromAPI(ctx, *raw.String)
		self.String, someDiags = types.ObjectValueFrom(ctx, string.AttributeTypes(), string)
		diags.Append(someDiags...)
	}

	return diags
}

func (self OrchestratorConfigurationValueModel) ToAPI(
	ctx context.Context,
) (OrchestratorConfigurationValueAPIModel, diag.Diagnostics) {

	var diags diag.Diagnostics
	raw := OrchestratorConfigurationValueAPIModel{}

	if self.Boolean.IsNull() || self.Boolean.IsUnknown() {
		raw.Boolean = nil
	} else {
		var boolean OrchestratorConfigurationBooleanModel
		diags.Append(self.Boolean.As(ctx, &boolean, basetypes.ObjectAsOptions{})...)
		booleanRaw, someDiags := boolean.ToAPI(ctx)
		raw.Boolean = &booleanRaw
		diags.Append(someDiags...)
	}

	if self.String.IsNull() || self.String.IsUnknown() {
		raw.String = nil
	} else {
		var string OrchestratorConfigurationStringModel
		diags.Append(self.String.As(ctx, &string, basetypes.ObjectAsOptions{})...)
		stringRaw, someDiags := string.ToAPI(ctx)
		raw.String = &stringRaw
		diags.Append(someDiags...)
	}

	return raw, diags
}

// Utils -------------------------------------------------------------------------------------------

// Used to convert structure to a types.Object.
func (self OrchestratorConfigurationValueModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"boolean": types.ObjectType{
			AttrTypes: OrchestratorConfigurationBooleanModel{}.AttributeTypes(),
		},
		"string": types.ObjectType{
			AttrTypes: OrchestratorConfigurationStringModel{}.AttributeTypes(),
		},
	}
}

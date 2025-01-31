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
	boolVal := OrchestratorConfigurationBooleanModel{}
	if raw.Boolean == nil {
		self.Boolean = types.ObjectNull(boolVal.AttributeTypes())
	} else {
		var someDiags diag.Diagnostics
		boolVal.FromAPI(*raw.Boolean)
		self.Boolean, someDiags = types.ObjectValueFrom(ctx, boolVal.AttributeTypes(), boolVal)
		diags.Append(someDiags...)
	}

	// Convert string from raw and then to object
	stringVal := OrchestratorConfigurationStringModel{}
	if raw.String == nil {
		self.String = types.ObjectNull(stringVal.AttributeTypes())
	} else {
		var someDiags diag.Diagnostics
		stringVal.FromAPI(*raw.String)
		self.String, someDiags = types.ObjectValueFrom(ctx, stringVal.AttributeTypes(), stringVal)
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
		var boolVal OrchestratorConfigurationBooleanModel
		diags.Append(self.Boolean.As(ctx, &boolVal, basetypes.ObjectAsOptions{})...)
		boolRaw := boolVal.ToAPI()
		raw.Boolean = &boolRaw
	}

	if self.String.IsNull() || self.String.IsUnknown() {
		raw.String = nil
	} else {
		var stringVal OrchestratorConfigurationStringModel
		diags.Append(self.String.As(ctx, &stringVal, basetypes.ObjectAsOptions{})...)
		stringRaw := stringVal.ToAPI()
		raw.String = &stringRaw
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

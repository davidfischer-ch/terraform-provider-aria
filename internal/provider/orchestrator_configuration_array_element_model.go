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

// OrchestratorConfigurationArrayElementModel describes the resource data model.
type OrchestratorConfigurationArrayElementModel struct {
	Boolean      types.Object `tfsdk:"boolean"`       // Of type ...BooleanModel
	String       types.Object `tfsdk:"string"`        // Of type ...StringModel
	SecureString types.Object `tfsdk:"secure_string"` // Of type ...SecureStringModel
	SDKObject    types.Object `tfsdk:"sdk_object"`    // Of type ...SDKObjectModel
}

// OrchestratorConfigurationArrayElementAPIModel describes the resource API model.
type OrchestratorConfigurationArrayElementAPIModel struct {
	Boolean      *OrchestratorConfigurationBooleanAPIModel      `json:"boolean,omitempty"`
	String       *OrchestratorConfigurationStringAPIModel       `json:"string,omitempty"`
	SecureString *OrchestratorConfigurationSecureStringAPIModel `json:"secure-string,omitempty"`
	SDKObject    *OrchestratorConfigurationSDKObjectAPIModel    `json:"sdk-object,omitempty"`
}

func (self *OrchestratorConfigurationArrayElementModel) FromAPI(
	ctx context.Context,
	raw OrchestratorConfigurationArrayElementAPIModel,
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

	// Convert secure string from raw and then to object
	secureVal := OrchestratorConfigurationSecureStringModel{}
	if raw.SecureString == nil {
		self.SecureString = types.ObjectNull(secureVal.AttributeTypes())
	} else {
		var someDiags diag.Diagnostics
		secureVal.FromAPI(*raw.SecureString)
		self.SecureString, someDiags = types.ObjectValueFrom(
			ctx, secureVal.AttributeTypes(), secureVal,
		)
		diags.Append(someDiags...)
	}

	// Convert secure string from raw and then to object
	objVal := OrchestratorConfigurationSDKObjectModel{}
	if raw.SDKObject == nil {
		self.SDKObject = types.ObjectNull(objVal.AttributeTypes())
	} else {
		var someDiags diag.Diagnostics
		objVal.FromAPI(*raw.SDKObject)
		self.SDKObject, someDiags = types.ObjectValueFrom(ctx, objVal.AttributeTypes(), objVal)
		diags.Append(someDiags...)
	}

	return diags
}

func (self OrchestratorConfigurationArrayElementModel) ToAPI(
	ctx context.Context,
) (OrchestratorConfigurationArrayElementAPIModel, diag.Diagnostics) {

	var diags diag.Diagnostics
	raw := OrchestratorConfigurationArrayElementAPIModel{}

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

	if self.SecureString.IsNull() || self.SecureString.IsUnknown() {
		raw.SecureString = nil
	} else {
		var secureVal OrchestratorConfigurationSecureStringModel
		diags.Append(self.SecureString.As(ctx, &secureVal, basetypes.ObjectAsOptions{})...)
		secureRaw := secureVal.ToAPI()
		raw.SecureString = &secureRaw
	}

	if self.SDKObject.IsNull() || self.SDKObject.IsUnknown() {
		raw.SDKObject = nil
	} else {
		var objVal OrchestratorConfigurationSDKObjectModel
		diags.Append(self.SDKObject.As(ctx, &objVal, basetypes.ObjectAsOptions{})...)
		objRaw := objVal.ToAPI()
		raw.SDKObject = &objRaw
	}

	return raw, diags
}

// Utils -------------------------------------------------------------------------------------------

// Used to convert structure to a types.Object.
func (self OrchestratorConfigurationArrayElementModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"boolean": types.ObjectType{
			AttrTypes: OrchestratorConfigurationBooleanModel{}.AttributeTypes(),
		},
		"string": types.ObjectType{
			AttrTypes: OrchestratorConfigurationStringModel{}.AttributeTypes(),
		},
		"secure_string": types.ObjectType{
			AttrTypes: OrchestratorConfigurationSecureStringModel{}.AttributeTypes(),
		},
		"sdk_object": types.ObjectType{
			AttrTypes: OrchestratorConfigurationSDKObjectModel{}.AttributeTypes(),
		},
	}
}

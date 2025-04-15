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

// OrchestratorConfigurationAttributeValueModel describes the resource data model.
type OrchestratorConfigurationAttributeValueModel struct {
	Array        types.Object `tfsdk:"array"`         // Of type ...ArrayModel
	Boolean      types.Object `tfsdk:"boolean"`       // Of type ...BooleanModel
	Number       types.Object `tfsdk:"number"`        // Of type ...NumberModel
	String       types.Object `tfsdk:"string"`        // Of type ...StringModel
	SecureString types.Object `tfsdk:"secure_string"` // Of type ...SecureStringModel
	SDKObject    types.Object `tfsdk:"sdk_object"`    // Of type ...SDKObjectModel
}

// OrchestratorConfigurationAttributeValueAPIModel describes the resource API model.
type OrchestratorConfigurationAttributeValueAPIModel struct {
	Array        *OrchestratorConfigurationArrayAPIModel        `json:"array,omitempty"`
	Boolean      *OrchestratorConfigurationBooleanAPIModel      `json:"boolean,omitempty"`
	Number       *OrchestratorConfigurationNumberAPIModel       `json:"number,omitempty"`
	String       *OrchestratorConfigurationStringAPIModel       `json:"string,omitempty"`
	SecureString *OrchestratorConfigurationSecureStringAPIModel `json:"secure-string,omitempty"`
	SDKObject    *OrchestratorConfigurationSDKObjectAPIModel    `json:"sdk-object,omitempty"`
}

func (self *OrchestratorConfigurationAttributeValueModel) FromAPI(
	ctx context.Context,
	raw OrchestratorConfigurationAttributeValueAPIModel,
) diag.Diagnostics {

	diags := diag.Diagnostics{}

	// Convert array from raw and then to object
	arrayVal := OrchestratorConfigurationArrayModel{}
	if raw.Array == nil {
		self.Array = types.ObjectNull(arrayVal.AttributeTypes())
	} else {
		var someDiags diag.Diagnostics
		diags.Append(arrayVal.FromAPI(ctx, *raw.Array)...)
		self.Array, someDiags = types.ObjectValueFrom(ctx, arrayVal.AttributeTypes(), arrayVal)
		diags.Append(someDiags...)
	}

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

	// Convert number from raw and then to object
	numberVal := OrchestratorConfigurationNumberModel{}
	if raw.Number == nil {
		self.Number = types.ObjectNull(numberVal.AttributeTypes())
	} else {
		var someDiags diag.Diagnostics
		numberVal.FromAPI(*raw.Number)
		self.Number, someDiags = types.ObjectValueFrom(ctx, numberVal.AttributeTypes(), numberVal)
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

func (self OrchestratorConfigurationAttributeValueModel) ToAPI(
	ctx context.Context,
) (OrchestratorConfigurationAttributeValueAPIModel, diag.Diagnostics) {

	var diags diag.Diagnostics
	raw := OrchestratorConfigurationAttributeValueAPIModel{}

	if self.Array.IsNull() || self.Array.IsUnknown() {
		raw.Array = nil
	} else {
		var arrayVal OrchestratorConfigurationArrayModel
		var someDiags diag.Diagnostics
		diags.Append(self.Array.As(ctx, &arrayVal, basetypes.ObjectAsOptions{})...)
		arrayRaw, someDiags := arrayVal.ToAPI(ctx)
		raw.Array = &arrayRaw
		diags.Append(someDiags...)
	}

	if self.Boolean.IsNull() || self.Boolean.IsUnknown() {
		raw.Boolean = nil
	} else {
		var boolVal OrchestratorConfigurationBooleanModel
		diags.Append(self.Boolean.As(ctx, &boolVal, basetypes.ObjectAsOptions{})...)
		boolRaw := boolVal.ToAPI()
		raw.Boolean = &boolRaw
	}

	if self.Number.IsNull() || self.Number.IsUnknown() {
		raw.Number = nil
	} else {
		var numberVal OrchestratorConfigurationNumberModel
		diags.Append(self.Number.As(ctx, &numberVal, basetypes.ObjectAsOptions{})...)
		numberRaw := numberVal.ToAPI()
		raw.Number = &numberRaw
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
func (self OrchestratorConfigurationAttributeValueModel) AttributeTypes(
	ctx context.Context,
) map[string]attr.Type {
	return map[string]attr.Type{
		"array": types.ObjectType{
			AttrTypes: OrchestratorConfigurationArrayModel{}.AttributeTypes(),
		},
		"boolean": types.ObjectType{
			AttrTypes: OrchestratorConfigurationBooleanModel{}.AttributeTypes(),
		},
		"number": types.ObjectType{
			AttrTypes: OrchestratorConfigurationNumberModel{}.AttributeTypes(),
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

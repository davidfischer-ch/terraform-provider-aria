// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CustomResourcPropertyeModel describes the resource data model.
type PropertyModel struct {
	Name             types.String         `tfsdk:"name"`
	Title            types.String         `tfsdk:"title"`
	Description      types.String         `tfsdk:"description"`
	Type             types.String         `tfsdk:"type"`
	Default          jsontypes.Normalized `tfsdk:"default"`
	Encrypted        types.Bool           `tfsdk:"encrypted"`
	ReadOnly         types.Bool           `tfsdk:"read_only"`
	RecreateOnUpdate types.Bool           `tfsdk:"recreate_on_update"`

	// Specifications
	Minimum   types.Int64  `tfsdk:"minimum"`
	Maximum   types.Int64  `tfsdk:"maximum"`
	MinLength types.Int32  `tfsdk:"min_length"`
	MaxLength types.Int32  `tfsdk:"max_length"`
	Pattern   types.String `tfsdk:"pattern"`
	/*Items*/
	OneOf []PropertyOneOfModel `tfsdk:"one_of"`
}

// PropertyAPIModel describes the resource API model.
type PropertyAPIModel struct {
	Title            string `json:"title" yaml:"title"`
	Description      string `json:"description" yaml:"description"`
	Type             string `json:"type" yaml:"type"`
	Default          any    `json:"default,omitempty" yaml:"default"`
	Encrypted        bool   `json:"encrypted" yaml:"encrypted"`
	ReadOnly         bool   `json:"readOnly" yaml:"readOnly"`
	RecreateOnUpdate bool   `json:"recreateOnUpdate" yaml:"recreateOnUpdate"`

	// Specifications
	Minimum   *int64  `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	Maximum   *int64  `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	MinLength *int32  `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	MaxLength *int32  `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	Pattern   *string `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	/*Items*/
	OneOf []PropertyOneOfAPIModel `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
}

func (self PropertyModel) String() string {
	return fmt.Sprintf("Property %s", self.Name.ValueString())
}

func (self *PropertyModel) FromAPI(
	ctx context.Context,
	name string,
	raw PropertyAPIModel,
) diag.Diagnostics {

	self.Name = types.StringValue(name)
	self.Title = types.StringValue(raw.Title)
	self.Description = types.StringValue(raw.Description)
	self.Type = types.StringValue(raw.Type)
	self.Encrypted = types.BoolValue(raw.Encrypted)
	self.ReadOnly = types.BoolValue(raw.ReadOnly)
	self.RecreateOnUpdate = types.BoolValue(raw.RecreateOnUpdate)
	self.Minimum = types.Int64PointerValue(raw.Minimum)
	self.Maximum = types.Int64PointerValue(raw.Maximum)
	self.MinLength = types.Int32PointerValue(raw.MinLength)
	self.MaxLength = types.Int32PointerValue(raw.MaxLength)
	self.Pattern = types.StringPointerValue(raw.Pattern)

	diags := diag.Diagnostics{}

	// Default API data -> JSON Encoded
	// TODO Deduplicate this routine, used on many places
	if raw.Default == nil {
		self.Default = jsontypes.NewNormalizedNull()
	} else {
		defaultJSON, err := json.Marshal(raw.Default)
		if err != nil {
			diags.AddError(
				"Client error",
				fmt.Sprintf(
					"Unable to JSON encode %s default \"%s\", got error: %s",
					self.String(), raw.Default, err))
		} else {
			self.Default = jsontypes.NewNormalizedValue(string(defaultJSON))
		}
	}

	if raw.OneOf == nil {
		self.OneOf = nil
	} else {
		self.OneOf = []PropertyOneOfModel{}
		for _, oneOfRaw := range raw.OneOf {
			oneOf := PropertyOneOfModel{}
			diags.Append(oneOf.FromAPI(ctx, oneOfRaw)...)
			self.OneOf = append(self.OneOf, oneOf)
		}
	}

	return diags
}

func (self PropertyModel) ToAPI(
	ctx context.Context,
) (string, PropertyAPIModel, diag.Diagnostics) {

	diags := diag.Diagnostics{}

	// Criteria JSON Encoded -> API data
	// TODO Deduplicate this routine, used on many places
	var defaultRaw any
	if self.Default.IsNull() {
		defaultRaw = nil
	} else {
		diags.Append(self.Default.Unmarshal(&defaultRaw)...)
	}

	var oneOfRawList []PropertyOneOfAPIModel
	if self.OneOf == nil {
		oneOfRawList = nil
	} else {
		oneOfRawList = []PropertyOneOfAPIModel{}
		for _, oneOf := range self.OneOf {
			oneOfRaw, oneOfDiags := oneOf.ToAPI(ctx)
			oneOfRawList = append(oneOfRawList, oneOfRaw)
			diags.Append(oneOfDiags...)
		}
	}

	return self.Name.ValueString(),
		PropertyAPIModel{
			Title:            self.Title.ValueString(),
			Description:      self.Description.ValueString(),
			Type:             self.Type.ValueString(),
			Default:          defaultRaw,
			Encrypted:        self.Encrypted.ValueBool(),
			ReadOnly:         self.ReadOnly.ValueBool(),
			RecreateOnUpdate: self.RecreateOnUpdate.ValueBool(),
			Minimum:          self.Minimum.ValueInt64Pointer(),
			Maximum:          self.Maximum.ValueInt64Pointer(),
			MinLength:        self.MinLength.ValueInt32Pointer(),
			MaxLength:        self.MaxLength.ValueInt32Pointer(),
			Pattern:          self.Pattern.ValueStringPointer(),
			OneOf:            oneOfRawList,
		},
		diags
}

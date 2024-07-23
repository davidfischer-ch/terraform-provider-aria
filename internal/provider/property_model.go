// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CustomResourcPropertyeModel describes the resource data model.
type PropertyModel struct {
	Name             types.String `tfsdk:"name"`
	Title            types.String `tfsdk:"title"`
	Description      types.String `tfsdk:"description"`
	Type             types.String `tfsdk:"type"`
	Default          types.String `tfsdk:"default"`
	Encrypted        types.Bool   `tfsdk:"encrypted"`
	ReadOnly         types.Bool   `tfsdk:"read_only"`
	RecreateOnUpdate types.Bool   `tfsdk:"recreate_on_update"`

	// Specifications
	Minimum   types.Int64          `tfsdk:"minimum"`
	Maximum   types.Int64          `tfsdk:"maximum"`
	MinLength types.Int32          `tfsdk:"min_length"`
	MaxLength types.Int32          `tfsdk:"max_length"`
	Pattern   types.String         `tfsdk:"pattern"`
	OneOf     []PropertyOneOfModel `tfsdk:"one_of"`
}

// PropertyAPIModel describes the resource API model.
type PropertyAPIModel struct {
	Title            string `json:"title"`
	Description      string `json:"description"`
	Type             string `json:"type"`
	Default          any    `json:"default,omitempty"`
	Encrypted        bool   `json:"encrypted"`
	ReadOnly         bool   `tfsdk:"readOnly"`
	RecreateOnUpdate bool   `json:"recreateOnUpdate"`

	// Specifications
	Minimum   *int64                  `json:"minimum,omitempty"`
	Maximum   *int64                  `json:"maximum,omitempty"`
	MinLength *int32                  `json:"minLength,omitempty"`
	MaxLength *int32                  `json:"maxLength,omitempty"`
	Pattern   string                  `json:"pattern"`
	OneOf     []PropertyOneOfAPIModel `json:"oneOf,omitempty"`
}

func (self PropertyModel) String() string {
	return fmt.Sprintf("Property %s", self.Name.ValueString())
}

func (self *PropertyModel) FromAPI(
	ctx context.Context,
	name string,
	raw PropertyAPIModel,
) diag.Diagnostics {

	diags := diag.Diagnostics{}

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
	self.Pattern = types.StringValue(raw.Pattern)

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

	if raw.Default == nil {
		self.Default = types.StringNull()
	} else {
		// Convert default value from any to string, warn if default type mismatch
		self.Default = types.StringValue(fmt.Sprintf("%s", raw.Default)) // Messy conversion first
		switch raw.Type {
		case "boolean":
			// Must be a boolean
			if defaultBool, ok := raw.Default.(bool); ok {
				self.Default = types.StringValue(strconv.FormatBool(defaultBool))
			} else {
				diags.AddWarning(
					"Configuration warning",
					fmt.Sprintf(
						"Property %s default \"%s\" is not a boolean",
						raw.Title, raw.Default))
			}
		case "integer":
			// Must be an integer
			if defaultInt, ok := raw.Default.(int64); ok {
				self.Default = types.StringValue(strconv.FormatInt(defaultInt, 10))
			} else {
				diags.AddWarning(
					"Configuration warning",
					fmt.Sprintf(
						"Property %s default \"%s\" is not an integer",
						raw.Title, raw.Default))
			}
		case "number":
			// Try integer first, then float
			if defaultInt, ok := raw.Default.(int64); ok {
				self.Default = types.StringValue(strconv.FormatInt(defaultInt, 10))
			} else if defaultFloat, ok := raw.Default.(float64); ok {
				self.Default = types.StringValue(strconv.FormatFloat(defaultFloat, 'g', -1, 64))
			} else {
				diags.AddWarning(
					"Configuration warning",
					fmt.Sprintf(
						"Property %s default \"%s\" is not a number",
						raw.Title, raw.Default))
			}
		case "string":
			// Must be a string
			if defaultString, ok := raw.Default.(string); ok {
				self.Default = types.StringValue(defaultString)
			} else {
				diags.AddWarning(
					"Configuration warning",
					fmt.Sprintf(
						"Property %s default \"%s\" is not a string",
						raw.Title, raw.Default))
			}
		default:
			// Not implemented or wrong type
			diags.AddError(
				"Configuration error",
				fmt.Sprintf(
					"Managing property %s of type %s is not yet implemented.",
					raw.Title, raw.Type))
		}
	}
	return diags
}

func (self PropertyModel) ToAPI(
	ctx context.Context,
) (string, PropertyAPIModel, diag.Diagnostics) {

	diags := diag.Diagnostics{}

	// Convert default value string to appropriate type
	titleRaw := self.Title.ValueString()
	typeRaw := self.Type.ValueString()

	var defaultRaw any
	if self.Default.IsNull() || self.Default.IsUnknown() {
		defaultRaw = nil
	} else {
		var err error
		defaultString := self.Default.ValueString()
		switch typeRaw {
		case "boolean":
			// Must be a boolean
			defaultRaw, err = strconv.ParseBool(defaultString)
		case "integer":
			// Must be an ineger
			defaultRaw, err = strconv.ParseInt(defaultString, 10, 64)
		case "number":
			// Try integer first, then float
			if defaultRaw, err = strconv.ParseInt(defaultString, 10, 64); err != nil {
				defaultRaw, err = strconv.ParseFloat(defaultString, 64)
			}
		case "string":
			// Nothing to do
			defaultRaw = defaultString
			err = nil
		default:
			// Not implemented or wrong type
			diags.AddError(
				"Configuration error",
				fmt.Sprintf(
					"Managing property %s of type %s is not yet implemented.",
					titleRaw, typeRaw))
		}
		if err != nil {
			diags.AddError(
				"Configuration error",
				fmt.Sprintf(
					"Unable to convert property %s default value \"%s\" to type %s, got error: %s",
					titleRaw, defaultString, typeRaw, err))
		}
		if diags.HasError() {
			return self.Name.ValueString(), PropertyAPIModel{}, diags
		}
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
			Title:            titleRaw,
			Description:      self.Description.ValueString(),
			Type:             typeRaw,
			Default:          defaultRaw,
			Encrypted:        self.Encrypted.ValueBool(),
			ReadOnly:         self.ReadOnly.ValueBool(),
			RecreateOnUpdate: self.RecreateOnUpdate.ValueBool(),
			Minimum:          self.Minimum.ValueInt64Pointer(),
			Maximum:          self.Maximum.ValueInt64Pointer(),
			MinLength:        self.MinLength.ValueInt32Pointer(),
			MaxLength:        self.MaxLength.ValueInt32Pointer(),
			Pattern:          self.Pattern.ValueString(),
			OneOf:            oneOfRawList,
		},
		diags
}

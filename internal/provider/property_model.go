// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CustomResourcPropertyeModel describes the resource data model.
type PropertyModel struct {
	Title            types.String `tfsdk:"title"`
	Description      types.String `tfsdk:"description"`
	Type             types.String `tfsdk:"type"`
	Default          types.String `tfsdk:"default"`
	Encrypted        types.Bool   `tfsdk:"encrypted"`
	RecreateOnUpdate types.Bool   `tfsdk:"recreate_on_update"`

	// Specifications
	Minimum   types.Int64  `tfsdk:"minimum"`
	Maximum   types.Int64  `tfsdk:"maximum"`
	MinLength types.Int64  `tfsdk:"min_length"`
	MaxLength types.Int64  `tfsdk:"max_length"`
	Pattern   types.String `tfsdk:"pattern"`
	OneOf     types.List   `tfsdk:"one_of"` // Of type PropertyOneOfModel
}

// PropertyAPIModel describes the resource API model.
type PropertyAPIModel struct {
	Title            string `json:"title"`
	Description      string `json:"description"`
	Type             string `json:"type"`
	Default          any    `json:"default"`
	Encrypted        bool   `json:"encrypted"`
	RecreateOnUpdate bool   `json:"recreateOnUpdate"`

	// Specifications
	Minimum   int64                   `json:"minimum"`
	Maximum   int64                   `json:"maximum"`
	MinLength int64                   `json:"minLength"`
	MaxLength int64                   `json:"maxLength"`
	Pattern   string                  `json:"pattern"`
	OneOf     []PropertyOneOfAPIModel `json:"oneOf"`
}

func (self *PropertyModel) FromAPI(
	ctx context.Context,
	raw PropertyAPIModel,
) diag.Diagnostics {

	diags := diag.Diagnostics{}

	self.Title = types.StringValue(raw.Title)
	self.Description = types.StringValue(raw.Description)
	self.Type = types.StringValue(raw.Type)
	self.Encrypted = types.BoolValue(raw.Encrypted)
	self.RecreateOnUpdate = types.BoolValue(raw.RecreateOnUpdate)
	self.Minimum = types.Int64Value(raw.Minimum)
	self.Maximum = types.Int64Value(raw.Maximum)
	self.MinLength = types.Int64Value(raw.MinLength)
	self.MaxLength = types.Int64Value(raw.MaxLength)
	self.Pattern = types.StringValue(raw.Pattern)
	// FIXME self.OneOf =

	defaultJSON, err := json.Marshal(raw.Default)
	if err != nil {
		diags.AddError(
			"Internal error",
			fmt.Sprintf(
				"Unable to JSON encode default value for property %s, got error: %s",
				raw.Title, err))
		return diags
	}

	self.Default = types.StringValue(string(defaultJSON))

	return diags
}

func (self *PropertyModel) ToAPI(
	ctx context.Context,
	resource string,
) (PropertyAPIModel, diag.Diagnostics) {

	diags := diag.Diagnostics{}

	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/list
	/*if self.OneOf.IsNull() || self.OneOf.IsUnknown() {
	    diags.AddError(
	        "Configuration error",
	        fmt.Sprintf("Unable to manage %s, one_of is either null or unknown", resource))
	    return PropertyAPIModel{}, diags
	}*/

	var defaultRaw any
	err := json.Unmarshal([]byte(self.Default.ValueString()), &defaultRaw)
	if err != nil {
		diags.AddError(
			"Internal error",
			fmt.Sprintf("Unable to JSON decode default value for property %s, got error: %s",
				self.Title.ValueString(), err))
		return PropertyAPIModel{}, diags
	}

	return PropertyAPIModel{
		Title:            self.Title.ValueString(),
		Description:      self.Description.ValueString(),
		Type:             self.Type.ValueString(),
		Default:          defaultRaw,
		Encrypted:        self.Encrypted.ValueBool(),
		RecreateOnUpdate: self.RecreateOnUpdate.ValueBool(),
		Minimum:          self.Minimum.ValueInt64(),
		Maximum:          self.Maximum.ValueInt64(),
		MinLength:        self.MinLength.ValueInt64(),
		MaxLength:        self.MaxLength.ValueInt64(),
		Pattern:          self.Pattern.ValueString(),
		// FIXME OneOf:
	}, diags
}

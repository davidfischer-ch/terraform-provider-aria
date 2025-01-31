// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PropertyOneOfModel describes the resource data model.
type PropertyOneOfModel struct {
	Const     types.String `tfsdk:"const"`
	Title     types.String `tfsdk:"title"`
	Encrypted types.Bool   `tfsdk:"encrypted"`
}

// PropertyOneOfAPIModel describes the resource API model.
type PropertyOneOfAPIModel struct {
	Const     string `json:"const"`
	Title     string `json:"title"`
	Encrypted bool   `json:"encrypted"`
}

func (self *PropertyOneOfModel) FromAPI(raw PropertyOneOfAPIModel) {
	self.Const = types.StringValue(raw.Const)
	self.Title = types.StringValue(raw.Title)
	self.Encrypted = types.BoolValue(raw.Encrypted)
}

func (self PropertyOneOfModel) ToAPI() PropertyOneOfAPIModel {
	return PropertyOneOfAPIModel{
		Const:     self.Const.ValueString(),
		Title:     self.Title.ValueString(),
		Encrypted: self.Encrypted.ValueBool(),
	}
}

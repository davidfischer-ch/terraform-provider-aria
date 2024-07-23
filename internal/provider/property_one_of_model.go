// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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

func (self *PropertyOneOfModel) FromAPI(
	ctx context.Context,
	raw PropertyOneOfAPIModel,
) diag.Diagnostics {
	self.Const = types.StringValue(raw.Const)
	self.Title = types.StringValue(raw.Title)
	self.Encrypted = types.BoolValue(raw.Encrypted)
	return diag.Diagnostics{}
}

func (self PropertyOneOfModel) ToAPI(
	ctx context.Context,
) (PropertyOneOfAPIModel, diag.Diagnostics) {
	return PropertyOneOfAPIModel{
		Const:     self.Const.ValueString(),
		Title:     self.Title.ValueString(),
		Encrypted: self.Encrypted.ValueBool(),
	}, diag.Diagnostics{}
}

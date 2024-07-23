// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ActionParameterModel describes the resource data model.
type ActionParameterModel struct {
	Type        types.String `tfsdk:"type"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

// ActionParameterAPIModel describes the resource API model.
type ActionParameterAPIModel struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (self *ActionParameterModel) FromAPI(
	ctx context.Context,
	raw ActionParameterAPIModel,
) diag.Diagnostics {
	self.Type = types.StringValue(raw.Type)
	self.Name = types.StringValue(raw.Name)
	self.Description = types.StringValue(raw.Description)
	return diag.Diagnostics{}
}

func (self ActionParameterModel) ToAPI(
	ctx context.Context,
) (ActionParameterAPIModel, diag.Diagnostics) {
	return ActionParameterAPIModel{
		Type:        self.Type.ValueString(),
		Name:        self.Name.ValueString(),
		Description: self.Description.ValueString(),
	}, diag.Diagnostics{}
}

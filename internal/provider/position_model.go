// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PositionModel describes the resource data model.
type PositionModel struct {
	X types.Float64 `tfsdk:"x"`
	Y types.Float64 `tfsdk:"y"`
}

// PositionAPIModel describes the resource API model.
type PositionAPIModel struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (self *PositionModel) FromAPI(
	ctx context.Context,
	raw PositionAPIModel,
) diag.Diagnostics {
	self.X = types.Float64Value(raw.X)
	self.Y = types.Float64Value(raw.Y)
	return diag.Diagnostics{}
}

func (self PositionModel) ToAPI(
	ctx context.Context,
) (PositionAPIModel, diag.Diagnostics) {
	return PositionAPIModel{
		X: self.X.ValueFloat64(),
		Y: self.Y.ValueFloat64(),
	}, diag.Diagnostics{}
}

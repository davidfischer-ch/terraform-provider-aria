// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
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

func (self *PositionModel) FromAPI(raw PositionAPIModel) {
	self.X = types.Float64Value(raw.X)
	self.Y = types.Float64Value(raw.Y)
}

func (self PositionModel) ToAPI() PositionAPIModel {
	return PositionAPIModel{
		X: self.X.ValueFloat64(),
		Y: self.Y.ValueFloat64(),
	}
}

// Utils -------------------------------------------------------------------------------------------

// Used to convert structure to a types.Object.
func (self PositionModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"x": types.Float64Type,
		"y": types.Float64Type,
	}
}

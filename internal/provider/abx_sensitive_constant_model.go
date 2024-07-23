// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ABXSensitiveConstantModel describes the resource data model.
type ABXSensitiveConstantModel struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Value     types.String `tfsdk:"value"`
	Encrypted types.Bool   `tfsdk:"encrypted"`
	OrgId     types.String `tfsdk:"org_id"`
}

// ABXSensitiveConstantAPIModel describes the resource API model.
type ABXSensitiveConstantAPIModel struct {
	Id            string `json:"id,omitempty"`
	Name          string `json:"name"`
	Value         string `json:"value"`
	Encrypted     bool   `json:"encrypted"`
	OrgId         string `json:"orgId"`
	CreatedMillis uint64 `json:"createdMillis"`
}

func (self ABXSensitiveConstantModel) String() string {
	return fmt.Sprintf(
		"ABX Sensitive Constant %s (%s)",
		self.Id.ValueString(),
		self.Name.ValueString())
}

func (self *ABXSensitiveConstantModel) FromAPI(
	ctx context.Context,
	raw ABXSensitiveConstantAPIModel,
) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	// The value is returned with the following '*****' awesome :)
	// self.Value = types.StringValue(raw.Value)
	self.Encrypted = types.BoolValue(raw.Encrypted)
	self.OrgId = types.StringValue(raw.OrgId)
	return diag.Diagnostics{}
}

func (self ABXSensitiveConstantModel) ToAPI() ABXSensitiveConstantAPIModel {
	return ABXSensitiveConstantAPIModel{
		Name:      self.Name.ValueString(),
		Value:     self.Value.ValueString(),
		Encrypted: self.Encrypted.ValueBool(),
	}
}

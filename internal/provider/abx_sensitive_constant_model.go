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

// Return an appropriate key that can be used for naming mutexes.
// Create: Identifier can be used to prevent concurrent creation of ABX constants.
// Read Update Delete: Identifier can be used to prevent concurrent modifications on the instance.
func (self ABXSensitiveConstantModel) LockKey() string {
	return "abx-constant-" + self.Id.ValueString() // Its not a mistake, shared.
}

func (self ABXSensitiveConstantModel) CreatePath() string {
	return "abx/api/resources/action-secrets"
}

func (self ABXSensitiveConstantModel) ReadPath() string {
	return "abx/api/resources/action-secrets/" + self.Id.ValueString()
}

func (self ABXSensitiveConstantModel) UpdatePath() string {
	return self.ReadPath()
}

func (self ABXSensitiveConstantModel) DeletePath() string {
	return self.ReadPath()
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

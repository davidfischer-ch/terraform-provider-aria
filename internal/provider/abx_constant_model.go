// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ABXConstantModel describes the resource data model.
type ABXConstantModel struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Value     types.String `tfsdk:"value"`
	Encrypted types.Bool   `tfsdk:"encrypted"`
	OrgId     types.String `tfsdk:"org_id"`
}

// ABXConstantAPIModel describes the resource API model.
type ABXConstantAPIModel struct {
	Id            string `json:"id,omitempty"`
	Name          string `json:"name"`
	Value         string `json:"value"`
	Encrypted     bool   `json:"encrypted"`
	OrgId         string `json:"orgId"`
	CreatedMillis uint64 `json:"createdMillis"`
}

func (self ABXConstantModel) String() string {
	return fmt.Sprintf(
		"ABX Constant %s (%s)",
		self.Id.ValueString(),
		self.Name.ValueString())
}

// Return an appropriate key that can be used for naming mutexes.
// Create: Identifier can be used to prevent concurrent creation of ABX constants.
// Read Update Delete: Identifier can be used to prevent concurrent modifications on the instance.
func (self ABXConstantModel) LockKey() string {
	return "abx-constant-" + self.Id.ValueString()
}

func (self ABXConstantModel) CreatePath() string {
	return "abx/api/resources/action-secrets"
}

func (self ABXConstantModel) ReadPath() string {
	return "abx/api/resources/action-secrets/" + self.Id.ValueString()
}

func (self ABXConstantModel) UpdatePath() string {
	return self.ReadPath()
}

func (self ABXConstantModel) DeletePath() string {
	return self.ReadPath()
}

func (self *ABXConstantModel) FromAPI(
	ctx context.Context,
	raw ABXConstantAPIModel,
) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Value = types.StringValue(raw.Value)
	self.Encrypted = types.BoolValue(raw.Encrypted)
	self.OrgId = types.StringValue(raw.OrgId)
	return diag.Diagnostics{}
}

func (self ABXConstantModel) ToAPI() ABXConstantAPIModel {
	return ABXConstantAPIModel{
		Name:      self.Name.ValueString(),
		Value:     self.Value.ValueString(),
		Encrypted: self.Encrypted.ValueBool(),
	}
}

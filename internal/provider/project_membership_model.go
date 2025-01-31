// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ProjectMembershipModel describes the resource data model.
type ProjectMembershipModel struct {
	Email types.String `tfsdk:"email"`
	Type  types.String `tfsdk:"type"`
	Role  types.String `tfsdk:"role"`
}

// ProjectMembershipAPIModel describes the resource API model.
type ProjectMembershipAPIModel struct {
	Email string `json:"email"`
	Type  string `json:"type"`
	Role  string `json:"role"`
}

func (self ProjectMembershipModel) String() string {
	return fmt.Sprintf(
		"Project %s Membership %s %s",
		self.Role.ValueString(),
		self.Type.ValueString(),
		self.Email.ValueString())
}

func (self *ProjectMembershipModel) FromAPI(raw ProjectMembershipAPIModel) {
	self.Email = types.StringValue(raw.Email)
	self.Type = types.StringValue(raw.Type)
	self.Role = types.StringValue(raw.Role)
}

func (self ProjectMembershipModel) ToAPI() ProjectMembershipAPIModel {
	return ProjectMembershipAPIModel{
		Email: self.Email.ValueString(),
		Type:  self.Type.ValueString(),
		Role:  self.Role.ValueString(),
	}
}

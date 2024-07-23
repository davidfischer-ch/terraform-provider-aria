// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CustomNamingProjectFilterModel describes the resource data model.
type CustomNamingProjectFilterModel struct {
	Id          types.String `tfsdk:"id"`
	Active      types.Bool   `tfsdk:"active"`
	OrgDefault  types.Bool   `tfsdk:"org_default"`
	OrgId       types.String `tfsdk:"org_id"`
	ProjectId   types.String `tfsdk:"project_id"`
	ProjectName types.String `tfsdk:"project_name"`
}

// CustomNamingProjectFilterAPIModel describes the resource API model.
type CustomNamingProjectFilterAPIModel struct {
	Id          string `json:"id,omitempty"`
	Active      bool   `json:"active"`
	OrgDefault  bool   `json:"defaultOrg"`
	OrgId       string `json:"orgId"`
	ProjectId   string `json:"projectId"`
	ProjectName string `json:"projectName"`
}

func (self CustomNamingProjectFilterModel) String() string {
	return fmt.Sprintf(
		"Custom Naming Project Filter %s (ID='%s', Name='%s')",
		self.Id.ValueString(),
		self.ProjectId.ValueString(),
		self.ProjectName.ValueString())
}

func (self *CustomNamingProjectFilterModel) FromAPI(
	ctx context.Context,
	raw CustomNamingProjectFilterAPIModel,
) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	self.Active = types.BoolValue(raw.Active)
	self.OrgDefault = types.BoolValue(raw.OrgDefault)
	self.OrgId = types.StringValue(raw.OrgId)
	self.ProjectId = types.StringValue(raw.ProjectId)
	self.ProjectName = types.StringValue(raw.ProjectName)
	return diag.Diagnostics{}
}

func (self CustomNamingProjectFilterModel) ToAPI(
	ctx context.Context,
) (CustomNamingProjectFilterAPIModel, diag.Diagnostics) {
	return CustomNamingProjectFilterAPIModel{
		Id:          self.Id.ValueString(),
		Active:      self.Active.ValueBool(),
		OrgDefault:  self.OrgDefault.ValueBool(),
		OrgId:       self.OrgId.ValueString(),
		ProjectId:   self.ProjectId.ValueString(),
		ProjectName: self.ProjectName.ValueString(),
	}, diag.Diagnostics{}
}

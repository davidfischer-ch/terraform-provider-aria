// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PropertyGroupModel describes the resource data model.
type PropertyGroupModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`

	Properties PropertiesModel `tfsdk:"properties"`

	ProjectId types.String `tfsdk:"project_id"`
	OrgId     types.String `tfsdk:"org_id"`
}

// PropertyGroupAPIModel describes the resource API model.
type PropertyGroupAPIModel struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`

	Properties PropertiesAPIModel `json:"properties"`

	ProjectId string `json:"projectId,omitempty"`
	OrgId     string `json:"orgId"`
}

func (self *PropertyGroupModel) String() string {
	return fmt.Sprintf(
		"Property Group %s (%s)",
		self.Id.ValueString(),
		self.Name.ValueString())
}

func (self *PropertyGroupModel) FromAPI(
	ctx context.Context,
	raw PropertyGroupAPIModel,
) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Description = types.StringValue(raw.Description)
	self.Type = types.StringValue(raw.Type)
	diags := self.Properties.FromAPI(ctx, raw.Properties)
	self.ProjectId = types.StringValue(raw.ProjectId)
	self.OrgId = types.StringValue(raw.OrgId)
	return diags
}

func (self *PropertyGroupModel) ToAPI(
	ctx context.Context,
) (PropertyGroupAPIModel, diag.Diagnostics) {
	propertiesRaw, diags := self.Properties.ToAPI(ctx)
	return PropertyGroupAPIModel{
		Name:        self.Name.ValueString(),
		Description: CleanString(self.Description.ValueString()),
		Type:        self.Type.ValueString(),
		Properties:  propertiesRaw,
		ProjectId:   self.ProjectId.ValueString(),
		OrgId:       self.OrgId.ValueString(),
	}, diags
}

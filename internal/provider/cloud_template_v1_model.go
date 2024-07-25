// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CloudTemplateModel describes the resource data model.
type CloudTemplateModel struct {
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	RequestScopeOrg types.Bool   `tfsdk:"request_scope_org"`

	Inputs    UnorderedPropertiesModel    `tfsdk:"properties"`
	Resources CloudTemplateResourcesModel `tfsdk:"resources"`

	ProjectId types.String `tfsdk:"project_id"`
	OrgId     types.String `tfsdk:"org_id"`
}

// CloudTemplateAPIModel describes the resource API model.
type CloudTemplateAPIModel struct {
	Id              string `json:"id,omitempty"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	RequestScopeOrg bool   `json:"requestScopeOrg"`

	Inputs    UnorderedPropertiesAPIModel    `json:"inputs"`
	Resources CloudTemplateResourcesAPIModel `json:"resources"`

	ProjectId string `json:"projectId"`
	OrgId     string `json:"orgId"`
}

func (self CloudTemplateModel) String() string {
	return fmt.Sprintf("Cloud Template v1 %s", self.Name.ValueString())
}

func (self *CloudTemplateModel) FromAPI(
	ctx context.Context,
	raw CloudTemplateAPIModel,
) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Description = types.StringValue(raw.Description)
	self.RequestScopeOrg = types.BoolValue(raw.RequestScopeOrg)
	self.ProjectId = types.StringValue(raw.ProjectId)
	self.OrgId = types.StringValue(raw.OrgId)
	diags := self.Inputs.FromAPI(ctx, raw.Inputs)
	diags.Append(self.Resources.FromAPI(ctx, raw.Resources)...)
	return diags
}

func (self CloudTemplateModel) ToAPI(
	ctx context.Context,
) (CloudTemplateAPIModel, diag.Diagnostics) {

	inputsRaw, diags := self.Inputs.ToAPI(ctx)
	resourcesRaw, resourcesDiags := self.Resources.ToAPI(ctx)
	diags.Append(resourcesDiags...)

	return CloudTemplateAPIModel{
		Id:          self.Id.ValueString(),
		Name:        self.Name.ValueString(),
		Description: CleanString(self.Description.ValueString()),
		ProjectId:   self.ProjectId.ValueString(),
		OrgId:       self.OrgId.ValueString(),
		Inputs:      inputsRaw,
		Resources:   resourcesRaw,
	}, diags
}

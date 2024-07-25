// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CloudTemplateV1Model describes the resource data model.
type CloudTemplateV1Model struct {
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	RequestScopeOrg types.Bool   `tfsdk:"request_scope_org"`

	Inputs    UnorderedPropertiesModel    `tfsdk:"properties"`
	Resources CloudTemplateResourcesModel `tfsdk:"resources"`

	ProjectId types.String `tfsdk:"project_id"`
	OrgId     types.String `tfsdk:"org_id"`
}

// CloudTemplateV1APIModel describes the resource API model.
type CloudTemplateV1APIModel struct {
	Id              string `json:"id,omitempty"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	RequestScopeOrg bool   `json:"requestScopeOrg"`

	Inputs    UnorderedPropertiesAPIModel    `json:"inputs"`
	Resources CloudTemplateResourcesAPIModel `json:"resources"`

	ProjectId string `json:"projectId"`
	OrgId     string `json:"orgId"`
}

func (self CloudTemplateV1Model) String() string {
	return fmt.Sprintf("Cloud Template v1 %s", self.Name.ValueString())
}

func (self *CloudTemplateV1Model) FromAPI(
	ctx context.Context,
	raw CloudTemplateV1APIModel,
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

func (self CloudTemplateV1Model) ToAPI(
	ctx context.Context,
) (CloudTemplateV1APIModel, diag.Diagnostics) {

	inputsRaw, diags := self.Inputs.ToAPI(ctx)
	resourcesRaw, resourcesDiags := self.Resources.ToAPI(ctx)
	diags.Append(resourcesDiags...)

	return CloudTemplateV1APIModel{
		Id:          self.Id.ValueString(),
		Name:        self.Name.ValueString(),
		Description: CleanString(self.Description.ValueString()),
		ProjectId:   self.ProjectId.ValueString(),
		OrgId:       self.OrgId.ValueString(),
		Inputs:      inputsRaw,
		Resources:   resourcesRaw,
	}, diags
}

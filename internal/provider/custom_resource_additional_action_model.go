// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CustomResourceAdditionalActionModel describes the resource data model.
type CustomResourceAdditionalActionModel struct {
	Id           types.String              `tfsdk:"id"`
	Name         types.String              `tfsdk:"name"`
	DisplayName  types.String              `tfsdk:"display_name"`
	Description  types.String              `tfsdk:"description"`
	ProviderName types.String              `tfsdk:"provider_name"`
	ResourceType types.String              `tfsdk:"resource_type"`
	RunnableItem CustomResourceActionModel `tfsdk:"runnable_item"`
	/*FormDefinition CustomFormModel           `tfsdk:"form_definition"`*/

	ProjectId types.String `tfsdk:"project_id"`
	OrgId     types.String `tfsdk:"org_id"`
}

// CustomResourceAdditionalActionAPIModel describes the resource API model.
type CustomResourceAdditionalActionAPIModel struct {
	Id           string                       `json:"id"`
	Name         string                       `json:"name"`
	DisplayName  string                       `json:"displayName"`
	Description  string                       `json:"description"`
	ProviderName string                       `json:"providerName"`
	ResourceType string                       `json:"resourceType"`
	RunnableItem CustomResourceActionAPIModel `json:"runnableItem"`
	/*FormDefinition CustomFormAPIModel           `json:"formDefinition"`*/

	ProjectId string `json:"projectId"`
	OrgId     string `json:"orgId"`
}

func (self *CustomResourceAdditionalActionModel) String() string {
	return fmt.Sprintf(
		"Custom Resource %s Additional Action %s (%s) project %s",
		self.ResourceType.ValueString(),
		self.Id.ValueString(),
		self.Name.ValueString(),
		self.ProjectId.ValueString())
}

func (self *CustomResourceAdditionalActionModel) FromAPI(
	ctx context.Context,
	raw CustomResourceAdditionalActionAPIModel,
) diag.Diagnostics {

	diags := diag.Diagnostics{}

	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.DisplayName = types.StringValue(raw.DisplayName)
	self.Description = types.StringValue(raw.Description)
	self.ProviderName = types.StringValue(raw.ProviderName)
	self.ResourceType = types.StringValue(raw.ResourceType)
	self.ProjectId = types.StringValue(raw.ProjectId)
	self.OrgId = types.StringValue(raw.OrgId)

	self.RunnableItem = CustomResourceActionModel{}
	diags.Append(self.RunnableItem.FromAPI(ctx, raw.RunnableItem)...)

	// FIXME self.FormDefinition =

	return diags
}

func (self *CustomResourceAdditionalActionModel) ToAPI(
	ctx context.Context,
) (CustomResourceAdditionalActionAPIModel, diag.Diagnostics) {

	diags := diag.Diagnostics{}

	runnableItemRaw, runnableItemDiags := self.RunnableItem.ToAPI(ctx)
	diags.Append(runnableItemDiags...)

	raw := CustomResourceAdditionalActionAPIModel{
		Name:         self.Name.ValueString(),
		DisplayName:  self.DisplayName.ValueString(),
		Description:  self.Description.ValueString(),
		ProviderName: self.ProviderName.ValueString(),
		ResourceType: self.ResourceType.ValueString(),
		RunnableItem: runnableItemRaw,
		// FIXME FormDefinition:
		ProjectId: self.ProjectId.ValueString(),
		OrgId:     self.OrgId.ValueString(),
	}

	// When updating resource
	if !self.Id.IsNull() {
		raw.Id = self.Id.ValueString()
	}

	return raw, diags
}

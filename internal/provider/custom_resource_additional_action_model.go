// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CustomResourceAdditionalActionModel describes the resource data model.
type CustomResourceAdditionalActionModel struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	DisplayName  types.String `tfsdk:"display_name"`
	ProviderName types.String `tfsdk:"provider_name"`
	ResourceType types.String `tfsdk:"resource_type"`
	Status       types.String `tfsdk:"status"`

	RunnableItem   CustomResourceActionModel `tfsdk:"runnable_item"`
	FormDefinition CustomFormModel           `tfsdk:"form_definition"`

	OrgId types.String `tfsdk:"org_id"`
}

// CustomResourceAdditionalActionAPIModel describes the resource API model.
type CustomResourceAdditionalActionAPIModel struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	DisplayName  string `json:"displayName"`
	ProviderName string `json:"providerName"`
	ResourceType string `json:"resourceType"`
	Status       string `json:"status"`

	RunnableItem   CustomResourceActionAPIModel `json:"runnableItem"`
	FormDefinition CustomFormAPIModel           `json:"formDefinition"`

	OrgId string `json:"orgId"`
}

func (self *CustomResourceAdditionalActionModel) FromAPI(
	ctx context.Context,
	raw CustomResourceAdditionalActionAPIModel,
) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.DisplayName = types.StringValue(raw.DisplayName)
	self.ProviderName = types.StringValue(raw.ProviderName)
	self.ResourceType = types.StringValue(raw.ResourceType)
	self.Status = types.StringValue(raw.Status)
	// FIXME self.RunnableItem =
	// FIXME self.FormDefinition =
	self.OrgId = types.StringValue(raw.OrgId)
	return diag.Diagnostics{}
}

func (self *CustomResourceAdditionalActionModel) ToAPI(
	ctx context.Context,
) (CustomResourceAdditionalActionAPIModel, diag.Diagnostics) {
	return CustomResourceAdditionalActionAPIModel{
		Id:           self.Id.ValueString(),
		Name:         self.Name.ValueString(),
		DisplayName:  self.DisplayName.ValueString(),
		ProviderName: self.ProviderName.ValueString(),
		ResourceType: self.ResourceType.ValueString(),
		Status:       self.Status.ValueString(),
		// FIXME RunnableItem:
		// FIXME FormDefinition:
		OrgId: self.OrgId.ValueString(),
	}, diag.Diagnostics{}
}

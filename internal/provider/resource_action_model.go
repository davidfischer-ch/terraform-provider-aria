// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResourceActionModel describes the resource data model.
type ResourceActionModel struct {
	Id           types.String                `tfsdk:"id"`
	Name         types.String                `tfsdk:"name"`
	DisplayName  types.String                `tfsdk:"display_name"`
	Description  types.String                `tfsdk:"description"`
	ProviderName types.String                `tfsdk:"provider_name"`
	ResourceType types.String                `tfsdk:"resource_type"`
	RunnableItem ResourceActionRunnableModel `tfsdk:"runnable_item"`
	Status       types.String                `tfsdk:"status"`

	// Of type CustomFormModel
	FormDefinition types.Object `tfsdk:"form_definition"`

	ProjectId types.String `tfsdk:"project_id"`
	OrgId     types.String `tfsdk:"org_id"`
}

// ResourceActionAPIModel describes the resource API model.
type ResourceActionAPIModel struct {
	Id           string                         `json:"id,omitempty"`
	Name         string                         `json:"name"`
	DisplayName  string                         `json:"displayName"`
	Description  string                         `json:"description"`
	ProviderName string                         `json:"providerName"`
	ResourceType string                         `json:"resourceType"`
	RunnableItem ResourceActionRunnableAPIModel `json:"runnableItem"`
	Status       string                         `json:"status"`

	FormDefinition *CustomFormAPIModel `json:"formDefinition,omitempty"`

	ProjectId string `json:"projectId,omitempty"`
	OrgId     string `json:"orgId,omitempty"`
}

func (self *ResourceActionModel) String() string {
	return fmt.Sprintf(
		"Resource %s Action %s (%s) project %s",
		self.ResourceType.ValueString(),
		self.Id.ValueString(),
		self.Name.ValueString(),
		self.ProjectId.ValueString())
}

func (self *ResourceActionModel) FromAPI(
	ctx context.Context,
	raw ResourceActionAPIModel,
) diag.Diagnostics {

	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.DisplayName = types.StringValue(raw.DisplayName)
	self.Description = types.StringValue(raw.Description)
	self.ProviderName = types.StringValue(raw.ProviderName)
	self.ResourceType = types.StringValue(raw.ResourceType)
	self.Status = types.StringValue(raw.Status)
	self.ProjectId = types.StringValue(raw.ProjectId)
	self.OrgId = types.StringValue(raw.OrgId)

	self.RunnableItem = ResourceActionRunnableModel{}
	diags := self.RunnableItem.FromAPI(ctx, raw.RunnableItem)

	var formDiags diag.Diagnostics
	self.FormDefinition, formDiags = raw.FormDefinition.ToObject(ctx)
	diags.Append(formDiags...)

	return diags
}

func (self *ResourceActionModel) ToAPI(
	ctx context.Context,
) (ResourceActionAPIModel, diag.Diagnostics) {

	formDefinitionRaw, diags := CustomFormAPIModelFromObject(ctx, self.FormDefinition)
	runnableItemRaw, runnableItemDiags := self.RunnableItem.ToAPI(ctx)
	diags.Append(runnableItemDiags...)

	return ResourceActionAPIModel{
		Id:             self.Id.ValueString(),
		Name:           self.Name.ValueString(),
		DisplayName:    self.DisplayName.ValueString(),
		Description:    self.Description.ValueString(),
		ProviderName:   self.ProviderName.ValueString(),
		ResourceType:   self.ResourceType.ValueString(),
		RunnableItem:   runnableItemRaw,
		FormDefinition: formDefinitionRaw,
		Status:         self.Status.ValueString(),
		ProjectId:      self.ProjectId.ValueString(),
		OrgId:          self.OrgId.ValueString(),
	}, diags
}

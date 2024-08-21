// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
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
	ResourceId   types.String                `tfsdk:"resource_id"`
	ResourceType types.String                `tfsdk:"resource_type"`
	RunnableItem ResourceActionRunnableModel `tfsdk:"runnable_item"`
	Criteria     jsontypes.Normalized        `tfsdk:"criteria"`
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
	Criteria     map[string]interface{}         `json:"criteria,omitempty"`
	Status       string                         `json:"status"`

	FormDefinition *CustomFormAPIModel `json:"formDefinition,omitempty"`

	ProjectId string `json:"projectId,omitempty"`
	OrgId     string `json:"orgId,omitempty"`
}

func (self ResourceActionModel) String() string {
	return fmt.Sprintf(
		"Resource %s Action %s (%s) project %s",
		self.ResourceType.ValueString(),
		self.Id.ValueString(),
		self.Name.ValueString(),
		self.ProjectId.ValueString())
}

// Return an appropriate key that can be used for naming mutexes.
// Create: Identifier can be used to prevent concurrent creation of resource actions.
// Read Update Delete: Identifier can be used to prevent concurrent modifications on the instance.
func (self ResourceActionModel) LockKey() string {
	if self.ForCustom() {
		// The custom resource is the object that will be manipulated to manage the action!
		return CustomResourceModel{Id: self.ResourceId}.LockKey()
	}
	return "resource-actions-" + self.Id.ValueString()
}

func (self ResourceActionModel) CreatePath() string {
	if len(self.ResourceId.ValueString()) > 0 {
		// Custom resource ...
		panic("not implemented")
	}
	// Native resource ...
	return "form-service/api/custom/resource-actions"
}

func (self ResourceActionModel) ReadPath() string {
	if len(self.ResourceId.ValueString()) > 0 {
		return fmt.Sprintf(
			// Custom Resource ...
			"form-service/api/custom/resource-types/%s/resource-actions/%s",
			self.ResourceId.ValueString(), self.Id.ValueString())
	}
	// Native resource ...
	return "form-service/api/custom/resource-actions/" + self.Id.ValueString()
}

func (self ResourceActionModel) UpdatePath() string {
	if len(self.ResourceId.ValueString()) > 0 {
		// Custom resource ...
		panic("not implemented")
	}
	// Native resource ...
	return "form-service/api/custom/resource-actions"
}

func (self ResourceActionModel) DeletePath() string {
	if len(self.ResourceId.ValueString()) > 0 {
		// Custom resource ...
		panic("not implemented")
	}
	// Native resource ...
	return "form-service/api/custom/resource-actions/" + self.Id.ValueString()
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

	// Criteria API data -> JSON Encoded
	if raw.Criteria == nil {
		self.Criteria = jsontypes.NewNormalizedNull()
	} else {
		criteriaJSON, err := json.Marshal(raw.Criteria)
		if err != nil {
			diags.AddError(
				"Client error",
				fmt.Sprintf("Unable to JSON encode %s criteria, got error: %s", self.String(), err))
		} else {
			self.Criteria = jsontypes.NewNormalizedValue(string(criteriaJSON))
		}
	}

	var formDiags diag.Diagnostics
	self.FormDefinition, formDiags = raw.FormDefinition.ToObject(ctx)
	diags.Append(formDiags...)

	return diags
}

func (self ResourceActionModel) ToAPI(
	ctx context.Context,
) (ResourceActionAPIModel, diag.Diagnostics) {

	formDefinitionRaw, diags := CustomFormAPIModelFromObject(ctx, self.FormDefinition)
	runnableItemRaw, runnableItemDiags := self.RunnableItem.ToAPI(ctx)
	diags.Append(runnableItemDiags...)

	// Defining name is mandatory when passing other form attributes such as styles...
	// Otherwise we trigger an Hibernate exception on Aria. Fortunately We known that
	// form's name is equal to resource action's name, so we set it like this.
	if formDefinitionRaw != nil {
		formDefinitionRaw.Name = self.Name.ValueString()
	}

	// Criteria JSON Encoded -> API data
	var criteriaRaw map[string]interface{}
	if self.Criteria.IsNull() {
		criteriaRaw = nil
	} else {
		diags.Append(self.Criteria.Unmarshal(&criteriaRaw)...)
	}

	return ResourceActionAPIModel{
		Id:             self.Id.ValueString(),
		Name:           self.Name.ValueString(),
		DisplayName:    self.DisplayName.ValueString(),
		Description:    self.Description.ValueString(),
		ProviderName:   self.ProviderName.ValueString(),
		ResourceType:   self.ResourceType.ValueString(),
		RunnableItem:   runnableItemRaw,
		Criteria:       criteriaRaw,
		FormDefinition: formDefinitionRaw,
		Status:         self.Status.ValueString(),
		ProjectId:      self.ProjectId.ValueString(),
		OrgId:          self.OrgId.ValueString(),
	}, diags
}

func (self ResourceActionModel) ForCustom() bool {
	return len(self.ResourceId.ValueString()) > 0
}

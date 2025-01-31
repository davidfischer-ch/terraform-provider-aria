// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CustomResourceModel describes the resource data model.
type CustomResourceModel struct {
	Id           types.String `tfsdk:"id"`
	DisplayName  types.String `tfsdk:"display_name"`
	Description  types.String `tfsdk:"description"`
	ResourceType types.String `tfsdk:"resource_type"`
	SchemaType   types.String `tfsdk:"schema_type"`
	Status       types.String `tfsdk:"status"`

	Properties UnorderedPropertiesModel `tfsdk:"properties"`

	Create ResourceActionRunnableModel `tfsdk:"create"`
	Read   ResourceActionRunnableModel `tfsdk:"read"`
	Update ResourceActionRunnableModel `tfsdk:"update"`
	Delete ResourceActionRunnableModel `tfsdk:"delete"`

	AdditionalActions []ResourceActionModel `tfsdk:"-"`

	ProjectId types.String `tfsdk:"project_id"`
	OrgId     types.String `tfsdk:"org_id"`
}

// CustomResourceAPIModel describes the resource API model.
type CustomResourceAPIModel struct {
	Id           string `json:"id,omitempty"`
	DisplayName  string `json:"displayName"`
	Description  string `json:"description"`
	ResourceType string `json:"resourceType"`
	SchemaType   string `json:"schemaType"`
	Status       string `json:"status"`

	Properties CustomResourcePropertiesAPIModel `json:"properties"`

	MainActions       map[string]ResourceActionRunnableAPIModel `json:"mainActions"`
	AdditionalActions []ResourceActionAPIModel                  `json:"additionalActions"`

	// Omit empty Project ID to prevent "projectId cannot be updated for type (...) on UPDATE"
	ProjectId string `json:"projectId,omitempty"`
	OrgId     string `json:"orgId,omitempty"`
}

func (self CustomResourceModel) String() string {
	return fmt.Sprintf(
		"Custom Resource %s (%s)",
		self.Id.ValueString(),
		self.DisplayName.ValueString())
}

// Return an appropriate key that can be used for naming mutexes.
// Create: Identifier can be used to prevent concurrent creation of custom resources.
// Read Update Delete: Identifier can be used to prevent concurrent modifications on the instance.
func (self CustomResourceModel) LockKey() string {
	return "custom-resource-" + self.Id.ValueString()
}

func (self CustomResourceModel) CreatePath() string {
	return "form-service/api/custom/resource-types"
}

func (self CustomResourceModel) ReadPath() string {
	return "form-service/api/custom/resource-types/" + self.Id.ValueString()
}

func (self CustomResourceModel) UpdatePath() string {
	return self.CreatePath() // Its not a mistake...
}

func (self CustomResourceModel) DeletePath() string {
	return self.ReadPath()
}

func (self *CustomResourceModel) FromAPI(
	ctx context.Context,
	raw CustomResourceAPIModel,
) diag.Diagnostics {

	self.Id = types.StringValue(raw.Id)
	self.DisplayName = types.StringValue(raw.DisplayName)
	self.Description = types.StringValue(raw.Description)
	self.ResourceType = types.StringValue(raw.ResourceType)
	self.SchemaType = types.StringValue(raw.SchemaType)
	self.Status = types.StringValue(raw.Status)
	self.ProjectId = types.StringValue(raw.ProjectId)
	self.OrgId = types.StringValue(raw.OrgId)
	diags := self.Properties.FromAPI(ctx, raw.Properties.Properties)

	self.Create = ResourceActionRunnableModel{}
	self.Create.FromAPI(raw.MainActions["create"])

	self.Read = ResourceActionRunnableModel{}
	self.Read.FromAPI(raw.MainActions["read"])

	self.Update = ResourceActionRunnableModel{}
	self.Update.FromAPI(raw.MainActions["update"])

	self.Delete = ResourceActionRunnableModel{}
	self.Delete.FromAPI(raw.MainActions["delete"])

	self.AdditionalActions = []ResourceActionModel{}
	for _, actionRaw := range raw.AdditionalActions {
		action := ResourceActionModel{}
		diags.Append(action.FromAPI(ctx, actionRaw)...)
		self.AdditionalActions = append(self.AdditionalActions, action)
	}

	return diags
}

func (self CustomResourceModel) ToAPI(
	ctx context.Context,
) (CustomResourceAPIModel, diag.Diagnostics) {

	propertiesRaw, diags := self.Properties.ToAPI(ctx)

	createRaw := self.Create.ToAPI()
	readRaw := self.Read.ToAPI()
	updateRaw := self.Update.ToAPI()
	deleteRaw := self.Delete.ToAPI()

	var additionalActionsRaw []ResourceActionAPIModel
	for _, action := range self.AdditionalActions {
		actionRaw, actionDiags := action.ToAPI(ctx)
		additionalActionsRaw = append(additionalActionsRaw, actionRaw)
		diags.Append(actionDiags...)
	}

	return CustomResourceAPIModel{
		Id:           self.Id.ValueString(),
		DisplayName:  self.DisplayName.ValueString(),
		Description:  CleanString(self.Description.ValueString()),
		ResourceType: self.ResourceType.ValueString(),
		SchemaType:   self.SchemaType.ValueString(),
		Status:       self.Status.ValueString(),
		ProjectId:    self.ProjectId.ValueString(),
		OrgId:        self.OrgId.ValueString(),
		Properties: CustomResourcePropertiesAPIModel{
			Properties: propertiesRaw,
		},
		MainActions: map[string]ResourceActionRunnableAPIModel{
			"create": createRaw,
			"read":   readRaw,
			"update": updateRaw,
			"delete": deleteRaw,
		},
		AdditionalActions: additionalActionsRaw,
	}, diags
}

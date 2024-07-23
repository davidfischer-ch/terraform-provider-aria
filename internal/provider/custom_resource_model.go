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

	MainActions map[string]ResourceActionRunnableAPIModel `json:"mainActions"`

	ProjectId string `json:"projectId"`
	OrgId     string `json:"orgId"`
}

func (self *CustomResourceModel) String() string {
	return fmt.Sprintf(
		"ABX Custom Resource %s (%s)",
		self.Id.ValueString(),
		self.DisplayName.ValueString())
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
	diags.Append(self.Create.FromAPI(ctx, raw.MainActions["create"])...)

	self.Read = ResourceActionRunnableModel{}
	diags.Append(self.Read.FromAPI(ctx, raw.MainActions["read"])...)

	self.Update = ResourceActionRunnableModel{}
	diags.Append(self.Update.FromAPI(ctx, raw.MainActions["update"])...)

	self.Delete = ResourceActionRunnableModel{}
	diags.Append(self.Delete.FromAPI(ctx, raw.MainActions["delete"])...)

	return diags
}

func (self *CustomResourceModel) ToAPI(
	ctx context.Context,
) (CustomResourceAPIModel, diag.Diagnostics) {

	propertiesRaw, diags := self.Properties.ToAPI(ctx)

	createRaw, createDiags := self.Create.ToAPI(ctx)
	diags.Append(createDiags...)

	readRaw, readDiags := self.Read.ToAPI(ctx)
	diags.Append(readDiags...)

	updateRaw, updateDiags := self.Update.ToAPI(ctx)
	diags.Append(updateDiags...)

	deleteRaw, deleteDiags := self.Delete.ToAPI(ctx)
	diags.Append(deleteDiags...)

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
	}, diags
}

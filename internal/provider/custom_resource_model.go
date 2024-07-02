// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

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

	Properties []PropertyModel `tfsdk:"properties"`

	Create CustomResourceActionModel `tfsdk:"create"`
	Read   CustomResourceActionModel `tfsdk:"read"`
	Update CustomResourceActionModel `tfsdk:"update"`
	Delete CustomResourceActionModel `tfsdk:"delete"`

	// Set of CustomResourceAdditionalActionModel
	/* AdditionalActions types.Set `tfsdk:"additional_actions"` */

	ProjectId types.String `tfsdk:"project_id"`
	OrgId     types.String `tfsdk:"org_id"`
}

// CustomResourcePropertiesAPIModel describes the resource API model.
type CustomResourcePropertiesAPIModel struct {
	Properties map[string]PropertyAPIModel `json:"properties"`
}

// CustomResourceAPIModel describes the resource API model.
type CustomResourceAPIModel struct {
	Id           string `json:"id"`
	DisplayName  string `json:"displayName"`
	Description  string `json:"description"`
	ResourceType string `json:"resourceType"`
	SchemaType   string `json:"schemaType"`
	Status       string `json:"status"`

	Properties CustomResourcePropertiesAPIModel `json:"properties"`

	MainActions map[string]CustomResourceActionAPIModel `json:"mainActions"`

	/* AdditionalActions []CustomResourceAdditionalActionAPIModel `json:"additionalActions"` */

	ProjectId string `json:"projectId"`
	OrgId     string `json:"orgId"`
}

// https://stackoverflow.com/questions/47339542/defining-custom-unmarshalling-for-non-built-in-types
/*func (self *CustomResourceAPIModel) UnmarshalJSON(bytes []byte) error {
	panic(errors.New(string(bytes)))
}*/

func (self *CustomResourceModel) FromAPI(
	ctx context.Context,
	raw CustomResourceAPIModel,
) diag.Diagnostics {

	diags := diag.Diagnostics{}

	self.Id = types.StringValue(raw.Id)
	self.DisplayName = types.StringValue(raw.DisplayName)
	self.Description = types.StringValue(raw.Description)
	self.ResourceType = types.StringValue(raw.ResourceType)
	self.SchemaType = types.StringValue(raw.SchemaType)
	self.Status = types.StringValue(raw.Status)
	self.ProjectId = types.StringValue(raw.ProjectId)
	self.OrgId = types.StringValue(raw.OrgId)

	self.Properties = []PropertyModel{}

	self.Create = CustomResourceActionModel{}
	diags.Append(self.Create.FromAPI(ctx, raw.MainActions["create"])...)

	self.Read = CustomResourceActionModel{}
	diags.Append(self.Read.FromAPI(ctx, raw.MainActions["read"])...)

	self.Update = CustomResourceActionModel{}
	diags.Append(self.Update.FromAPI(ctx, raw.MainActions["update"])...)

	self.Delete = CustomResourceActionModel{}
	diags.Append(self.Delete.FromAPI(ctx, raw.MainActions["delete"])...)

	/* self.AdditionalActions */

	return diags
}

func (self *CustomResourceModel) ToAPI(
	ctx context.Context,
) (CustomResourceAPIModel, diag.Diagnostics) {

	diags := diag.Diagnostics{}

	// FIXME TODO THIS
	properties := map[string]PropertyAPIModel{}

	createRaw, createDiags := self.Create.ToAPI(ctx)
	diags.Append(createDiags...)

	readRaw, readDiags := self.Read.ToAPI(ctx)
	diags.Append(readDiags...)

	updateRaw, updateDiags := self.Update.ToAPI(ctx)
	diags.Append(updateDiags...)

	deleteRaw, deleteDiags := self.Delete.ToAPI(ctx)
	diags.Append(deleteDiags...)

	raw := CustomResourceAPIModel{
		DisplayName:  self.DisplayName.ValueString(),
		Description:  self.Description.ValueString(),
		ResourceType: self.ResourceType.ValueString(),
		SchemaType:   self.SchemaType.ValueString(),
		Status:       self.Status.ValueString(),
		ProjectId:    self.ProjectId.ValueString(),
		OrgId:        self.OrgId.ValueString(),
		Properties: CustomResourcePropertiesAPIModel{
			Properties: properties,
		},
		MainActions: map[string]CustomResourceActionAPIModel{
			"create": createRaw,
			"read":   readRaw,
			"update": updateRaw,
			"delete": deleteRaw,
		},
		/* AdditionalActions */
	}

	// When updating resource
	if !self.Id.IsNull() {
		raw.Id = self.Id.ValueString()
	}

	return raw, diags
}

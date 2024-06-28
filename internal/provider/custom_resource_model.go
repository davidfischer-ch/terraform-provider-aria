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

	// List of PropertyModel
	Properties types.List `tfsdk:"properties"`

	// Of type CustomResourceActionModel
	Create types.Object `tfsdk:"create"`
	Read   types.Object `tfsdk:"read"`
	Update types.Object `tfsdk:"update"`
	Delete types.Object `tfsdk:"delete"`

	// Set of CustomResourceAdditionalActionModel
	AdditionalActions types.Set `tfsdk:"additional_actions"`

	ProjectId types.String `tfsdk:"project_id"`
	OrgId     types.String `tfsdk:"org_id"`
}

// CustomResourceAPIModel describes the resource API model.
type CustomResourceAPIModel struct {
	Id           string `json:"id"`
	DisplayName  string `json:"displayName"`
	Description  string `json:"description"`
	ResourceType string `json:"resourceType"`
	SchemaType   string `json:"schemaType"`
	Status       string `json:"status"`

	Properties map[string]map[string]PropertyAPIModel `json:"properties"`

	MainActions       map[string]CustomResourceActionAPIModel  `json:"mainActions"`
	AdditionalActions []CustomResourceAdditionalActionAPIModel `json:"additionalActions"`

	ProjectId string `json:"projectId"`
	OrgId     string `json:"orgId"`
}

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

	// FIXME self.Properties
	// FIXME self.Create
	// FIXME self.Read
	// FIXME self.Update
	// FIXME self.Delete
	// FIXME self.AdditionalActions

	return diags
}

func (self *CustomResourceModel) ToAPI(
	ctx context.Context,
) (CustomResourceAPIModel, diag.Diagnostics) {

	diags := diag.Diagnostics{}

	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/list
	if self.Properties.IsNull() || self.Properties.IsUnknown() {
		diags.AddError(
			"Configuration error",
			fmt.Sprintf(
				"Unable to manage custom resource %s, properties is either null or unknown",
				self.Id.ValueString()))
		return CustomResourceAPIModel{}, diags
	}

	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/list
	if self.Create.IsNull() || self.Create.IsUnknown() {
		diags.AddError(
			"Configuration error",
			fmt.Sprintf(
				"Unable to manage custom resource %s, create is either null or unknown",
				self.Id.ValueString()))
		return CustomResourceAPIModel{}, diags
	}

	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/list
	if self.Read.IsNull() || self.Read.IsUnknown() {
		diags.AddError(
			"Configuration error",
			fmt.Sprintf(
				"Unable to manage custom resource %s, read is either null or unknown",
				self.Id.ValueString()))
		return CustomResourceAPIModel{}, diags
	}

	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/list
	if self.Update.IsNull() || self.Update.IsUnknown() {
		diags.AddError(
			"Configuration error",
			fmt.Sprintf(
				"Unable to manage custom resource %s, update is either null or unknown",
				self.Id.ValueString()))
		return CustomResourceAPIModel{}, diags
	}

	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/list
	if self.Delete.IsNull() || self.Delete.IsUnknown() {
		diags.AddError(
			"Configuration error",
			fmt.Sprintf(
				"Unable to manage custom resource %s, delete is either null or unknown",
				self.Id.ValueString()))
		return CustomResourceAPIModel{}, diags
	}

	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/list
	if self.AdditionalActions.IsNull() || self.AdditionalActions.IsUnknown() {
		diags.AddError(
			"Configuration error",
			fmt.Sprintf(
				"Unable to manage custom resource %s, additional_actions is either null or unknown",
				self.Id.ValueString()))
		return CustomResourceAPIModel{}, diags
	}

	raw := CustomResourceAPIModel{
		DisplayName:  self.DisplayName.ValueString(),
		Description:  self.Description.ValueString(),
		ResourceType: self.ResourceType.ValueString(),
		SchemaType:   self.SchemaType.ValueString(),
		Status:       self.Status.ValueString(),
		// FIXME Properties
		// FIXME MainActions
		// FIXME AdditionalActions
		ProjectId: self.ProjectId.ValueString(),
		// Let platform manage this field
	}

	// When updating resource
	if !self.Id.IsNull() {
		raw.Id = self.Id.ValueString()
	}

	return raw, diags
}

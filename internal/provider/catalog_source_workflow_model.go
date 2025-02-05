// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// CatalogSourceWorkflowModel describes the resource data model.
type CatalogSourceWorkflowModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Version     types.String `tfsdk:"version"`

	// Of type IntegrationModel
	Integration types.Object `tfsdk:"integration"`
}

// CatalogSourceWorkflowAPIModel describes the resource API model.
type CatalogSourceWorkflowAPIModel struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`

	Integration IntegrationAPIModel `json:"integration,omitempty"`
}

func (self *CatalogSourceWorkflowModel) String() string {
	return fmt.Sprintf(
		"Catalog Source Workflow %s (%s)",
		self.Id.ValueString(),
		self.Name.ValueString())
}

func (self *CatalogSourceWorkflowModel) FromAPI(
	ctx context.Context,
	raw CatalogSourceWorkflowAPIModel,
) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Description = types.StringValue(raw.Description)
	self.Version = types.StringValue(raw.Version)

	diags := diag.Diagnostics{}
	var someDiags diag.Diagnostics

	// Convert integration from raw and then to object
	integration := IntegrationModel{}
	integration.FromAPI(raw.Integration)
	self.Integration, someDiags = types.ObjectValueFrom(
		ctx, integration.AttributeTypes(), integration,
	)
	diags.Append(someDiags...)

	return diags
}

func (self CatalogSourceWorkflowModel) ToAPI(
	ctx context.Context,
) (CatalogSourceWorkflowAPIModel, diag.Diagnostics) {

	diags := diag.Diagnostics{}
	integrationRaw := IntegrationAPIModel{}

	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/object
	if self.Integration.IsNull() || self.Integration.IsUnknown() {
		diags.AddError(
			"Configuration error",
			fmt.Sprintf(
				"Unable to manage %s, integration is either null or unknown",
				self.String()))
	} else {
		// Convert integration from object to raw
		integration := IntegrationModel{}
		diags.Append(self.Integration.As(ctx, &integration, basetypes.ObjectAsOptions{})...)
		integrationRaw = integration.ToAPI()
	}

	return CatalogSourceWorkflowAPIModel{
		Id:          self.Id.ValueString(),
		Name:        self.Name.ValueString(),
		Description: self.Description.ValueString(),
		Version:     self.Version.ValueString(),
		Integration: integrationRaw,
	}, diags
}

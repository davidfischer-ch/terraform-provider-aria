// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// OrchestratorConfigurationModel describes the resource data model.
type OrchestratorConfigurationModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	CategoryId  types.String `tfsdk:"category_id"`
	Version     types.String `tfsdk:"version"`
	VersionId   types.String `tfsdk:"version_id"`

	Attributes types.List `tfsdk:"attributes"`

	ForceDelete types.Bool `tfsdk:"force_delete"`
}

// OrchestratorConfigurationAPIModel describes the resource API model.
type OrchestratorConfigurationAPIModel struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CategoryId  string `json:"category-id"`
	Version     string `json:"version"`

	Attributes []OrchestratorConfigurationAttributeAPIModel `json:"attributes"`
}

func (self OrchestratorConfigurationModel) String() string {
	return fmt.Sprintf(
		"Orchestrator Configuration %s (%s)",
		self.Id.ValueString(),
		self.Name.ValueString())
}

// Return an appropriate key that can be used for naming mutexes.
// Create: Identifier can be used to prevent concurrent creation of vRO configurations.
// Read Update Delete: Identifier can be used to prevent concurrent modifications on the instance.
func (self OrchestratorConfigurationModel) LockKey() string {
	return "orchestrator-configuration-" + self.Id.ValueString()
}

func (self OrchestratorConfigurationModel) CreatePath() string {
	return "vco/api/configurations"
}

func (self OrchestratorConfigurationModel) ReadPath() string {
	return fmt.Sprintf("vco/api/configurations/%s", self.Id.ValueString())
}

func (self OrchestratorConfigurationModel) UpdatePath() string {
	return self.ReadPath()
}

func (self OrchestratorConfigurationModel) DeletePath() string {
	path := self.ReadPath()
	if self.ForceDelete.ValueBool() {
		return path + "?force=true"
	}
	return path
}

func (self *OrchestratorConfigurationModel) FromAPI(
	ctx context.Context,
	raw OrchestratorConfigurationAPIModel,
	response *resty.Response,
) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Description = types.StringValue(raw.Description)
	self.CategoryId = types.StringValue(raw.CategoryId)
	self.Version = types.StringValue(raw.Version)
	self.VersionId = types.StringValue(response.Header().Get("x-vro-changeset-sha"))

	diags := diag.Diagnostics{}

	// Convert attributes from raw
	attributes := []OrchestratorConfigurationAttributeModel{}
	for _, attributeRaw := range raw.Attributes {
		attribute := OrchestratorConfigurationAttributeModel{}
		diags.Append(attribute.FromAPI(ctx, attributeRaw)...)
		attributes = append(attributes, attribute)
	}

	// Store attributes to list value
	var someDiags diag.Diagnostics
	self.Attributes, someDiags = types.ListValueFrom(
		ctx, self.Attributes.ElementType(ctx), attributes,
	)
	diags.Append(someDiags...)

	return diags
}

func (self OrchestratorConfigurationModel) ToAPI(
	ctx context.Context,
) (OrchestratorConfigurationAPIModel, diag.Diagnostics) {

	diags := diag.Diagnostics{}
	attributesRaw := []OrchestratorConfigurationAttributeAPIModel{}

	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/list
	if self.Attributes.IsNull() || self.Attributes.IsUnknown() {
		diags.AddError(
			"Configuration error",
			fmt.Sprintf("Unable to manage %s, attributes is either null or unknown", self.String()))
	} else {
		// Extract attributes from list value and then convert to raw
		attributes := make(
			[]OrchestratorConfigurationAttributeModel, 0, len(self.Attributes.Elements()),
		)
		diags.Append(self.Attributes.ElementsAs(ctx, &attributes, false)...)
		if !diags.HasError() {
			for _, attribute := range attributes {
				attributeRaw, someDiags := attribute.ToAPI(ctx)
				attributesRaw = append(attributesRaw, attributeRaw)
				diags.Append(someDiags...)
			}
		}
	}

	return OrchestratorConfigurationAPIModel{
		Id:          self.Id.ValueString(),
		Name:        self.Name.ValueString(),
		CategoryId:  self.CategoryId.ValueString(),
		Description: self.Description.ValueString(),
		Version:     self.Version.ValueString(),
		Attributes:  attributesRaw,
	}, diag.Diagnostics{}
}

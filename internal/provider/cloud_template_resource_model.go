// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CustomResourcPropertyeModel describes the resource data model.
type CloudTemplateResourceModel struct {
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`

	//Metadata   CloudTemplateResourceMetadataModel   `json:"metadata"`
	//Properties CloudTemplateResourcePropertiesModel `json:"properties"`

	AllocatePerInstance types.Bool `tfsdk:"allocate_per_instance"`
}

// CloudTemplateResourceAPIModel describes the resource API model.
type CloudTemplateResourceAPIModel struct {
	Type string `json:"type"`

	//Metadata   CloudTemplateResourceMetadataModel   `json:"metadata"`
	//Properties CloudTemplateResourcePropertiesModel `json:"properties"`

	AllocatePerInstance *bool `tfsdk:"allocate_per_instance,omitempty"`
}

func (self CloudTemplateResourceModel) String() string {
	return fmt.Sprintf("Cloud Template Resource %s", self.Name.ValueString())
}

func (self *CloudTemplateResourceModel) FromAPI(
	ctx context.Context,
	name string,
	raw CloudTemplateResourceAPIModel,
) diag.Diagnostics {

	diags := diag.Diagnostics{}

	self.Name = types.StringValue(name)
	self.Type = types.StringValue(raw.Type)

	// self.Metadata =
	// self.Properties =

	self.AllocatePerInstance = types.BoolPointerValue(raw.AllocatePerInstance)

	return diags
}

func (self CloudTemplateResourceModel) ToAPI(
	ctx context.Context,
) (string, CloudTemplateResourceAPIModel, diag.Diagnostics) {

	diags := diag.Diagnostics{}

	return self.Name.ValueString(),
		CloudTemplateResourceAPIModel{
			Type: self.Type.ValueString(),
			// Metadata:
			// Properties:
			AllocatePerInstance: self.AllocatePerInstance.ValueBoolPointer(),
		},
		diags
}

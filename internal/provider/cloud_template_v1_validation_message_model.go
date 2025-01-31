// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// CloudTemplateV1ValidationMessageModel describes the resource data model.
type CloudTemplateV1ValidationMessageModel struct {
	ResourceName types.String `tfsdk:"resource_name"`
	Path         types.String `tfsdk:"path"`
	Message      types.String `tfsdk:"message"`
}

// CloudTemplateV1ValidationMessageAPIModel describes the resource API model.
type CloudTemplateV1ValidationMessageAPIModel struct {
	ResourceName string `json:"resourceName,omitempty"`
	Path         string `json:"path"`
	Message      string `json:"message"`
}

func (self CloudTemplateV1ValidationMessageModel) String() string {
	return "Cloud Template v1 Validation Message"
}

func (self *CloudTemplateV1ValidationMessageModel) FromAPI(
	raw CloudTemplateV1ValidationMessageAPIModel,
) {
	self.ResourceName = types.StringValue(raw.ResourceName)
	self.Path = types.StringValue(raw.Path)
	self.Message = types.StringValue(raw.Message)
}

func (self CloudTemplateV1ValidationMessageModel) ToAPI() CloudTemplateV1ValidationMessageAPIModel {
	return CloudTemplateV1ValidationMessageAPIModel{
		ResourceName: self.ResourceName.ValueString(),
		Path:         self.Path.ValueString(),
		Message:      self.Message.ValueString(),
	}
}

// Convert an object to a CloudTemplateV1ValidationMessageAPIModel.
func CloudTemplateV1ValidationMessageAPIModelFromObject(
	ctx context.Context,
	object types.Object,
) (*CloudTemplateV1ValidationMessageAPIModel, diag.Diagnostics) {

	if object.IsNull() || object.IsUnknown() {
		return nil, diag.Diagnostics{}
	}

	var message CloudTemplateV1ValidationMessageModel
	diags := object.As(ctx, &message, basetypes.ObjectAsOptions{})
	raw := message.ToAPI()
	return &raw, diags
}

// Convert a CloudTemplateV1ValidationMessageAPIModel to an object.
func (self *CloudTemplateV1ValidationMessageAPIModel) ToObject(
	ctx context.Context,
) (types.Object, diag.Diagnostics) {
	message := CloudTemplateV1ValidationMessageModel{}
	if self == nil {
		return types.ObjectNull(message.AttributeTypes()), diag.Diagnostics{}
	}
	message.FromAPI(*self)
	return types.ObjectValueFrom(ctx, message.AttributeTypes(), message)
}

// Used to convert structure to a types.Object.
func (self CloudTemplateV1ValidationMessageModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"resource_name": types.StringType,
		"path":          types.StringType,
		"message":       types.StringType,
	}
}

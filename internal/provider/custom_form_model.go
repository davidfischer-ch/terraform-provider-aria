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

// CustomFormModel describes the resource data model.
type CustomFormModel struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Type       types.String `tfsdk:"type"`
	Form       types.String `tfsdk:"form"` // TODO A struct to define this attribute
	FormFormat types.String `tfsdk:"form_format"`
	Styles     types.String `tfsdk:"styles"`
	SourceId   types.String `tfsdk:"source_id"`
	SourceType types.String `tfsdk:"source_type"`
	Tenant     types.String `tfsdk:"tenant"`
	Status     types.String `tfsdk:"status"`
}

// CustomFormAPIModel describes the resource API model.
type CustomFormAPIModel struct {
	Id         string `json:"id,omitempty"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Form       string `json:"form"` // TODO A struct to define this attribute
	FormFormat string `json:"formFormat"`
	Styles     string `json:"styles"`
	SourceId   string `json:"sourceId"`
	SourceType string `json:"sourceType"`
	Tenant     string `json:"tenant"`
	Status     string `json:"status"`
}

func (self *CustomFormModel) FromAPI(
	ctx context.Context,
	raw CustomFormAPIModel,
) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Type = types.StringValue(raw.Type)
	self.Form = types.StringValue(raw.Form)
	self.FormFormat = types.StringValue(raw.FormFormat)
	self.Styles = types.StringValue(raw.Styles)
	self.SourceId = types.StringValue(raw.SourceId)
	self.SourceType = types.StringValue(raw.SourceType)
	self.Tenant = types.StringValue(raw.Tenant)
	self.Status = types.StringValue(raw.Status)
	return diag.Diagnostics{}
}

func (self *CustomFormModel) ToAPI(
	ctx context.Context,
) (CustomFormAPIModel, diag.Diagnostics) {
	return CustomFormAPIModel{
		Id:         self.Id.ValueString(),
		Name:       self.Name.ValueString(),
		Type:       self.Type.ValueString(),
		Form:       self.Form.ValueString(),
		FormFormat: self.FormFormat.ValueString(),
		Styles:     self.Styles.ValueString(),
		SourceId:   self.SourceId.ValueString(),
		SourceType: self.SourceType.ValueString(),
		Tenant:     self.Tenant.ValueString(),
		Status:     self.Status.ValueString(),
	}, diag.Diagnostics{}
}

// Convert an object to a CustomFormAPIModel
func CustomFormAPIModelFromObject(
	ctx context.Context,
	object types.Object,
) (*CustomFormAPIModel, diag.Diagnostics) {

	if object.IsNull() || object.IsUnknown() {
		return nil, diag.Diagnostics{}
	}

	var formDefinition CustomFormModel
	diags := object.As(ctx, &formDefinition, basetypes.ObjectAsOptions{})
	raw, formDiags := formDefinition.ToAPI(ctx)
	diags.Append(formDiags...)
	return &raw, diags
}

// Convert a CustomFormAPIModel to an object
func (self *CustomFormAPIModel) ToObject(
	ctx context.Context,
) (types.Object, diag.Diagnostics) {
	if self == nil {
		return types.ObjectNull(CustomFormModelAttributeTypes()), diag.Diagnostics{}
	}
	form := CustomFormModel{}
	diags := form.FromAPI(ctx, *self)
	object, objectDiags := types.ObjectValueFrom(ctx, CustomFormModelAttributeTypes(), form)
	diags.Append(objectDiags...)
	return object, diags
}

// Used to convert structure to a types.Object
func CustomFormModelAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":          types.StringType,
		"name":        types.StringType,
		"type":        types.StringType,
		"form":        types.StringType,
		"form_format": types.StringType,
		"styles":      types.StringType,
		"source_id":   types.StringType,
		"source_type": types.StringType,
		"tenant":      types.StringType,
		"status":      types.StringType,
	}
}

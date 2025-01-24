// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// CustomFormModel describes the resource data model.
type CustomFormModel struct {
	Id         types.String         `tfsdk:"id"`
	Name       types.String         `tfsdk:"name"`
	Type       types.String         `tfsdk:"type"`
	Form       jsontypes.Normalized `tfsdk:"form"` // TODO A struct to define this attribute
	FormFormat types.String         `tfsdk:"form_format"`
	Styles     types.String         `tfsdk:"styles"`
	SourceId   types.String         `tfsdk:"source_id"`
	SourceType types.String         `tfsdk:"source_type"`
	Tenant     types.String         `tfsdk:"tenant"`
	Status     types.String         `tfsdk:"status"`
}

// CustomFormAPIModel describes the resource API model.
type CustomFormAPIModel struct {
	Id         string `json:"id,omitempty"`
	Name       string `json:"name"`
	Type       string `json:"type,omitempty"`
	Form       string `json:"form,omitempty"` // TODO A struct to define this attribute
	FormFormat string `json:"formFormat,omitempty"`
	Styles     string `json:"styles,omitempty"`
	SourceId   string `json:"sourceId,omitempty"`
	SourceType string `json:"sourceType,omitempty"`
	Tenant     string `json:"tenant,omitempty"`
	Status     string `json:"status,omitempty"`
}

func (self *CustomFormModel) String() string {
	return fmt.Sprintf(
		"Custom Form %s (%s)",
		self.Id.ValueString(),
		self.Name.ValueString())
}

// Return an appropriate key that can be used for naming mutexes.
// Create: Identifier can be used to prevent concurrent creation of custom forms.
// Read Update Delete: Identifier can be used to prevent concurrent modifications on the instance.
func (self CustomFormModel) LockKey() string {
	return "custom-form-" + self.Id.ValueString()
}

func (self CustomFormModel) CreatePath() string {
	return "form-service/api/forms"
}

func (self CustomFormModel) ReadPath() string {
	return "form-service/api/forms/" + self.Id.ValueString()
}

func (self CustomFormModel) UpdatePath() string {
	return self.CreatePath() // Its not a mistake ...
}

func (self CustomFormModel) DeletePath() string {
	return self.ReadPath()
}

func (self *CustomFormModel) GenerateId() {
	if len(self.Id.ValueString()) == 0 {
		self.Id = types.StringValue(uuid.New().String())
	}
}

func (self *CustomFormModel) FromAPI(
	ctx context.Context,
	raw CustomFormAPIModel,
) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Type = types.StringValue(raw.Type)
	self.Form = jsontypes.NewNormalizedValue(raw.Form)
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

// Utils -------------------------------------------------------------------------------------------

// Convert an object to a CustomFormAPIModel.
func CustomFormAPIModelFromObject(
	ctx context.Context,
	object types.Object,
) (*CustomFormAPIModel, diag.Diagnostics) {

	if object.IsNull() || object.IsUnknown() {
		return nil, diag.Diagnostics{}
	}

	formDefinition := CustomFormModel{}
	diags := object.As(ctx, &formDefinition, basetypes.ObjectAsOptions{})
	raw, someDiags := formDefinition.ToAPI(ctx)
	diags.Append(someDiags...)
	return &raw, diags
}

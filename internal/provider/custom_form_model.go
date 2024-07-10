// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

/*
"form": "{\"layout\":{\"pages\":[{\"id\":\"page_1\",\"title\":\"Page 1\",\"sections\":[]}]},\"schema\":{}}"
"formFormat": "JSON"
"name": "update-sonde"
"sourceId": "Custom.POC.API.FAX.DOPI_v1.custom.update-sonde"
"sourceType": "resource.action"
"status": "ON"
"tenant": "2817c6e5-7408-449f-a86d-8f511105e5ba"
"type": "requestForm"
*/

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

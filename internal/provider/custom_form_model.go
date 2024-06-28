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
	Id types.String `tfsdk:"id"`
	// FIXME other fields
}

// CustomFormAPIModel describes the resource API model.
type CustomFormAPIModel struct {
	Id string `json:"id"`
	// FIXME other fields
}

func (self *CustomFormModel) FromAPI(
	ctx context.Context,
	raw CustomFormAPIModel,
) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	// FIXME other fields
	return diag.Diagnostics{}
}

func (self *CustomFormModel) ToAPI(
	ctx context.Context,
) (CustomFormAPIModel, diag.Diagnostics) {
	return CustomFormAPIModel{
		Id: self.Id.ValueString(),
		// FIXME other fields
	}, diag.Diagnostics{}
}

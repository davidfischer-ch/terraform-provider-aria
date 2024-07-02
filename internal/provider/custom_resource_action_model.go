// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CustomResourceActionModel describes the resource data model.
type CustomResourceActionModel struct {
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Type            types.String `tfsdk:"type"`
	ProjectId       types.String `tfsdk:"project_id"`
	InputParameters types.List   `tfsdk:"input_parameters"`
}

// CustomResourceActionAPIModel describes the resource API model.
type CustomResourceActionAPIModel struct {
	Id              string   `json:"id"`
	Name            string   `json:"name"`
	Type            string   `json:"type"`
	ProjectId       string   `json:"projectId"`
	InputParameters []string `json:"inputParameters"`
}

func (self *CustomResourceActionModel) FromAPI(
	ctx context.Context,
	raw CustomResourceActionAPIModel,
) diag.Diagnostics {
	InputParameters, diags := types.ListValueFrom(ctx, types.StringType, raw.InputParameters)

	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Type = types.StringValue(raw.Type)
	self.ProjectId = types.StringValue(raw.ProjectId)
	self.InputParameters = InputParameters

	return diags
}

func (self *CustomResourceActionModel) ToAPI(
	ctx context.Context,
) (CustomResourceActionAPIModel, diag.Diagnostics) {

	var diags diag.Diagnostics

	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/list
	if self.InputParameters.IsNull() || self.InputParameters.IsUnknown() {
		diags.AddError(
			"Configuration error",
			fmt.Sprintf(
				"Unable to manage custom resource %s, input_parameters is either null or unknown",
				self.Id.ValueString()))
		return CustomResourceActionAPIModel{}, diags
	}

	InputParameters := make([]string, 0, len(self.InputParameters.Elements()))
	diags = self.InputParameters.ElementsAs(ctx, &InputParameters, false)
	if diags.HasError() {
		return CustomResourceActionAPIModel{}, diags
	}

	return CustomResourceActionAPIModel{
		Id:              self.Id.ValueString(),
		Name:            self.Name.ValueString(),
		Type:            self.Type.ValueString(),
		ProjectId:       self.ProjectId.ValueString(),
		InputParameters: InputParameters,
	}, diag.Diagnostics{}
}

/* func (self *CustomResourceActionModel) ToObject(
    ctx context.Context,
) (types.Object, diag.Diagnostics) {
    return types.ObjectValueFrom(ctx, map[string]attr.Type{
        "id": types.StringType,
        "name": types.StringType,
        "type": types.StringType,
        "project_id": types.StringType,
        "input_parameters": types.ListType{
            ElemType: types.StringType,
        },
    }, self)
} */

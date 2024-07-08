// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CustomResourceActionModel describes the resource data model.
type CustomResourceActionModel struct {
	Id               types.String           `tfsdk:"id"`
	Name             types.String           `tfsdk:"name"`
	Type             types.String           `tfsdk:"type"`
	ProjectId        types.String           `tfsdk:"project_id"`
	InputParameters  []ActionParameterModel `tfsdk:"input_parameters"`
	OutputParameters []ActionParameterModel `tfsdk:"output_parameters"`
}

// CustomResourceActionAPIModel describes the resource API model.
type CustomResourceActionAPIModel struct {
	Id               string                    `json:"id"`
	Name             string                    `json:"name"`
	Type             string                    `json:"type"`
	ProjectId        string                    `json:"projectId"`
	InputParameters  []ActionParameterAPIModel `json:"inputParameters"`
	OutputParameters []ActionParameterAPIModel `json:"outputParameters"`
}

func (self *CustomResourceActionModel) FromAPI(
	ctx context.Context,
	raw CustomResourceActionAPIModel,
) diag.Diagnostics {

	diags := diag.Diagnostics{}

	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Type = types.StringValue(raw.Type)
	self.ProjectId = types.StringValue(raw.ProjectId)

	self.InputParameters = []ActionParameterModel{}
	for _, parameterItem := range raw.InputParameters {
		parameter := ActionParameterModel{}
		diags.Append(parameter.FromAPI(ctx, parameterItem)...)
		self.InputParameters = append(self.InputParameters, parameter)
	}

	self.OutputParameters = []ActionParameterModel{}
	for _, parameterItem := range raw.OutputParameters {
		parameter := ActionParameterModel{}
		diags.Append(parameter.FromAPI(ctx, parameterItem)...)
		self.OutputParameters = append(self.OutputParameters, parameter)
	}

	return diags
}

func (self *CustomResourceActionModel) ToAPI(
	ctx context.Context,
) (CustomResourceActionAPIModel, diag.Diagnostics) {

	diags := diag.Diagnostics{}

	inputParametersRaw := []ActionParameterAPIModel{}
	for _, parameter := range self.InputParameters {
		parameterRaw, parameterDiags := parameter.ToAPI(ctx)
		inputParametersRaw = append(inputParametersRaw, parameterRaw)
		diags.Append(parameterDiags...)
	}

	outputParametersRaw := []ActionParameterAPIModel{}
	for _, parameter := range self.InputParameters {
		parameterRaw, parameterDiags := parameter.ToAPI(ctx)
		outputParametersRaw = append(outputParametersRaw, parameterRaw)
		diags.Append(parameterDiags...)
	}

	return CustomResourceActionAPIModel{
		Id:               self.Id.ValueString(),
		Name:             self.Name.ValueString(),
		Type:             self.Type.ValueString(),
		ProjectId:        self.ProjectId.ValueString(),
		InputParameters:  inputParametersRaw,
		OutputParameters: outputParametersRaw,
	}, diag.Diagnostics{}
}

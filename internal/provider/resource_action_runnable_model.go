// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResourceActionRunnableModel describes the resource data model.
type ResourceActionRunnableModel struct {
	Id               types.String     `tfsdk:"id"`
	Name             types.String     `tfsdk:"name"`
	Type             types.String     `tfsdk:"type"`
	ProjectId        types.String     `tfsdk:"project_id"`
	InputParameters  []ParameterModel `tfsdk:"input_parameters"`
	OutputParameters []ParameterModel `tfsdk:"output_parameters"`
}

// ResourceActionRunnableAPIModel describes the resource API model.
type ResourceActionRunnableAPIModel struct {
	Id               string              `json:"id,omitempty"`
	Name             string              `json:"name"`
	Type             string              `json:"type"`
	ProjectId        string              `json:"projectId"`
	InputParameters  []ParameterAPIModel `json:"inputParameters"`
	OutputParameters []ParameterAPIModel `json:"outputParameters"`
}

func (self *ResourceActionRunnableModel) String() string {
	return fmt.Sprintf(
		"Resource Action Runnable %s (%s) project %s",
		self.Id.ValueString(),
		self.Name.ValueString(),
		self.ProjectId.ValueString())
}

func (self *ResourceActionRunnableModel) FromAPI(
	ctx context.Context,
	raw ResourceActionRunnableAPIModel,
) diag.Diagnostics {

	diags := diag.Diagnostics{}

	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Type = types.StringValue(raw.Type)
	self.ProjectId = types.StringValue(raw.ProjectId)

	self.InputParameters = []ParameterModel{}
	for _, parameterItem := range raw.InputParameters {
		parameter := ParameterModel{}
		diags.Append(parameter.FromAPI(ctx, parameterItem)...)
		self.InputParameters = append(self.InputParameters, parameter)
	}

	self.OutputParameters = []ParameterModel{}
	for _, parameterItem := range raw.OutputParameters {
		parameter := ParameterModel{}
		diags.Append(parameter.FromAPI(ctx, parameterItem)...)
		self.OutputParameters = append(self.OutputParameters, parameter)
	}

	return diags
}

func (self ResourceActionRunnableModel) ToAPI(
	ctx context.Context,
) (ResourceActionRunnableAPIModel, diag.Diagnostics) {

	diags := diag.Diagnostics{}

	inputParametersRaw := []ParameterAPIModel{}
	for _, parameter := range self.InputParameters {
		parameterRaw, parameterDiags := parameter.ToAPI(ctx)
		inputParametersRaw = append(inputParametersRaw, parameterRaw)
		diags.Append(parameterDiags...)
	}

	outputParametersRaw := []ParameterAPIModel{}
	for _, parameter := range self.InputParameters {
		parameterRaw, parameterDiags := parameter.ToAPI(ctx)
		outputParametersRaw = append(outputParametersRaw, parameterRaw)
		diags.Append(parameterDiags...)
	}

	return ResourceActionRunnableAPIModel{
		Id:               self.Id.ValueString(),
		Name:             self.Name.ValueString(),
		Type:             self.Type.ValueString(),
		ProjectId:        self.ProjectId.ValueString(),
		InputParameters:  inputParametersRaw,
		OutputParameters: outputParametersRaw,
	}, diag.Diagnostics{}
}

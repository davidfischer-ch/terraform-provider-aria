// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResourceActionRunnableModel describes the resource data model.
type ResourceActionRunnableModel struct {
	Id               types.String         `tfsdk:"id"`
	Name             types.String         `tfsdk:"name"`
	Type             types.String         `tfsdk:"type"`
	ProjectId        types.String         `tfsdk:"project_id"`
	EndpointLink     types.String         `tfsdk:"endpoint_link"`
	InputBindings    jsontypes.Normalized `tfsdk:"input_bindings"`
	InputParameters  []ParameterModel     `tfsdk:"input_parameters"`
	OutputParameters []ParameterModel     `tfsdk:"output_parameters"`
}

// ResourceActionRunnableAPIModel describes the resource API model.
type ResourceActionRunnableAPIModel struct {
	Id               string              `json:"id,omitempty"`
	Name             string              `json:"name"`
	Type             string              `json:"type"`
	ProjectId        string              `json:"projectId,omitempty"`
	EndpointLink     string              `json:"endpointLink,omitempty"`
	InputBindings    any                 `json:"inputBindings"`
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
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Type = types.StringValue(raw.Type)
	self.ProjectId = types.StringValue(raw.ProjectId)
	self.EndpointLink = types.StringValue(raw.EndpointLink)

	var diags diag.Diagnostics
	self.InputBindings, diags = JSONNormalizedFromAny(self.String(), raw.InputBindings)

	self.InputParameters = []ParameterModel{}
	for _, parameterItem := range raw.InputParameters {
		parameter := ParameterModel{}
		parameter.FromAPI(parameterItem)
		self.InputParameters = append(self.InputParameters, parameter)
	}

	self.OutputParameters = []ParameterModel{}
	for _, parameterItem := range raw.OutputParameters {
		parameter := ParameterModel{}
		parameter.FromAPI(parameterItem)
		self.OutputParameters = append(self.OutputParameters, parameter)
	}

	return diags
}

func (self ResourceActionRunnableModel) ToAPI(
	ctx context.Context,
) (ResourceActionRunnableAPIModel, diag.Diagnostics) {

	// InputBindings JSON Encoded -> API data
	inputBindingsRaw, diags := JSONNormalizedToAny(self.InputBindings)

	inputParametersRaw := []ParameterAPIModel{}
	for _, parameter := range self.InputParameters {
		inputParametersRaw = append(inputParametersRaw, parameter.ToAPI())
	}

	outputParametersRaw := []ParameterAPIModel{}
	for _, parameter := range self.OutputParameters {
		outputParametersRaw = append(outputParametersRaw, parameter.ToAPI())
	}

	return ResourceActionRunnableAPIModel{
		Id:               self.Id.ValueString(),
		Name:             self.Name.ValueString(),
		Type:             self.Type.ValueString(),
		ProjectId:        self.ProjectId.ValueString(),
		EndpointLink:     self.EndpointLink.ValueString(),
		InputBindings:    inputBindingsRaw,
		InputParameters:  inputParametersRaw,
		OutputParameters: outputParametersRaw,
	}, diags
}

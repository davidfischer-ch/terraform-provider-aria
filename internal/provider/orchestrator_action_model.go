// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// OrchestratorActionModel describes the resource data model.
type OrchestratorActionModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Module      types.String `tfsdk:"module"`
	FQN         types.String `tfsdk:"fqn"`
	Description types.String `tfsdk:"description"`
	Version     types.String `tfsdk:"version"`

	Runtime            types.String `tfsdk:"runtime"`
	RuntimeMemoryLimit types.Int64  `tfsdk:"runtime_memory_limit"`
	RuntimeTimeout     types.Int32  `tfsdk:"runtime_timeout"`

	Script types.String `tfsdk:"script"`

	InputParameters types.List `tfsdk:"input_parameters"`
	// Of type OrchestratorActionInputParameterModel

	OutputType types.String `tfsdk:"output_type"`
}

// OrchestratorActionAPIModel describes the resource API model.
type OrchestratorActionAPIModel struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Module      string `json:"module"`
	FQN         string `json:"fqn"`
	Description string `json:"description"`
	Version     string `json:"version"`

	Runtime            string `json:"runtime,omitempty"`
	RuntimeMemoryLimit int64  `json:"runtimeMemoryLimit"`
	RuntimeTimeout     int32  `json:"runtimeTimeout"`

	Script string `json:"script"`

	InputParameters []OrchestratorActionInputParameterAPIModel `json:"input-parameters,omitempty"`

	OutputType string `json:"output-type"`
}

func (self OrchestratorActionModel) String() string {
	return fmt.Sprintf(
		"Orchestrator Action %s (%s)",
		self.Id.ValueString(),
		self.FQN.ValueString())
}

// Return an appropriate key that can be used for naming mutexes.
// Create: Identifier can be used to prevent concurrent creation of vRO actions.
// Read Update Delete: Identifier can be used to prevent concurrent modifications on the instance.
func (self OrchestratorActionModel) LockKey() string {
	return "orchestrator-action-" + self.Id.ValueString()
}

func (self OrchestratorActionModel) CreatePath() string {
	return "vco/api/actions"
}

func (self OrchestratorActionModel) ReadPath() string {
	return fmt.Sprintf("vco/api/actions/%s", self.Id.ValueString())
}

func (self OrchestratorActionModel) UpdatePath() string {
	return self.ReadPath()
}

func (self OrchestratorActionModel) DeletePath() string {
	return self.ReadPath()
}

func (self *OrchestratorActionModel) FromAPI(
	ctx context.Context,
	raw OrchestratorActionAPIModel,
) diag.Diagnostics {

	diags := diag.Diagnostics{}

	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Module = types.StringValue(raw.Module)
	self.FQN = types.StringValue(raw.FQN)
	self.Description = types.StringValue(raw.Description)
	self.Version = types.StringValue(raw.Version)

	self.Runtime = types.StringValue(raw.Runtime)
	self.RuntimeMemoryLimit = types.Int64Value(raw.RuntimeMemoryLimit)
	self.RuntimeTimeout = types.Int32Value(raw.RuntimeTimeout)
	self.Script = types.StringValue(raw.Script)
	self.OutputType = types.StringValue(raw.OutputType)

	// Convert input parameters from raw
	// Ensure array (not nil) to make practitioner's life easier
	parameters := []OrchestratorActionInputParameterModel{}
	if raw.InputParameters != nil {
		for _, parameterRaw := range raw.InputParameters {
			parameter := OrchestratorActionInputParameterModel{}
			diags.Append(parameter.FromAPI(ctx, parameterRaw)...)
			parameters = append(parameters, parameter)
		}
	}

	// Store inputs parameters to list value
	var parametersDiags diag.Diagnostics
	parameterAttrs := types.ObjectType{AttrTypes: OrchestratorActionInputParameterAttributeTypes()}
	self.InputParameters, parametersDiags = types.ListValueFrom(ctx, parameterAttrs, parameters)
	diags.Append(parametersDiags...)

	return diags
}

func (self OrchestratorActionModel) ToAPI(
	ctx context.Context,
) (OrchestratorActionAPIModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/list
	if self.InputParameters.IsNull() || self.InputParameters.IsUnknown() {
		diags.AddError(
			"Configuration error",
			fmt.Sprintf(
				"Unable to manage %s, input_parameters is either null or unknown",
				self.String()))
		return OrchestratorActionAPIModel{}, diags
	}

	// Extract input parameters from list value
	parameters := make(
		[]OrchestratorActionInputParameterModel,
		0, len(self.InputParameters.Elements()),
	)
	diags.Append(self.InputParameters.ElementsAs(ctx, &parameters, false)...)

	// Convert input parameters to raw
	parametersRaw := []OrchestratorActionInputParameterAPIModel{}
	for _, parameter := range parameters {
		parameterRaw, parameterDiags := parameter.ToAPI(ctx)
		parametersRaw = append(parametersRaw, parameterRaw)
		diags.Append(parameterDiags...)
	}

	return OrchestratorActionAPIModel{
		Id: self.Id.ValueString(),
		Name: self.Name.ValueString(),
		Module: self.Module.ValueString(),
		FQN: self.FQN.ValueString(),
		Description: CleanString(self.Description.ValueString()),
		Version: self.Version.ValueString(),
		Runtime: self.Runtime.ValueString(),
		RuntimeMemoryLimit: self.RuntimeMemoryLimit.ValueInt64(),
		RuntimeTimeout: self.RuntimeTimeout.ValueInt32(),
		Script: self.Script.ValueString(),
		InputParameters: parametersRaw,
		OutputType: self.OutputType.ValueString(),
	}, diag.Diagnostics{}
}

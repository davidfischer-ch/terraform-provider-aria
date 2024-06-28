// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ABXActionModel describes the resource data model.
type ABXActionModel struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	FAASProvider types.String `tfsdk:"faas_provider"`
	Type         types.String `tfsdk:"type"`

	RuntimeName    types.String `tfsdk:"runtime_name"`
	RuntimeVersion types.String `tfsdk:"runtime_version"`

	CPUShares                types.Int64 `tfsdk:"cpu_shares"`
	MemoryInMB               types.Int64 `tfsdk:"memory_in_mb"`
	TimeoutSeconds           types.Int64 `tfsdk:"timeout_seconds"`
	DeploymentTimeoutSeconds types.Int64 `tfsdk:"deployment_timeout_seconds"`

	Entrypoint   types.String `tfsdk:"entrypoint"`
	Dependencies types.List   `tfsdk:"dependencies"`
	Constants    types.Set    `tfsdk:"constants"`
	Secrets      types.Set    `tfsdk:"secrets"`

	Source types.String `tfsdk:"source"`

	Shared        types.Bool `tfsdk:"shared"`
	System        types.Bool `tfsdk:"system"`
	AsyncDeployed types.Bool `tfsdk:"async_deployed"`

	ProjectId types.String `tfsdk:"project_id"`
	OrgId     types.String `tfsdk:"org_id"`
}

// ABXActionAPIModel describes the resource API model.
type ABXActionAPIModel struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	FAASProvider string `json:"provider"`
	Type         string `json:"actionType"`

	RuntimeName    string `json:"runtime"`
	RuntimeVersion string `json:"runtimeVersion"`

	CPUShares                int64 `json:"cpuShares"`
	MemoryInMB               int64 `json:"memoryInMB"`
	TimeoutSeconds           int64 `json:"timeoutSeconds"`
	DeploymentTimeoutSeconds int64 `json:"deploymentTimeoutSeconds"`

	Entrypoint   string            `json:"entrypoint"`
	Dependencies string            `json:"dependencies"`
	Inputs       map[string]string `json:"inputs"`

	Source string `json:"source"`

	Shared        bool `json:"shared"`
	System        bool `json:"system"`
	AsyncDeployed bool `json:"asyncDeployed"`

	ProjectId string `json:"projectId"`
	OrgId     string `json:"orgId"`
}

func (self *ABXActionModel) FromAPI(
	ctx context.Context,
	raw ABXActionAPIModel,
) diag.Diagnostics {

	diags := diag.Diagnostics{}

	faasProvider := raw.FAASProvider
	if faasProvider == "" {
		faasProvider = "auto"
	}

	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Description = types.StringValue(raw.Description)
	self.FAASProvider = types.StringValue(faasProvider)
	self.RuntimeName = types.StringValue(raw.RuntimeName)
	self.RuntimeVersion = types.StringValue(raw.RuntimeVersion)
	self.CPUShares = types.Int64Value(raw.CPUShares)
	self.MemoryInMB = types.Int64Value(raw.MemoryInMB)
	self.TimeoutSeconds = types.Int64Value(raw.TimeoutSeconds)
	self.DeploymentTimeoutSeconds = types.Int64Value(raw.DeploymentTimeoutSeconds)
	self.Entrypoint = types.StringValue(raw.Entrypoint)
	self.Source = types.StringValue(CleanString(raw.Source))
	self.Shared = types.BoolValue(raw.Shared)
	self.System = types.BoolValue(raw.System)
	self.AsyncDeployed = types.BoolValue(raw.AsyncDeployed)
	self.ProjectId = types.StringValue(raw.ProjectId)
	self.OrgId = types.StringValue(raw.OrgId)

	constantsIds := []string{}
	secretsIds := []string{}
	inputsKeys := []string{}

	for key := range raw.Inputs {
		res := strings.SplitN(key, ":", 2)
		if len(res) == 1 {
			inputsKeys = append(inputsKeys, key)
		} else if res[0] == "secret" {
			constantsIds = append(constantsIds, res[1])
		} else if res[0] == "psecret" {
			secretsIds = append(secretsIds, res[1])
		} else {
			// Unhandled -> inputsKeys
			inputsKeys = append(inputsKeys, key)
		}
	}

	if len(inputsKeys) > 0 {
		diags.AddError(
			"Client error",
			fmt.Sprintf(
				"Unable to manage ABX action %s, unhandled inputs keys %s",
				self.Id.ValueString(), strings.Join(inputsKeys, ", ")))
	}

	constants, constantsDiags := types.SetValueFrom(ctx, types.StringType, constantsIds)
	self.Constants = constants
	diags.Append(constantsDiags...)

	secrets, secretsDiags := types.SetValueFrom(ctx, types.StringType, secretsIds)
	self.Secrets = secrets
	diags.Append(secretsDiags...)

	dependencies, dependenciesDiags := types.ListValueFrom(
		ctx, types.StringType, SkipEmpty(strings.Split(raw.Dependencies, "\n")),
	)
	self.Dependencies = dependencies
	diags.Append(dependenciesDiags...)

	return diags
}

func (self *ABXActionModel) ToAPI(ctx context.Context) (ABXActionAPIModel, diag.Diagnostics) {

	diags := diag.Diagnostics{}

	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/list
	if self.Dependencies.IsNull() || self.Dependencies.IsUnknown() {
		diags.AddError(
			"Configuration error",
			fmt.Sprintf(
				"Unable to manage ABX action %s, dependencies is either null or unknown",
				self.Id.ValueString()))
		return ABXActionAPIModel{}, diags
	}

	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/list
	if self.Constants.IsNull() || self.Constants.IsUnknown() {
		diags.AddError(
			"Configuration error",
			fmt.Sprintf(
				"Unable to manage ABX action %s, constants is either null or unknown",
				self.Id.ValueString()))
		return ABXActionAPIModel{}, diags
	}

	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/list
	if self.Secrets.IsNull() || self.Secrets.IsUnknown() {
		diags.AddError(
			"Configuration error",
			fmt.Sprintf(
				"Unable to manage ABX action %s, secrets is either null or unknown",
				self.Id.ValueString()))
		return ABXActionAPIModel{}, diags
	}

	dependencies := make([]string, 0, len(self.Dependencies.Elements()))
	diags.Append(self.Dependencies.ElementsAs(ctx, &dependencies, false)...)

	constants := make([]string, 0, len(self.Constants.Elements()))
	diags.Append(self.Constants.ElementsAs(ctx, &constants, false)...)

	secrets := make([]string, 0, len(self.Secrets.Elements()))
	diags.Append(self.Secrets.ElementsAs(ctx, &secrets, false)...)

	if diags.HasError() {
		return ABXActionAPIModel{}, diags
	}

	inputs := map[string]string{}
	for _, constant := range constants {
		inputs["secret:"+constant] = ""
	}
	for _, secret := range secrets {
		inputs["psecret:"+secret] = ""
	}

	faasProvider := self.FAASProvider.ValueString()
	if faasProvider == "auto" {
		faasProvider = ""
	}

	return ABXActionAPIModel{
		Name:                     self.Name.ValueString(),
		Description:              self.Description.ValueString(),
		FAASProvider:             faasProvider,
		Type:                     self.Type.ValueString(),
		RuntimeName:              self.RuntimeName.ValueString(),
		RuntimeVersion:           self.RuntimeVersion.ValueString(),
		CPUShares:                self.CPUShares.ValueInt64(),
		MemoryInMB:               self.MemoryInMB.ValueInt64(),
		TimeoutSeconds:           self.TimeoutSeconds.ValueInt64(),
		DeploymentTimeoutSeconds: self.DeploymentTimeoutSeconds.ValueInt64(),
		Entrypoint:               self.Entrypoint.ValueString(),
		Dependencies:             strings.Join(SkipEmpty(dependencies), "\n"),
		Inputs:                   inputs,
		Source:                   CleanString(self.Source.ValueString()),
		Shared:                   self.Shared.ValueBool(),
		System:                   self.System.ValueBool(),
		AsyncDeployed:            self.AsyncDeployed.ValueBool(),
		ProjectId:                self.ProjectId.ValueString(),
	}, diags
}

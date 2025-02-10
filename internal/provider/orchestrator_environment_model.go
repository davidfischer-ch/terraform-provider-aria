// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// OrchestratorEnvironmentModel describes the resource data model.
type OrchestratorEnvironmentModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Version     types.String `tfsdk:"version"`
	VersionId   types.String `tfsdk:"version_id"`

	Runtime            types.String `tfsdk:"runtime"`
	RuntimeMemoryLimit types.Int64  `tfsdk:"runtime_memory_limit"`
	RuntimeTimeout     types.Int32  `tfsdk:"runtime_timeout"`

	Dependencies types.Map `tfsdk:"dependencies"`
	Repositories types.Map `tfsdk:"repositories"`
	Variables    types.Map `tfsdk:"variables"`

	BundleHash                     types.String `tfsdk:"bundle_hash"`
	DependenciesInstallExecutionId types.String `tfsdk:"dependencies_install_execution_id"`
	Status                         types.String `tfsdk:"status"`
	ValidationMessage              types.String `tfsdk:"validation_message"`

	WaitUpToDate types.Bool `tfsdk:"wait_up_to_date"`
}

// OrchestratorEnvironmentAPIModel describes the resource API model.
type OrchestratorEnvironmentAPIModel struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`

	Runtime            string `json:"runtime"`
	RuntimeMemoryLimit int64  `json:"runtimeMemoryLimit"`
	RuntimeTimeout     int32  `json:"runtimeTimeout"`

	Dependencies map[string]string `json:"dependencies"`
	Repositories map[string]string `json:"repositories"`
	Variables    map[string]string `json:"variables"`

	BundleHash                     string `json:"bundleHash,omitempty"`
	DependenciesInstallExecutionId string `json:"dependenciesInstallExecutionId,omitempty"`
	Status                         string `json:"status,omitempty"`
	ValidationMessage              string `json:"validationMessage,omitempty"`
}

func (self OrchestratorEnvironmentModel) String() string {
	return fmt.Sprintf(
		"Orchestrator Environment %s (%s) of %s",
		self.Id.ValueString(),
		self.Name.ValueString(),
		self.Runtime.ValueString())
}

// Return an appropriate key that can be used for naming mutexes.
// Create: Identifier can be used to prevent concurrent creation of environments.
// Read Update Delete: Identifier can be used to prevent concurrent modifications on the instance.
func (self OrchestratorEnvironmentModel) LockKey() string {
	return "orchestrator-environment-" + self.Id.ValueString()
}

func (self OrchestratorEnvironmentModel) CreatePath() string {
	return "vco/api/environments"
}

func (self OrchestratorEnvironmentModel) ReadPath() string {
	return "vco/api/environments/" + self.Id.ValueString()
}

func (self OrchestratorEnvironmentModel) UpdatePath() string {
	return self.ReadPath()
}

func (self OrchestratorEnvironmentModel) DeletePath() string {
	return self.ReadPath()
}

func (self *OrchestratorEnvironmentModel) FromAPI(
	ctx context.Context,
	raw OrchestratorEnvironmentAPIModel,
	response *resty.Response,
) diag.Diagnostics {

	prevVersionId := self.VersionId.ValueString()

	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Description = types.StringValue(raw.Description)
	self.Version = types.StringValue(raw.Version)
	self.VersionId = types.StringValue(response.Header().Get("x-vro-changeset-sha"))

	self.Runtime = types.StringValue(raw.Runtime)
	self.RuntimeMemoryLimit = types.Int64Value(raw.RuntimeMemoryLimit)
	self.RuntimeTimeout = types.Int32Value(raw.RuntimeTimeout)

	diags := diag.Diagnostics{}
	var someDiags diag.Diagnostics

	self.Dependencies, someDiags = types.MapValueFrom(ctx, types.StringType, raw.Dependencies)
	diags.Append(someDiags...)

	self.Repositories, someDiags = types.MapValueFrom(ctx, types.StringType, raw.Repositories)
	diags.Append(someDiags...)

	self.Variables, someDiags = types.MapValueFrom(ctx, types.StringType, raw.Variables)
	diags.Append(someDiags...)

	self.BundleHash = types.StringValue(raw.BundleHash)
	self.DependenciesInstallExecutionId = types.StringValue(raw.DependenciesInstallExecutionId)
	self.Status = types.StringValue(raw.Status)
	self.ValidationMessage = types.StringValue(raw.ValidationMessage)

	tflog.Debug(
		ctx,
		fmt.Sprintf(
			"%s FromAPI: versionId from response headers %s -> %s",
			self.String(), prevVersionId, self.VersionId.ValueString()))

	return diags
}

func (self OrchestratorEnvironmentModel) ToAPI(
	ctx context.Context,
) (OrchestratorEnvironmentAPIModel, diag.Diagnostics) {

	dependenciesRaw := make(map[string]string, len(self.Dependencies.Elements()))
	diags := self.Dependencies.ElementsAs(ctx, &dependenciesRaw, false)

	repositoriesRaw := make(map[string]string, len(self.Repositories.Elements()))
	diags.Append(self.Repositories.ElementsAs(ctx, &repositoriesRaw, false)...)

	variablesRaw := make(map[string]string, len(self.Variables.Elements()))
	diags.Append(self.Variables.ElementsAs(ctx, &variablesRaw, false)...)

	tflog.Debug(
		ctx,
		fmt.Sprintf("%s ToAPI: versionId %s", self.String(), self.VersionId.ValueString()))

	return OrchestratorEnvironmentAPIModel{
		Id:                             self.Id.ValueString(),
		Name:                           self.Name.ValueString(),
		Description:                    self.Description.ValueString(),
		Version:                        self.Version.ValueString(),
		Runtime:                        self.Runtime.ValueString(),
		RuntimeMemoryLimit:             self.RuntimeMemoryLimit.ValueInt64(),
		RuntimeTimeout:                 self.RuntimeTimeout.ValueInt32(),
		Dependencies:                   dependenciesRaw,
		Repositories:                   repositoriesRaw,
		Variables:                      variablesRaw,
		BundleHash:                     self.BundleHash.ValueString(),
		DependenciesInstallExecutionId: self.DependenciesInstallExecutionId.ValueString(),
		Status:                         self.Status.ValueString(),
		ValidationMessage:              self.ValidationMessage.ValueString(),
	}, diags
}

// Utils -------------------------------------------------------------------------------------------

func (self OrchestratorEnvironmentModel) IsUpToDate() bool {
	return strings.ToUpper(self.Status.ValueString()) == "UP_TO_DATE"
}

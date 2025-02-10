// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// OrchestratorEnvironmentRepositoryModel describes the resource data model.
type OrchestratorEnvironmentRepositoryModel struct {
	Id                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Runtime           types.String `tfsdk:"runtime"`
	Location          types.String `tfsdk:"location"`
	BasicAuth         types.Bool   `tfsdk:"basic_auth"`
	SystemUser        types.String `tfsdk:"system_user"`
	SystemCredentials types.String `tfsdk:"system_credentials"`
}

// OrchestratorEnvironmentRepositoryAPIModel describes the resource API model.
type OrchestratorEnvironmentRepositoryAPIModel struct {
	Id                string `json:"id,omitempty"`
	Name              string `json:"name"`
	Runtime           string `json:"runtime"`
	Location          string `json:"location"`
	BasicAuth         bool   `json:"basicAuth"`
	SystemUser        string `json:"systemUser,omitempty"`
	SystemCredentials string `json:"systemCredentials,omitempty"`
}

func (self OrchestratorEnvironmentRepositoryModel) String() string {
	return fmt.Sprintf(
		"Orchestrator Environment Repository %s (%s) %s",
		self.Id.ValueString(),
		self.Name.ValueString(),
		self.Runtime.ValueString())
}

// Return an appropriate key that can be used for naming mutexes.
// Create: Identifier can be used to prevent concurrent creation of repositories.
// Read Update Delete: Identifier can be used to prevent concurrent modifications on the instance.
func (self OrchestratorEnvironmentRepositoryModel) LockKey() string {
	return "orchestrator-environment-repository-" + self.Id.ValueString()
}

func (self OrchestratorEnvironmentRepositoryModel) CreatePath() string {
	return "vco/api/environments/repositories"
}

func (self OrchestratorEnvironmentRepositoryModel) ReadPath() string {
	return "vco/api/environments/repositories/" + self.Id.ValueString()
}

func (self OrchestratorEnvironmentRepositoryModel) UpdatePath() string {
	return self.ReadPath()
}

func (self OrchestratorEnvironmentRepositoryModel) DeletePath() string {
	return self.ReadPath()
}

func (self *OrchestratorEnvironmentRepositoryModel) FromAPI(
	raw OrchestratorEnvironmentRepositoryAPIModel,
) {
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Runtime = types.StringValue(raw.Runtime)
	self.Location = types.StringValue(raw.Location)
	self.BasicAuth = types.BoolValue(raw.BasicAuth)
	self.SystemUser = types.StringValue(raw.SystemUser)

	// The value is not returned
	// self.SystemCredentials = types.StringValue("")
}

func (self OrchestratorEnvironmentRepositoryModel) ToAPI() OrchestratorEnvironmentRepositoryAPIModel {
	self.BasicAuth = types.BoolValue(len(self.SystemUser.ValueString()) > 0)
	return OrchestratorEnvironmentRepositoryAPIModel{
		Id:                self.Id.ValueString(),
		Name:              self.Name.ValueString(),
		Runtime:           self.Runtime.ValueString(),
		Location:          self.Location.ValueString(),
		BasicAuth:         self.BasicAuth.ValueBool(),
		SystemUser:        self.SystemUser.ValueString(),
		SystemCredentials: self.SystemCredentials.ValueString(),
	}
}

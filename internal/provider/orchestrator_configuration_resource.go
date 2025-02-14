// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &OrchestratorConfigurationResource{}
var _ resource.ResourceWithImportState = &OrchestratorConfigurationResource{}

func NewOrchestratorConfigurationResource() resource.Resource {
	return &OrchestratorConfigurationResource{}
}

// OrchestratorConfigurationResource defines the resource implementation.
type OrchestratorConfigurationResource struct {
	client *AriaClient
}

func (self *OrchestratorConfigurationResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_orchestrator_configuration"
}

func (self *OrchestratorConfigurationResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = OrchestratorConfigurationSchema()
}

func (self *OrchestratorConfigurationResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *OrchestratorConfigurationResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var configuration OrchestratorConfigurationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &configuration)...)
	if resp.Diagnostics.HasError() {
		return
	}

	configurationToAPI, diags := configuration.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var configurationFromAPI OrchestratorConfigurationAPIModel
	path := configuration.CreatePath()
	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", GetVersionFromPath(path)).
		SetBody(configurationToAPI).
		SetResult(&configurationFromAPI).
		Post(path)
	err = handleAPIResponse(ctx, response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", configuration.String(), err))
		return
	}

	// Save configuration into Terraform state
	resp.Diagnostics.Append(configuration.FromAPI(ctx, configurationFromAPI, response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &configuration)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", configuration.String()))
}

func (self *OrchestratorConfigurationResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var configuration OrchestratorConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &configuration)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var configurationFromAPI OrchestratorConfigurationAPIModel
	found, response, someDiags := self.client.ReadIt(ctx, &configuration, &configurationFromAPI)
	resp.Diagnostics.Append(someDiags...)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	if !resp.Diagnostics.HasError() {
		// Save updated configuration into Terraform state
		resp.Diagnostics.Append(configuration.FromAPI(ctx, configurationFromAPI, response)...)
		resp.Diagnostics.Append(resp.State.Set(ctx, &configuration)...)
	}
}

func (self *OrchestratorConfigurationResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var configuration OrchestratorConfigurationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &configuration)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform state data into the model
	var configurationFromState OrchestratorConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &configurationFromState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	configurationToAPI, diags := configuration.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// No response body from API, only the changeset (version) available in response headers
	path := configuration.UpdatePath()
	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", GetVersionFromPath(path)).
		SetHeader("x-vro-changeset-sha", configurationFromState.VersionId.ValueString()).
		SetBody(configurationToAPI).
		Put(path)

	err = handleAPIResponse(ctx, response, err, []int{204})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", configuration.String(), err))
		return
	}

	// Save updated configuration into Terraform state
	resp.Diagnostics.Append(configuration.FromAPI(ctx, configurationToAPI, response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &configuration)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", configuration.String()))
}

func (self *OrchestratorConfigurationResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var configuration OrchestratorConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &configuration)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(self.client.DeleteIt(ctx, &configuration)...)
	}
}

func (self *OrchestratorConfigurationResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

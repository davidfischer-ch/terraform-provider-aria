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
var _ resource.Resource = &OrchestratorEnvironmentRepositoryResource{}
var _ resource.ResourceWithImportState = &OrchestratorEnvironmentRepositoryResource{}

func NewOrchestratorEnvironmentRepositoryResource() resource.Resource {
	return &OrchestratorEnvironmentRepositoryResource{}
}

// OrchestratorEnvironmentRepositoryResource defines the resource implementation.
type OrchestratorEnvironmentRepositoryResource struct {
	client *AriaClient
}

func (self *OrchestratorEnvironmentRepositoryResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_orchestrator_environment_repository"
}

func (self *OrchestratorEnvironmentRepositoryResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = OrchestratorEnvironmentRepositorySchema()
}

func (self *OrchestratorEnvironmentRepositoryResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *OrchestratorEnvironmentRepositoryResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var repository OrchestratorEnvironmentRepositoryModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &repository)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var repositoryFromAPI OrchestratorEnvironmentRepositoryAPIModel
	path := repository.CreatePath()
	response, err := self.client.R(path).
		SetBody(repository.ToAPI()).
		SetResult(&repositoryFromAPI).
		Post(path)
	err = self.client.HandleAPIResponse(response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", repository.String(), err))
		return
	}

	// Save repository into Terraform state
	repository.FromAPI(repositoryFromAPI)
	resp.Diagnostics.Append(resp.State.Set(ctx, &repository)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", repository.String()))
}

func (self *OrchestratorEnvironmentRepositoryResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var repository OrchestratorEnvironmentRepositoryModel
	resp.Diagnostics.Append(req.State.Get(ctx, &repository)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var repositoryFromAPI OrchestratorEnvironmentRepositoryAPIModel
	found, _, readDiags := self.client.ReadIt(&repository, &repositoryFromAPI)
	resp.Diagnostics.Append(readDiags...)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated repository into Terraform state
	repository.FromAPI(repositoryFromAPI)
	resp.Diagnostics.Append(resp.State.Set(ctx, &repository)...)
}

func (self *OrchestratorEnvironmentRepositoryResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var repository OrchestratorEnvironmentRepositoryModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &repository)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var repositoryFromAPI OrchestratorEnvironmentRepositoryAPIModel
	path := repository.UpdatePath()
	body := repository.ToAPI()
	response, err := self.client.R(path).SetBody(body).SetResult(&repositoryFromAPI).Put(path)
	err = self.client.HandleAPIResponse(response, err, []int{202})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", repository.String(), err))
		return
	}

	// Save updated repository into Terraform state
	repository.FromAPI(repositoryFromAPI)
	resp.Diagnostics.Append(resp.State.Set(ctx, &repository)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", repository.String()))
}

func (self *OrchestratorEnvironmentRepositoryResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var repository OrchestratorEnvironmentRepositoryModel
	resp.Diagnostics.Append(req.State.Get(ctx, &repository)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(self.client.DeleteIt(&repository)...)
	}
}

func (self *OrchestratorEnvironmentRepositoryResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("system_credentials"), "")...)
}

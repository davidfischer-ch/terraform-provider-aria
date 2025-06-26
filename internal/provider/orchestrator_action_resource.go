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
var _ resource.Resource = &OrchestratorActionResource{}
var _ resource.ResourceWithImportState = &OrchestratorActionResource{}

func NewOrchestratorActionResource() resource.Resource {
	return &OrchestratorActionResource{}
}

// OrchestratorActionResource defines the resource implementation.
type OrchestratorActionResource struct {
	client *AriaClient
}

func (self *OrchestratorActionResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_orchestrator_action"
}

func (self *OrchestratorActionResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = OrchestratorActionSchema()
}

func (self *OrchestratorActionResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *OrchestratorActionResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var action OrchestratorActionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &action)...)
	if resp.Diagnostics.HasError() {
		return
	}

	actionToAPI, diags := action.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var actionFromAPI OrchestratorActionAPIModel
	path := action.CreatePath()

	self.client.Mutex.Lock(ctx, action.LockKey())
	response, err := self.client.R(path).SetBody(actionToAPI).SetResult(&actionFromAPI).Post(path)
	self.client.Mutex.Unlock(ctx, action.LockKey())

	err = self.client.HandleAPIResponse(response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", action.String(), err))
		return
	}

	// Save action into Terraform state
	resp.Diagnostics.Append(action.FromAPI(ctx, actionFromAPI)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &action)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", action.String()))
}

func (self *OrchestratorActionResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var action OrchestratorActionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &action)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var actionFromAPI OrchestratorActionAPIModel
	self.client.Mutex.RLock(ctx, action.LockKey())
	found, _, readDiags := self.client.ReadIt(&action, &actionFromAPI)
	self.client.Mutex.RUnlock(ctx, action.LockKey())

	resp.Diagnostics.Append(readDiags...)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated action into Terraform state
	resp.Diagnostics.Append(action.FromAPI(ctx, actionFromAPI)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &action)...)
}

func (self *OrchestratorActionResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var action OrchestratorActionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &action)...)
	if resp.Diagnostics.HasError() {
		return
	}

	actionToAPI, diags := action.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	self.client.Mutex.Lock(ctx, action.LockKey())

	path := action.UpdatePath()
	response, err := self.client.R(path).SetBody(actionToAPI).Put(path)
	err = self.client.HandleAPIResponse(response, err, []int{200})
	if err != nil {
		self.client.Mutex.Unlock(ctx, action.LockKey())
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", action.String(), err))
		return
	}

	// Read (using API) to retrieve the action content (and not empty stuff)
	var actionFromAPI OrchestratorActionAPIModel
	path = action.ReadPath()
	response, err = self.client.R(path).SetResult(&actionFromAPI).Get(path)

	self.client.Mutex.Unlock(ctx, action.LockKey())

	err = self.client.HandleAPIResponse(response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read %s, got error: %s", action.String(), err))
		return
	}

	// Save updated action into Terraform state
	resp.Diagnostics.Append(action.FromAPI(ctx, actionFromAPI)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &action)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", action.String()))
}

func (self *OrchestratorActionResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var action OrchestratorActionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &action)...)
	if !resp.Diagnostics.HasError() {
		// Do not serialize deletion (with a mutex) to allow convering (if possible) when deletion
		// is not forced by some of the actions
		resp.Diagnostics.Append(self.client.DeleteIt(&action)...)
	}
}

func (self *OrchestratorActionResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("force_delete"), false)...)
}

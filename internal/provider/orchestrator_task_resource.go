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
var _ resource.Resource = &OrchestratorTaskResource{}
var _ resource.ResourceWithImportState = &OrchestratorTaskResource{}

func NewOrchestratorTaskResource() resource.Resource {
	return &OrchestratorTaskResource{}
}

// OrchestratorTaskResource defines the resource implementation.
type OrchestratorTaskResource struct {
	client *AriaClient
}

func (self *OrchestratorTaskResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_orchestrator_task"
}

func (self *OrchestratorTaskResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = OrchestratorTaskSchema()
}

func (self *OrchestratorTaskResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *OrchestratorTaskResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var task OrchestratorTaskModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &task)...)
	if resp.Diagnostics.HasError() {
		return
	}

	taskToAPI, diags := task.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var taskFromAPI OrchestratorTaskAPIModel
	path := task.CreatePath()
	response, err := self.client.R(path).SetBody(taskToAPI).SetResult(&taskFromAPI).Post(path)
	err = self.client.HandleAPIResponse(response, err, []int{202})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", task.String(), err))
		return
	}

	// Save task into Terraform state
	resp.Diagnostics.Append(task.FromAPI(ctx, taskFromAPI)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &task)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", task.String()))
}

func (self *OrchestratorTaskResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var task OrchestratorTaskModel
	resp.Diagnostics.Append(req.State.Get(ctx, &task)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var taskFromAPI OrchestratorTaskAPIModel
	found, _, readDiags := self.client.ReadIt(&task, &taskFromAPI)
	resp.Diagnostics.Append(readDiags...)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated task into Terraform state
	resp.Diagnostics.Append(task.FromAPI(ctx, taskFromAPI)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &task)...)
}

func (self *OrchestratorTaskResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var task OrchestratorTaskModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &task)...)
	if resp.Diagnostics.HasError() {
		return
	}

	taskToAPI, diags := task.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var taskFromAPI OrchestratorTaskAPIModel
	path := task.UpdatePath()
	response, err := self.client.R(path).SetBody(taskToAPI).SetResult(&taskFromAPI).Post(path)
	err = self.client.HandleAPIResponse(response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", task.String(), err))
		return
	}

	// Save updated task into Terraform state
	resp.Diagnostics.Append(task.FromAPI(ctx, taskFromAPI)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &task)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", task.String()))
}

func (self *OrchestratorTaskResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var task OrchestratorTaskModel
	resp.Diagnostics.Append(req.State.Get(ctx, &task)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(self.client.DeleteIt(&task)...)
	}
}

func (self *OrchestratorTaskResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

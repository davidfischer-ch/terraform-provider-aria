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
var _ resource.Resource = &OrchestratorWorkflowResource{}
var _ resource.ResourceWithImportState = &OrchestratorWorkflowResource{}

func NewOrchestratorWorkflowResource() resource.Resource {
	return &OrchestratorWorkflowResource{}
}

// OrchestratorWorkflowResource defines the resource implementation.
type OrchestratorWorkflowResource struct {
	client *AriaClient
}

func (self *OrchestratorWorkflowResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_orchestrator_workflow"
}

func (self *OrchestratorWorkflowResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = OrchestratorWorkflowSchema()
}

func (self *OrchestratorWorkflowResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *OrchestratorWorkflowResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var workflow OrchestratorWorkflowModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &workflow)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var workflowFromCreateAPI OrchestratorWorkflowCreateAPIModel
	path := workflow.CreatePath()
	response, err := self.client.R(path).
		SetBody(workflow.ToCreateAPI()).
		SetResult(&workflowFromCreateAPI).
		Post(path)
	err = self.client.HandleAPIResponse(response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", workflow.String(), err))
		return
	}

	// Save workflow into Terraform state
	workflow.FromCreateAPI(workflowFromCreateAPI)
	resp.Diagnostics.Append(resp.State.Set(ctx, &workflow)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully, now updating", workflow.String()))

	// Update ... TODO deduplicate with Update()

	workflowToVersionAPI, diags := workflow.ToVersionAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var workflowFromVersionAPI OrchestratorWorkflowVersionResponseAPIModel
	path = workflow.UpdatePath()
	response, err = self.client.R(path).
		SetBody(workflowToVersionAPI).
		SetResult(&workflowFromVersionAPI).
		Post(path)
	err = self.client.HandleAPIResponse(response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", workflow.String(), err))
		return
	}

	// Save updated workflow into Terraform state
	workflow.FromVersionAPI(workflowFromVersionAPI)
	resp.Diagnostics.Append(resp.State.Set(ctx, &workflow)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", workflow.String()))

	// Read ... TODO deduplicate with Read()

	// Read content
	var workflowFromContentAPI OrchestratorWorkflowContentAPIModel
	found, response, readDiags := self.client.ReadIt(&workflow, &workflowFromContentAPI)
	resp.Diagnostics.Append(readDiags...)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	// Read forms
	var fromsFromAPI any
	_, _, readDiags = self.client.ReadIt(&workflow, &fromsFromAPI, workflow.ReadFormPath())
	resp.Diagnostics.Append(readDiags...)

	if !resp.Diagnostics.HasError() {
		// Save updated workflow into Terraform state
		resp.Diagnostics.Append(workflow.FromContentAPI(ctx, workflowFromContentAPI, response)...)
		resp.Diagnostics.Append(workflow.FromFormAPI(ctx, fromsFromAPI)...)
		resp.Diagnostics.Append(resp.State.Set(ctx, &workflow)...)
	}
}

func (self *OrchestratorWorkflowResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var workflow OrchestratorWorkflowModel
	resp.Diagnostics.Append(req.State.Get(ctx, &workflow)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read content
	var workflowFromContentAPI OrchestratorWorkflowContentAPIModel
	found, response, readDiags := self.client.ReadIt(&workflow, &workflowFromContentAPI)
	resp.Diagnostics.Append(readDiags...)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	// Read forms
	var formsFromAPI any
	_, _, readDiags = self.client.ReadIt(&workflow, &formsFromAPI, workflow.ReadFormPath())
	resp.Diagnostics.Append(readDiags...)

	if !resp.Diagnostics.HasError() {
		// Save updated workflow into Terraform state
		resp.Diagnostics.Append(workflow.FromContentAPI(ctx, workflowFromContentAPI, response)...)
		resp.Diagnostics.Append(workflow.FromFormAPI(ctx, formsFromAPI)...)
		resp.Diagnostics.Append(resp.State.Set(ctx, &workflow)...)
	}
}

func (self *OrchestratorWorkflowResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var workflow OrchestratorWorkflowModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &workflow)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workflowToVersionAPI, diags := workflow.ToVersionAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var workflowFromVersionAPI OrchestratorWorkflowVersionResponseAPIModel
	path := workflow.UpdatePath()
	response, err := self.client.R(path).
		SetBody(workflowToVersionAPI).
		SetResult(&workflowFromVersionAPI).
		Post(path)
	err = self.client.HandleAPIResponse(response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", workflow.String(), err))
		return
	}

	// Save updated workflow into Terraform state
	workflow.FromVersionAPI(workflowFromVersionAPI)
	resp.Diagnostics.Append(resp.State.Set(ctx, &workflow)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", workflow.String()))
}

func (self *OrchestratorWorkflowResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var workflow OrchestratorWorkflowModel
	resp.Diagnostics.Append(req.State.Get(ctx, &workflow)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(self.client.DeleteIt(&workflow)...)
	}
}

func (self *OrchestratorWorkflowResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("force_delete"), false)...)
}

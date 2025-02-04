// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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

	workflowCreateRaw := workflow.ToCreateAPI()
	response, err := self.client.Client.R().
		// TODO SetQueryParam("apiVersion", ORCHESTRATOR_API_VERSION).
		SetBody(workflowCreateRaw).
		SetResult(&workflowCreateRaw).
		Post(workflow.CreatePath())
	err = handleAPIResponse(ctx, response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", workflow.String(), err))
		return
	}

	// Save workflow into Terraform state
	workflow.FromCreateAPI(workflowCreateRaw)
	resp.Diagnostics.Append(resp.State.Set(ctx, &workflow)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully, now updating", workflow.String()))

	// Update ... TODO deduplicate with Update()

	workflowVersionRaw, diags := workflow.ToVersionAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var workflowVersionResponsRaw OrchestratorWorkflowVersionResponseAPIModel
	response, err = self.client.Client.R().
		// TODO SetQueryParam("apiVersion", ORCHESTRATOR_API_VERSION).
		SetBody(workflowVersionRaw).
		SetResult(&workflowVersionResponsRaw).
		Post(workflow.UpdatePath())

	err = handleAPIResponse(ctx, response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", workflow.String(), err))
		return
	}

	// Save updated workflow into Terraform state
	workflow.FromVersionAPI(workflowVersionResponsRaw)
	resp.Diagnostics.Append(resp.State.Set(ctx, &workflow)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", workflow.String()))

	// Read ... TODO deduplicate with Read()

	// Read content
	var workflowContentRaw OrchestratorWorkflowContentAPIModel
	found, response, readDiags := self.client.ReadIt(ctx, &workflow, &workflowContentRaw)
	resp.Diagnostics.Append(readDiags...)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	// Read forms
	var formsRaw any
	_, _, readDiags = self.client.ReadIt(ctx, &workflow, &formsRaw, workflow.ReadFormPath())
	resp.Diagnostics.Append(readDiags...)

	if !resp.Diagnostics.HasError() {
		// Save updated workflow into Terraform state
		resp.Diagnostics.Append(workflow.FromContentAPI(ctx, workflowContentRaw, response)...)
		resp.Diagnostics.Append(workflow.FromFormAPI(ctx, formsRaw)...)
		resp.Diagnostics.Append(resp.State.Set(ctx, &workflow)...)
	}

	// Optionally wait available on catalog
	resp.Diagnostics.Append(self.WaitOnCatalog(ctx, &workflow)...)
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
	var workflowContentRaw OrchestratorWorkflowContentAPIModel
	found, response, readDiags := self.client.ReadIt(ctx, &workflow, &workflowContentRaw)
	resp.Diagnostics.Append(readDiags...)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	// Read forms
	var formsRaw any
	_, _, readDiags = self.client.ReadIt(ctx, &workflow, &formsRaw, workflow.ReadFormPath())
	resp.Diagnostics.Append(readDiags...)

	if !resp.Diagnostics.HasError() {
		// Save updated workflow into Terraform state
		resp.Diagnostics.Append(workflow.FromContentAPI(ctx, workflowContentRaw, response)...)
		resp.Diagnostics.Append(workflow.FromFormAPI(ctx, formsRaw)...)
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

	workflowVersionRaw, diags := workflow.ToVersionAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var workflowVersionResponsRaw OrchestratorWorkflowVersionResponseAPIModel
	response, err := self.client.Client.R().
		// TODO SetQueryParam("apiVersion", ORCHESTRATOR_API_VERSION).
		SetBody(workflowVersionRaw).
		SetResult(&workflowVersionResponsRaw).
		Post(workflow.UpdatePath())

	err = handleAPIResponse(ctx, response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", workflow.String(), err))
		return
	}

	// Save updated workflow into Terraform state
	workflow.FromVersionAPI(workflowVersionResponsRaw)
	resp.Diagnostics.Append(resp.State.Set(ctx, &workflow)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", workflow.String()))

	// Optionally wait available on catalog
	resp.Diagnostics.Append(self.WaitOnCatalog(ctx, &workflow)...)
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
		resp.Diagnostics.Append(self.client.DeleteIt(ctx, &workflow)...)
	}
}

func (self *OrchestratorWorkflowResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("force_delete"), false)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("wait_on_catalog"), true)...)
}

// -------------------------------------------------------------------------------------------------

func (self *OrchestratorWorkflowResource) WaitOnCatalog(
	ctx context.Context,
	workflow *OrchestratorWorkflowModel,
) diag.Diagnostics {

	diags := diag.Diagnostics{}

	if !workflow.WaitOnCatalog.ValueBool() {
		return diags
	}

	name := workflow.String()
	path := workflow.ReadCatalogPath()
	tflog.Debug(ctx, fmt.Sprintf("Wait %s to be available on catalog...", name))

	// Poll for matching catalog item to be available up to 10 minutes (60 x 10 seconds)
	maxAttempts := 60
	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Poll resource until imported
		time.Sleep(time.Duration(10) * time.Second)
		tflog.Debug(
			ctx,
			fmt.Sprintf("Poll %d of %d - Check %s is imported...", attempt+1, maxAttempts, name))

		response, err := self.client.Client.R().
			SetQueryParam("apiVersion", GetVersionFromPath(path)).
			Get(path)

		if response.StatusCode() == 404 {
			tflog.Debug(ctx, fmt.Sprintf("%s not found", name))
			continue
		}

		err = handleAPIResponse(ctx, response, err, []int{200})
		if err != nil {
			diags.AddError(
				"Client error",
				fmt.Sprintf("Unable to read %s, got error: %s", name, err))
		}

		// Found
		return diags
	}

	diags.AddError(
		"Client error", fmt.Sprintf("Timeout while waiting for %s to be imported.", name))
	return diags
}

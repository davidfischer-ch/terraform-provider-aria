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
var _ resource.Resource = &OrchestratorCategoryResource{}
var _ resource.ResourceWithImportState = &OrchestratorCategoryResource{}

func NewOrchestratorCategoryResource() resource.Resource {
	return &OrchestratorCategoryResource{}
}

// OrchestratorCategoryResource defines the resource implementation.
type OrchestratorCategoryResource struct {
	client *AriaClient
}

func (self *OrchestratorCategoryResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_orchestrator_category"
}

func (self *OrchestratorCategoryResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = OrchestratorCategorySchema()
}

func (self *OrchestratorCategoryResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *OrchestratorCategoryResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var category OrchestratorCategoryModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &category)...)
	if resp.Diagnostics.HasError() {
		return
	}

	categoryRaw, diags := category.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.Client.R().
		// TODO SetQueryParam("apiVersion", ORCHESTRATOR_API_VERSION).
		SetBody(categoryRaw).
		SetResult(&categoryRaw).
		Post(category.CreatePath())
	err = handleAPIResponse(ctx, response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", category.String(), err))
		return
	}

	// Save category into Terraform state
	resp.Diagnostics.Append(category.FromAPI(ctx, categoryRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &category)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", category.String()))
}

func (self *OrchestratorCategoryResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var category OrchestratorCategoryModel
	resp.Diagnostics.Append(req.State.Get(ctx, &category)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var categoryRaw OrchestratorCategoryAPIModel
	found, readDiags := self.client.ReadIt(ctx, &category, &categoryRaw)
	resp.Diagnostics.Append(readDiags...)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	if !resp.Diagnostics.HasError() {
		// Save updated category into Terraform state
		resp.Diagnostics.Append(category.FromAPI(ctx, categoryRaw)...)
		resp.Diagnostics.Append(resp.State.Set(ctx, &category)...)
	}
}

func (self *OrchestratorCategoryResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var category OrchestratorCategoryModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &category)...)
	if resp.Diagnostics.HasError() {
		return
	}

	categoryRaw, diags := category.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.Client.R().
		// TODO SetQueryParam("apiVersion", ORCHESTRATOR_API_VERSION).
		SetBody(categoryRaw).
		Put(category.UpdatePath())

	err = handleAPIResponse(ctx, response, err, []int{204})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", category.String(), err))
		return
	}

	// Read (using API) to retrieve the category content (and not empty stuff)
	response, err = self.client.Client.R().
		// TODO SetQueryParam("apiVersion", ORCHESTRATOR_API_VERSION).
		SetResult(&categoryRaw).
		Get(category.ReadPath())

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read %s, got error: %s", category.String(), err))
		return
	}

	// Save updated category into Terraform state
	resp.Diagnostics.Append(category.FromAPI(ctx, categoryRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &category)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", category.String()))
}

func (self *OrchestratorCategoryResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var category OrchestratorCategoryModel
	resp.Diagnostics.Append(req.State.Get(ctx, &category)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(self.client.DeleteIt(ctx, &category)...)
	}
}

func (self *OrchestratorCategoryResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

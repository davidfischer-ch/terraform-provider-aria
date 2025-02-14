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
var _ resource.Resource = &ProjectResource{}
var _ resource.ResourceWithImportState = &ProjectResource{}

func NewProjectResource() resource.Resource {
	return &ProjectResource{}
}

// ProjectResource defines the resource implementation.
type ProjectResource struct {
	client *AriaClient
}

func (self *ProjectResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (self *ProjectResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = ProjectSchema()
}

func (self *ProjectResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *ProjectResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var project ProjectModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &project)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectToAPI, diags := project.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var projectFromAPI ProjectAPIModel
	path := project.CreatePath()
	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", GetVersionFromPath(path)).
		SetQueryParam("validatePrincipals", "true").
		SetQueryParam("syncPrincipals", "true").
		SetBody(projectToAPI).
		SetResult(&projectFromAPI).
		Post(project.CreatePath())
	err = handleAPIResponse(ctx, response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", project.String(), err))
		return
	}

	// Save property group into Terraform state
	resp.Diagnostics.Append(project.FromAPI(ctx, projectFromAPI)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &project)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", project.String()))
}

func (self *ProjectResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var project ProjectModel
	resp.Diagnostics.Append(req.State.Get(ctx, &project)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var projectFromAPI ProjectAPIModel
	found, _, readDiags := self.client.ReadIt(ctx, &project, &projectFromAPI)
	resp.Diagnostics.Append(readDiags...)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	if !resp.Diagnostics.HasError() {
		// Save updated property group into Terraform state
		resp.Diagnostics.Append(project.FromAPI(ctx, projectFromAPI)...)
		resp.Diagnostics.Append(resp.State.Set(ctx, &project)...)
	}
}

func (self *ProjectResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var project ProjectModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &project)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectToAPI, diags := project.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var projectFromAPI ProjectAPIModel
	path := project.UpdatePath()
	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", GetVersionFromPath(path)).
		SetQueryParam("validatePrincipals", "true").
		SetQueryParam("syncPrincipals", "true").
		SetBody(projectToAPI).
		SetResult(&projectFromAPI).
		Patch(path)

	// TODO Also call PATCH project-service/api/projects/{id}/cost
	// TODO Also call PATCH project-service/api/projects/{id}/principals
	// TODO Also call PATCH project-service/api/projects/{id}/resource-metadata

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", project.String(), err))
		return
	}

	// Save updated property group into Terraform state
	resp.Diagnostics.Append(project.FromAPI(ctx, projectFromAPI)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &project)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", project.String()))
}

func (self *ProjectResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var project ProjectModel
	resp.Diagnostics.Append(req.State.Get(ctx, &project)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(self.client.DeleteIt(ctx, &project)...)
	}
}

func (self *ProjectResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// FIXME must be filtered by id and projectId
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

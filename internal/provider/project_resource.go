// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
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
	client *resty.Client
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
	var Project ProjectModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &Project)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ProjectRaw, diags := Project.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.R().
		SetQueryParam("apiVersion", PROJECT_API_VERSION).
		SetQueryParam("validatePrincipals", "true").
		SetQueryParam("syncPrincipals", "true").
		SetBody(ProjectRaw).
		SetResult(&ProjectRaw).
		Post("project-service/api/projects")
	err = handleAPIResponse(ctx, response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", Project.String(), err))
		return
	}

	// Save property group into Terraform state
	resp.Diagnostics.Append(Project.FromAPI(ctx, ProjectRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &Project)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", Project.String()))
}

func (self *ProjectResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var Project ProjectModel
	resp.Diagnostics.Append(req.State.Get(ctx, &Project)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ProjectId := Project.Id.ValueString()
	var ProjectRaw ProjectAPIModel
	response, err := self.client.R().
		SetQueryParam("apiVersion", PROJECT_API_VERSION).
		SetResult(&ProjectRaw).
		Get("project-service/api/projects/" + ProjectId)

	// Handle gracefully a resource that has vanished on the platform
	// Beware that some APIs respond with HTTP 404 instead of 403 ...
	if response.StatusCode() == 404 {
		tflog.Debug(ctx, fmt.Sprintf("%s not found", Project.String()))
		resp.State.RemoveResource(ctx)
		return
	}

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read %s, got error: %s", Project.String(), err))
		return
	}

	// Save updated property group into Terraform state
	resp.Diagnostics.Append(Project.FromAPI(ctx, ProjectRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &Project)...)
}

func (self *ProjectResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var Project ProjectModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &Project)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ProjectId := Project.Id.ValueString()
	ProjectRaw, diags := Project.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.R().
		SetQueryParam("apiVersion", PROJECT_API_VERSION).
		SetQueryParam("validatePrincipals", "true").
		SetQueryParam("syncPrincipals", "true").
		SetBody(ProjectRaw).
		SetResult(&ProjectRaw).
		Patch("project-service/api/projects/" + ProjectId)

	// TODO Also call PATCH project-service/api/projects/{id}/cost
	// TODO Also call PATCH project-service/api/projects/{id}/principals
	// TODO Also call PATCH project-service/api/projects/{id}/resource-metadata

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", Project.String(), err))
		return
	}

	// Save updated property group into Terraform state
	resp.Diagnostics.Append(Project.FromAPI(ctx, ProjectRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &Project)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", Project.String()))
}

func (self *ProjectResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var Project ProjectModel
	resp.Diagnostics.Append(req.State.Get(ctx, &Project)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ProjectId := Project.Id.ValueString()
	if len(ProjectId) == 0 {
		return
	}

	resp.Diagnostics.Append(
		DeleteIt(
			self.client,
			ctx,
			Project.String(),
			"project-service/api/projects/"+ProjectId,
			PROJECT_API_VERSION,
		)...,
	)
}

func (self *ProjectResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// FIXME must be filtered by id and projectId
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

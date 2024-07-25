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
var _ resource.Resource = &ABXActionResource{}
var _ resource.ResourceWithImportState = &ABXActionResource{}

func NewABXActionResource() resource.Resource {
	return &ABXActionResource{}
}

// ABXActionResource defines the resource implementation.
type ABXActionResource struct {
	client *resty.Client
}

func (self *ABXActionResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_abx_action"
}

func (self *ABXActionResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = ABXActionSchema()
}

func (self *ABXActionResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *ABXActionResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var action ABXActionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &action)...)
	if resp.Diagnostics.HasError() {
		return
	}

	actionRaw, diags := action.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.R().
		SetQueryParam("apiVersion", ABX_API_VERSION).
		SetBody(actionRaw).
		SetResult(&actionRaw).
		Post("abx/api/resources/actions")
	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", action.String(), err))
		return
	}

	// Save action into Terraform state
	resp.Diagnostics.Append(action.FromAPI(ctx, actionRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &action)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", action.String()))
}

func (self *ABXActionResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var action ABXActionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &action)...)
	if resp.Diagnostics.HasError() {
		return
	}

	actionId := action.Id.ValueString()
	projectId := action.ProjectId.ValueString()
	var actionRaw ABXActionAPIModel
	response, err := self.client.R().
		SetQueryParam("apiVersion", ABX_API_VERSION).
		SetResult(&actionRaw).
		Get(fmt.Sprintf("abx/api/resources/actions/%s?projectId=%s", actionId, projectId))

	// Handle gracefully a resource that has vanished on the platform
	// Beware that some APIs respond with HTTP 404 instead of 403 ...
	if response.StatusCode() == 404 {
		tflog.Debug(ctx, fmt.Sprintf("%s not found", action.String()))
		resp.State.RemoveResource(ctx)
		return
	}

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read %s, got error: %s", action.String(), err))
		return
	}

	// Save updated action into Terraform state
	resp.Diagnostics.Append(action.FromAPI(ctx, actionRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &action)...)
}

func (self *ABXActionResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var action ABXActionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &action)...)
	if resp.Diagnostics.HasError() {
		return
	}

	actionId := action.Id.ValueString()
	projectId := action.ProjectId.ValueString()
	actionRaw, diags := action.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.R().
		SetQueryParam("apiVersion", ABX_API_VERSION).
		SetBody(actionRaw).
		SetResult(&actionRaw).
		Put(fmt.Sprintf("abx/api/resources/actions/%s?projectId=%s", actionId, projectId))

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", action.String(), err))
		return
	}

	// Save updated action into Terraform state
	resp.Diagnostics.Append(action.FromAPI(ctx, actionRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &action)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", action.String()))
}

func (self *ABXActionResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var action ABXActionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &action)...)
	if resp.Diagnostics.HasError() {
		return
	}

	actionId := action.Id.ValueString()
	projectId := action.ProjectId.ValueString()
	if len(actionId) == 0 {
		return
	}

	resp.Diagnostics.Append(
		DeleteIt(
			self.client,
			ctx,
			action.String(),
			fmt.Sprintf("abx/api/resources/actions/%s?projectId=%s", actionId, projectId),
			ABX_API_VERSION,
		)...,
	)
}

func (self *ABXActionResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// FIXME must be filtered by id and projectId
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

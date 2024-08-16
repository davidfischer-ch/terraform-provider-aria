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
var _ resource.Resource = &ResourceActionResource{}
var _ resource.ResourceWithImportState = &ResourceActionResource{}

func NewResourceActionResource() resource.Resource {
	return &ResourceActionResource{}
}

// ResourceActionResource defines the resource implementation.
type ResourceActionResource struct {
	client *AriaClient
}

func (self *ResourceActionResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_resource_action"
}

func (self *ResourceActionResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = ResourceActionSchema()
}

func (self *ResourceActionResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *ResourceActionResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var action ResourceActionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &action)...)
	if resp.Diagnostics.HasError() {
		return
	}

	actionRaw, diags := action.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceId := action.ResourceId.ValueString()
	if len(resourceId) > 0 {
		// Custom resource ...
		resp.Diagnostics.AddError("Client error", "Not yet implemented")
		return
	} else {
		// Native resource ...
		response, err := self.client.Client.R().
			SetQueryParam("apiVersion", FORM_API_VERSION).
			SetBody(actionRaw).
			SetResult(&actionRaw).
			Post(action.CreatePath())
		err = handleAPIResponse(ctx, response, err, []int{200})
		if err != nil {
			resp.Diagnostics.AddError(
				"Client error",
				fmt.Sprintf("Unable to create %s, got error: %s", action.String(), err))
			return
		}
	}

	// Save resource action into Terraform state
	resp.Diagnostics.Append(action.FromAPI(ctx, actionRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &action)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", action.String()))
}

func (self *ResourceActionResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var action ResourceActionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &action)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var actionRaw ResourceActionAPIModel
	found, readDiags := self.client.ReadIt(ctx, &action, &actionRaw)
	resp.Diagnostics.Append(readDiags...)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	if !resp.Diagnostics.HasError() {
		// Save updated resource action into Terraform state
		resp.Diagnostics.Append(action.FromAPI(ctx, actionRaw)...)
		resp.Diagnostics.Append(resp.State.Set(ctx, &action)...)
	}
}

func (self *ResourceActionResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var action ResourceActionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &action)...)
	if resp.Diagnostics.HasError() {
		return
	}

	actionRaw, diags := action.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceId := action.ResourceId.ValueString()
	if len(resourceId) > 0 {
		// Custom resource ...
		resp.Diagnostics.AddError("Client error", "Not yet implemented")
		return
	} else {
		// Native resource ...
		response, err := self.client.Client.R().
			SetQueryParam("apiVersion", FORM_API_VERSION).
			SetBody(actionRaw).
			SetResult(&actionRaw).
			Post(action.UpdatePath())

		err = handleAPIResponse(ctx, response, err, []int{200})
		if err != nil {
			resp.Diagnostics.AddError(
				"Client error",
				fmt.Sprintf("Unable to update %s, got error: %s", action.String(), err))
			return
		}
	}

	// Save updated resource action into Terraform state
	resp.Diagnostics.Append(action.FromAPI(ctx, actionRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &action)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", action.String()))
}

func (self *ResourceActionResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var action ResourceActionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &action)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(action.ResourceId.ValueString()) > 0 {
		// Custom resource ...
		resp.Diagnostics.AddError("Client error", "Not yet implemented")
		return
	} else {
		// Native resource ...
		resp.Diagnostics.Append(self.client.DeleteIt(ctx, &action)...)
	}
}

func (self *ResourceActionResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// FIXME must be filtered by id and projectId
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

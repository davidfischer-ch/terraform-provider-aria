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
var _ resource.Resource = &CustomResourceResource{}
var _ resource.ResourceWithImportState = &CustomResourceResource{}

func NewCustomResourceResource() resource.Resource {
	return &CustomResourceResource{}
}

// CustomResourceResource defines the resource implementation.
type CustomResourceResource struct {
	client *AriaClient
}

func (self *CustomResourceResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_custom_resource"
}

func (self *CustomResourceResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = CustomResourceSchema()
}

func (self *CustomResourceResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *CustomResourceResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var resource CustomResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resource)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceToAPI, diags := resource.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var resourceFromAPI CustomResourceAPIModel
	path := resource.CreatePath()
	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", GetVersionFromPath(path)).
		SetBody(resourceToAPI).
		SetResult(&resourceFromAPI).
		Post(path)
	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", resource.String(), err))
		return
	}

	// Save custom resource into Terraform state
	resp.Diagnostics.Append(resource.FromAPI(ctx, resourceFromAPI)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &resource)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", resource.String()))
}

func (self *CustomResourceResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var resource CustomResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &resource)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var resourceFromAPI CustomResourceAPIModel
	self.client.Mutex.RLock(ctx, resource.LockKey())
	found, _, diags := self.client.ReadIt(ctx, &resource, &resourceFromAPI)
	self.client.Mutex.RUnlock(ctx, resource.LockKey())
	resp.Diagnostics.Append(diags...)

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	if !resp.Diagnostics.HasError() {
		// Save updated custom resource into Terraform state
		resp.Diagnostics.Append(resource.FromAPI(ctx, resourceFromAPI)...)
		resp.Diagnostics.Append(resp.State.Set(ctx, &resource)...)
	}
}

func (self *CustomResourceResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var resource CustomResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resource)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceToAPI, diags := resource.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	self.client.Mutex.Lock(ctx, resource.LockKey())

	// Read resource to retrieve latest value for additional actions
	var resourceFromAPI CustomResourceAPIModel
	found, _, diags := self.client.ReadIt(ctx, &resource, &resourceFromAPI)
	resp.Diagnostics.Append(diags...)

	if !found || resp.Diagnostics.HasError() {
		if !found {
			resp.Diagnostics.AddError(
				"Client error",
				fmt.Sprintf(
					"Unable to update %s: Not found.",
					resource.String()))
		}

		self.client.Mutex.Unlock(ctx, resource.LockKey())
		return
	}

	// Ensure additional actions are left untouched
	resourceToAPI.AdditionalActions = resourceFromAPI.AdditionalActions

	// Reset to prevent muxing of old/new data
	resourceFromAPI = CustomResourceAPIModel{}
	path := resource.UpdatePath()
	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", GetVersionFromPath(path)).
		SetBody(resourceToAPI).
		SetResult(&resourceFromAPI).
		Post(path)

	self.client.Mutex.Unlock(ctx, resource.LockKey())

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", resource.String(), err))
		return
	}

	// Save updated custom resource into Terraform state
	resp.Diagnostics.Append(resource.FromAPI(ctx, resourceFromAPI)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &resource)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", resource.String()))
}

func (self *CustomResourceResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var resource CustomResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &resource)...)
	if !resp.Diagnostics.HasError() {
		self.client.Mutex.Lock(ctx, resource.LockKey())
		resp.Diagnostics.Append(self.client.DeleteIt(ctx, &resource)...)
		self.client.Mutex.Unlock(ctx, resource.LockKey())
	}
}

func (self *CustomResourceResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// FIXME must be filtered by id and projectId
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

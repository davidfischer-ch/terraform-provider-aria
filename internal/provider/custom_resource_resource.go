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

	resourceRaw, diags := resource.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", FORM_API_VERSION).
		SetBody(resourceRaw).
		SetResult(&resourceRaw).
		Post(resource.CreatePath())
	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", resource.String(), err))
		return
	}

	// Save custom resource into Terraform state
	resp.Diagnostics.Append(resource.FromAPI(ctx, resourceRaw)...)
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

	var raw CustomResourceAPIModel
	found, readDiags := self.client.ReadIt(ctx, &resource, &raw)
	resp.Diagnostics.Append(readDiags...)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	if !resp.Diagnostics.HasError() {
		// Save updated custom resource into Terraform state
		resp.Diagnostics.Append(resource.FromAPI(ctx, raw)...)
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

	resourceRaw, diags := resource.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// FIXME Mutex
	// FIXME Read resource to retrieve additional actions to keep them untouched

	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", FORM_API_VERSION).
		SetBody(resourceRaw).
		SetResult(&resourceRaw).
		Post(resource.UpdatePath())

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", resource.String(), err))
		return
	}

	// Save updated custom resource into Terraform state
	resp.Diagnostics.Append(resource.FromAPI(ctx, resourceRaw)...)
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
		resp.Diagnostics.Append(self.client.DeleteIt(ctx, &resource)...)
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

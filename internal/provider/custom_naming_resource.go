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
var _ resource.Resource = &CustomNamingResource{}
var _ resource.ResourceWithImportState = &CustomNamingResource{}

func NewCustomNamingResource() resource.Resource {
	return &CustomNamingResource{}
}

// CustomNamingResource defines the resource implementation.
type CustomNamingResource struct {
	client *AriaClient
}

func (self *CustomNamingResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_custom_naming"
}

func (self *CustomNamingResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = CustomNamingSchema()
}

func (self *CustomNamingResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *CustomNamingResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var naming CustomNamingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &naming)...)
	if resp.Diagnostics.HasError() {
		return
	}

	namingRaw, diags := naming.ToAPI(ctx, CustomNamingModel{})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", IAAS_API_VERSION).
		SetBody(namingRaw).
		SetResult(&namingRaw).
		Post(naming.CreatePath())
	err = handleAPIResponse(ctx, response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", naming.String(), err))
		return
	}

	// Save custom naming into Terraform state
	resp.Diagnostics.Append(naming.FromAPI(ctx, namingRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &naming)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", naming.String()))
}

func (self *CustomNamingResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var naming CustomNamingModel
	resp.Diagnostics.Append(req.State.Get(ctx, &naming)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var namingRaw CustomNamingAPIModel
	found, _, readDiags := self.client.ReadIt(ctx, &naming, &namingRaw)
	resp.Diagnostics.Append(readDiags...)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	if !resp.Diagnostics.HasError() {
		// Save updated custom naming into Terraform state
		resp.Diagnostics.Append(naming.FromAPI(ctx, namingRaw)...)
		resp.Diagnostics.Append(resp.State.Set(ctx, &naming)...)
	}
}

func (self *CustomNamingResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan and state data into the model
	var naming, namingState CustomNamingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &naming)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &namingState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	namingRaw, diags := naming.ToAPI(ctx, namingState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", IAAS_API_VERSION).
		SetBody(namingRaw).
		SetResult(&namingRaw).
		Put(naming.UpdatePath())

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", naming.String(), err))
		return
	}

	// Save updated custom naming into Terraform state
	resp.Diagnostics.Append(naming.FromAPI(ctx, namingRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &naming)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", naming.String()))
}

func (self *CustomNamingResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var naming CustomNamingModel
	resp.Diagnostics.Append(req.State.Get(ctx, &naming)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(self.client.DeleteIt(ctx, &naming)...)
	}
}

func (self *CustomNamingResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// FIXME must be filtered by id and projectId
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

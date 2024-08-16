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
var _ resource.Resource = &ABXConstantResource{}
var _ resource.ResourceWithImportState = &ABXConstantResource{}

func NewABXConstantResource() resource.Resource {
	return &ABXConstantResource{}
}

// ABXConstantResource defines the resource implementation.
type ABXConstantResource struct {
	client *AriaClient
}

func (self *ABXConstantResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_abx_constant"
}

func (self *ABXConstantResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = ABXConstantSchema()
}

func (self *ABXConstantResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *ABXConstantResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var constant ABXConstantModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &constant)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var constantRaw ABXConstantAPIModel
	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", ABX_API_VERSION).
		SetBody(constant.ToAPI()).
		SetResult(&constantRaw).
		Post(constant.CreatePath())

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", constant.String(), err))
		return
	}

	// Save constant into Terraform state
	resp.Diagnostics.Append(constant.FromAPI(ctx, constantRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &constant)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", constant.String()))
}

func (self *ABXConstantResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var constant ABXConstantModel
	resp.Diagnostics.Append(req.State.Get(ctx, &constant)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var constantRaw ABXConstantAPIModel
	found, readDiags := self.client.ReadIt(ctx, &constant, &constantRaw)
	resp.Diagnostics.Append(readDiags...)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	if !resp.Diagnostics.HasError() {
		// Save updated constant into Terraform state
		resp.Diagnostics.Append(constant.FromAPI(ctx, constantRaw)...)
		resp.Diagnostics.Append(resp.State.Set(ctx, &constant)...)
	}
}

func (self *ABXConstantResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var constant ABXConstantModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &constant)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var constantRaw ABXConstantAPIModel
	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", ABX_API_VERSION).
		SetBody(constant.ToAPI()).
		SetResult(&constantRaw).
		Put(constant.UpdatePath())

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", constant.String(), err))
		return
	}

	// Save constant into Terraform state
	resp.Diagnostics.Append(constant.FromAPI(ctx, constantRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &constant)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", constant.String()))
}

func (self *ABXConstantResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var constant ABXConstantModel
	resp.Diagnostics.Append(req.State.Get(ctx, &constant)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(self.client.DeleteIt(ctx, &constant)...)
	}
}

func (self *ABXConstantResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

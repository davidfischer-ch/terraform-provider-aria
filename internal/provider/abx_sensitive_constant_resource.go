// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ABXSensitiveConstantResource{}

func NewABXSensitiveConstantResource() resource.Resource {
	return &ABXSensitiveConstantResource{}
}

// ABXSensitiveConstantResource defines the resource implementation.
type ABXSensitiveConstantResource struct {
	client *AriaClient
}

func (self *ABXSensitiveConstantResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_abx_sensitive_constant"
}

func (self *ABXSensitiveConstantResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = ABXSensitiveConstantSchema()
}

func (self *ABXSensitiveConstantResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *ABXSensitiveConstantResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var constant ABXSensitiveConstantModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &constant)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var constantFromAPI ABXSensitiveConstantAPIModel
	path := constant.CreatePath()
	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", GetVersionFromPath(path)).
		SetBody(constant.ToAPI()).
		SetResult(&constantFromAPI).
		Post(path)

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create ABX Sensitive Constant, got error: %s", err))
		return
	}

	// Save sensitive constant into Terraform state
	constant.FromAPI(constantFromAPI)
	resp.Diagnostics.Append(resp.State.Set(ctx, &constant)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", constant.String()))
}

func (self *ABXSensitiveConstantResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var constant ABXSensitiveConstantModel
	resp.Diagnostics.Append(req.State.Get(ctx, &constant)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var constantFromAPI ABXSensitiveConstantAPIModel
	found, _, readDiags := self.client.ReadIt(ctx, &constant, &constantFromAPI)
	resp.Diagnostics.Append(readDiags...)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	if !resp.Diagnostics.HasError() {
		// Save updated secret into Terraform state
		constant.FromAPI(constantFromAPI)
		resp.Diagnostics.Append(resp.State.Set(ctx, &constant)...)
	}
}

func (self *ABXSensitiveConstantResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var constant ABXSensitiveConstantModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &constant)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var constantFromAPI ABXSensitiveConstantAPIModel
	path := constant.UpdatePath()
	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", GetVersionFromPath(path)).
		SetBody(constant.ToAPI()).
		SetResult(&constantFromAPI).
		Put(path)

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", constant.String(), err))
		return
	}

	// Save sensitive constant into Terraform state
	constant.FromAPI(constantFromAPI)
	resp.Diagnostics.Append(resp.State.Set(ctx, &constant)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", constant.String()))
}

func (self *ABXSensitiveConstantResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var constant ABXSensitiveConstantModel
	resp.Diagnostics.Append(req.State.Get(ctx, &constant)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(self.client.DeleteIt(ctx, &constant)...)
	}
}

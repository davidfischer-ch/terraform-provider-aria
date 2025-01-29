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
var _ resource.Resource = &PolicyResource{}
var _ resource.ResourceWithImportState = &PolicyResource{}

func NewPolicyResource() resource.Resource {
	return &PolicyResource{}
}

// PolicyResource defines the resource implementation.
type PolicyResource struct {
	client *AriaClient
}

func (self *PolicyResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_policy"
}

func (self *PolicyResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = PolicySchema()
}

func (self *PolicyResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *PolicyResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var policy PolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &policy)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyRaw, diags := policy.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", POLICY_API_VERSION).
		SetBody(policyRaw).
		SetResult(&policyRaw).
		Post(policy.CreatePath())
	err = handleAPIResponse(ctx, response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", policy.String(), err))
		return
	}

	// Save policy into Terraform state
	resp.Diagnostics.Append(policy.FromAPI(ctx, policyRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &policy)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", policy.String()))
}

func (self *PolicyResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var policy PolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &policy)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var policyRaw PolicyAPIModel
	found, _, readDiags := self.client.ReadIt(ctx, &policy, &policyRaw)
	resp.Diagnostics.Append(readDiags...)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	if !resp.Diagnostics.HasError() {
		// Save updated policy into Terraform state
		resp.Diagnostics.Append(policy.FromAPI(ctx, policyRaw)...)
		resp.Diagnostics.Append(resp.State.Set(ctx, &policy)...)
	}
}

func (self *PolicyResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var policy PolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &policy)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyRaw, diags := policy.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", FORM_API_VERSION).
		SetBody(policyRaw).
		SetResult(&policyRaw).
		Post(policy.UpdatePath())
	err = handleAPIResponse(ctx, response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", policy.String(), err))
		return
	}

	// Save policy into Terraform state
	resp.Diagnostics.Append(policy.FromAPI(ctx, policyRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &policy)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", policy.String()))
}

func (self *PolicyResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var policy PolicyModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &policy)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(self.client.DeleteIt(ctx, &policy)...)
	}
}

func (self *PolicyResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

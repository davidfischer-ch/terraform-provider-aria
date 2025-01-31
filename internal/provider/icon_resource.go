// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IconResource{}

func NewIconResource() resource.Resource {
	return &IconResource{}
}

// IconResource defines the resource implementation.
type IconResource struct {
	client *AriaClient
}

func (self *IconResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_icon"
}

func (self *IconResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = IconSchema()
}

func (self *IconResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *IconResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var icon IconModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &icon)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Creating "icons" with the same content multiple times in parallel
	// will lead to a platform's Internal Error (HTTP 500)...
	// The platform is not handling properly concurrent requests to icon create/delete API
	// So we implement this protection (mutex) at the client side (provider)

	lockKey := icon.LockKey()
	self.client.Mutex.Lock(ctx, lockKey)
	response, err := self.client.Client.R().
		// TODO SetQueryParam("apiVersion", ICON_API_VERSION).
		SetFile("file", icon.Path.ValueString()).
		Post(icon.CreatePath())

	err = handleAPIResponse(ctx, response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", icon.String(), err))
		return
	}

	iconId, err := GetIdFromLocation(response)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to parse Icon ID, got error: %s", err))
		return
	}

	// Save icon into Terraform state
	icon.Id = types.StringValue(iconId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &icon)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", icon.String()))

	// Read the icon to retrieve its content (duplicated code with read)

	response, err = self.client.Client.R().
		// TODO SetQueryParam("apiVersion", ICON_API_VERSION).
		Get(icon.ReadPath())

	self.client.Mutex.Unlock(ctx, lockKey)

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read %s, got error: %s", icon.String(), err))
		return
	}

	// Save updated icon into Terraform state
	icon.Hash = types.StringValue(fmt.Sprintf("%x", sha256.Sum256(response.Body())))
	resp.Diagnostics.Append(resp.State.Set(ctx, &icon)...)
	tflog.Debug(ctx, fmt.Sprintf("Refreshed %s successfully", icon.String()))

}

func (self *IconResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var icon IconModel
	resp.Diagnostics.Append(req.State.Get(ctx, &icon)...)
	if resp.Diagnostics.HasError() {
		return
	}

	self.client.Mutex.RLock(ctx, icon.LockKey())
	response, err := self.client.Client.R().
		// TODO SetQueryParam("apiVersion", ICON_API_VERSION).
		Get(icon.ReadPath())
	self.client.Mutex.RUnlock(ctx, icon.LockKey())

	// Handle gracefully a resource that has vanished on the platform
	// Beware that some APIs respond with HTTP 404 instead of 403 ...
	if response.StatusCode() == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read %s, got error: %s", icon.String(), err))
		return
	}

	// Save updated icon into Terraform state
	icon.Hash = types.StringValue(fmt.Sprintf("%x", sha256.Sum256(response.Body())))
	resp.Diagnostics.Append(resp.State.Set(ctx, &icon)...)
}

func (self *IconResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform prior state data into the model
	var icon IconModel
	resp.Diagnostics.Append(req.State.Get(ctx, &icon)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddError(
		"Client error",
		fmt.Sprintf("Cannot update %s, this type of resource is immutable.", icon.String()))
}

func (self *IconResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var icon IconModel
	resp.Diagnostics.Append(req.State.Get(ctx, &icon)...)
	if !resp.Diagnostics.HasError() && !icon.KeepOnDestroy.ValueBool() {
		self.client.Mutex.Lock(ctx, icon.LockKey())
		resp.Diagnostics.Append(self.client.DeleteIt(ctx, &icon)...)
		self.client.Mutex.Unlock(ctx, icon.LockKey())
	}
}

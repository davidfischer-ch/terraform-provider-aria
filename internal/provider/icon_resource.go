// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IconResource{}
var _ resource.ResourceWithImportState = &IconResource{}

func NewIconResource() resource.Resource {
	return &IconResource{}
}

// IconResource defines the resource implementation.
type IconResource struct {
	client *resty.Client
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

	response, err := self.client.R().
		// TODO SetQueryParam("apiVersion", ICON_API_VERSION).
		SetFileReader("file", "file", strings.NewReader(icon.Content.ValueString())).
		Post("icon/api/icons")

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

	iconId := icon.Id.ValueString()
	response, err := self.client.R().
		// TODO SetQueryParam("apiVersion", ICON_API_VERSION).
		Get("icon/api/icons/" + iconId)

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
	icon.Content = types.StringValue(response.String())
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
		fmt.Sprintf("Unable to update %s, method is not implement.", icon.String()))
}

func (self *IconResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var icon IconModel
	resp.Diagnostics.Append(req.State.Get(ctx, &icon)...)
	if resp.Diagnostics.HasError() {
		return
	}

	iconId := icon.Id.ValueString()
	if len(iconId) == 0 {
		return
	}

	resp.Diagnostics.Append(
		DeleteIt(
			self.client,
			ctx,
			icon.String(),
			"icon/api/icons/"+iconId,
			ICON_API_VERSION,
		)...,
	)
}

func (self *IconResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

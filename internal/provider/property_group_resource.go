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
var _ resource.Resource = &PropertyGroupResource{}
var _ resource.ResourceWithImportState = &PropertyGroupResource{}

func NewPropertyGroupResource() resource.Resource {
	return &PropertyGroupResource{}
}

// PropertyGroupResource defines the resource implementation.
type PropertyGroupResource struct {
	client *resty.Client
}

func (self *PropertyGroupResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_property_group"
}

func (self *PropertyGroupResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = PropertyGroupSchema()
}

func (self *PropertyGroupResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *PropertyGroupResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var propertyGroup PropertyGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &propertyGroup)...)
	if resp.Diagnostics.HasError() {
		return
	}

	propertyGroupRaw, diags := propertyGroup.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.R().
		SetQueryParam("apiVersion", BLUEPRINT_API_VERSION).
		SetBody(propertyGroupRaw).
		SetResult(&propertyGroupRaw).
		Post("properties/api/property-groups")
	err = handleAPIResponse(ctx, response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", propertyGroup.String(), err))
		return
	}

	// Save property group into Terraform state
	resp.Diagnostics.Append(propertyGroup.FromAPI(ctx, propertyGroupRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &propertyGroup)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", propertyGroup.String()))
}

func (self *PropertyGroupResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var propertyGroup PropertyGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &propertyGroup)...)
	if resp.Diagnostics.HasError() {
		return
	}

	propertyGroupId := propertyGroup.Id.ValueString()
	var propertyGroupRaw PropertyGroupAPIModel
	response, err := self.client.R().
		SetQueryParam("apiVersion", BLUEPRINT_API_VERSION).
		SetResult(&propertyGroupRaw).
		Get("properties/api/property-groups/" + propertyGroupId)

	// Handle gracefully a resource that has vanished on the platform
	// Beware that some APIs respond with HTTP 404 instead of 403 ...
	if response.StatusCode() == 404 {
		tflog.Debug(ctx, fmt.Sprintf("%s not found", propertyGroup.String()))
		resp.State.RemoveResource(ctx)
		return
	}

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read %s, got error: %s", propertyGroup.String(), err))
		return
	}

	// Save updated property group into Terraform state
	resp.Diagnostics.Append(propertyGroup.FromAPI(ctx, propertyGroupRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &propertyGroup)...)
}

func (self *PropertyGroupResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var propertyGroup PropertyGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &propertyGroup)...)
	if resp.Diagnostics.HasError() {
		return
	}

	propertyGroupId := propertyGroup.Id.ValueString()
	propertyGroupRaw, diags := propertyGroup.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.R().
		SetQueryParam("apiVersion", BLUEPRINT_API_VERSION).
		SetBody(propertyGroupRaw).
		SetResult(&propertyGroupRaw).
		Put("properties/api/property-groups/" + propertyGroupId)

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", propertyGroup.String(), err))
		return
	}

	// Save updated property group into Terraform state
	resp.Diagnostics.Append(propertyGroup.FromAPI(ctx, propertyGroupRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &propertyGroup)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", propertyGroup.String()))
}

func (self *PropertyGroupResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var propertyGroup PropertyGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &propertyGroup)...)
	if resp.Diagnostics.HasError() {
		return
	}

	propertyGroupId := propertyGroup.Id.ValueString()
	if len(propertyGroupId) == 0 {
		return
	}

	resp.Diagnostics.Append(
		DeleteIt(
			self.client,
			ctx,
			propertyGroup.String(),
			"properties/api/property-groups/"+propertyGroupId,
			BLUEPRINT_API_VERSION,
		)...,
	)
}

func (self *PropertyGroupResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// FIXME must be filtered by id and projectId
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

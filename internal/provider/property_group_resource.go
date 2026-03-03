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
var _ resource.Resource = &PropertyGroupResource{}
var _ resource.ResourceWithImportState = &PropertyGroupResource{}

func NewPropertyGroupResource() resource.Resource {
	return &PropertyGroupResource{}
}

// PropertyGroupResource defines the resource implementation.
type PropertyGroupResource struct {
	client *AriaClient
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

	propertyGroupToAPI, diags := propertyGroup.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var propertyGroupFromAPI PropertyGroupAPIModel
	_, createDiags := self.client.CreateIt(&propertyGroup, &propertyGroupFromAPI, propertyGroupToAPI)
	resp.Diagnostics.Append(createDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save property group into Terraform state
	resp.Diagnostics.Append(propertyGroup.FromAPI(ctx, propertyGroupFromAPI)...)
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

	var propertyGroupFromAPI PropertyGroupAPIModel
	found, _, readDiags := self.client.ReadIt(&propertyGroup, &propertyGroupFromAPI)
	resp.Diagnostics.Append(readDiags...)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated property group into Terraform state
	resp.Diagnostics.Append(propertyGroup.FromAPI(ctx, propertyGroupFromAPI)...)
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

	propertyGroupToAPI, diags := propertyGroup.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var propertyGroupFromAPI PropertyGroupAPIModel
	_, updateDiags := self.client.UpdateIt(&propertyGroup, &propertyGroupFromAPI, propertyGroupToAPI, "PUT")
	resp.Diagnostics.Append(updateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated property group into Terraform state
	resp.Diagnostics.Append(propertyGroup.FromAPI(ctx, propertyGroupFromAPI)...)
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

	resp.Diagnostics.Append(self.client.DeleteIt(&propertyGroup)...)
}

func (self *PropertyGroupResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// FIXME must be filtered by id and projectId
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

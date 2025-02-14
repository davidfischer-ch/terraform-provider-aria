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
var _ resource.Resource = &CatalogItemIconResource{}
var _ resource.ResourceWithImportState = &CatalogItemIconResource{}

func NewCatalogItemIconResource() resource.Resource {
	return &CatalogItemIconResource{}
}

// CatalogItemIconResource defines the resource implementation.
type CatalogItemIconResource struct {
	client *AriaClient
}

func (self *CatalogItemIconResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_catalog_item_icon"
}

func (self *CatalogItemIconResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = CatalogItemIconSchema()
}

func (self *CatalogItemIconResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *CatalogItemIconResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var itemIcon CatalogItemIconModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &itemIcon)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var itemIconFromAPI CatalogItemIconAPIModel
	path := itemIcon.CreatePath()
	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", GetVersionFromPath(path)).
		SetBody(itemIcon.ToAPI()).
		SetResult(&itemIconFromAPI).
		Patch(path)
	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", itemIcon.String(), err))
		return
	}

	// Save item's icon into Terraform state
	itemIcon.FromAPI(itemIconFromAPI)
	resp.Diagnostics.Append(resp.State.Set(ctx, &itemIcon)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", itemIcon.String()))
}

func (self *CatalogItemIconResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var itemIcon CatalogItemIconModel
	resp.Diagnostics.Append(req.State.Get(ctx, &itemIcon)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var itemIconFromAPI CatalogItemIconAPIModel
	found, _, readDiags := self.client.ReadIt(ctx, &itemIcon, &itemIconFromAPI)
	resp.Diagnostics.Append(readDiags...)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	if !resp.Diagnostics.HasError() {
		// Save updated item's icon into Terraform state
		itemIcon.FromAPI(itemIconFromAPI)
		resp.Diagnostics.Append(resp.State.Set(ctx, &itemIcon)...)
	}
}

func (self *CatalogItemIconResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var itemIcon CatalogItemIconModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &itemIcon)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var itemIconFromAPI CatalogItemIconAPIModel
	path := itemIcon.UpdatePath()
	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", GetVersionFromPath(path)).
		SetBody(itemIcon.ToAPI()).
		SetResult(&itemIconFromAPI).
		Patch(path)

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", itemIcon.String(), err))
		return
	}

	// Save updated item's icon into Terraform state
	itemIcon.FromAPI(itemIconFromAPI)
	resp.Diagnostics.Append(resp.State.Set(ctx, &itemIcon)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", itemIcon.String()))
}

func (self *CatalogItemIconResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Nothing to do.
}

func (self *CatalogItemIconResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("item_id"), req, resp)
}

// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &CatalogSourceResource{}

func NewCatalogSourceResource() resource.Resource {
	return &CatalogSourceResource{}
}

// CatalogSourceResource defines the resource implementation.
type CatalogSourceResource struct {
	client *AriaClient
}

func (self *CatalogSourceResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_catalog_source"
}

func (self *CatalogSourceResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = CatalogSourceSchema()
}

func (self *CatalogSourceResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *CatalogSourceResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var source CatalogSourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &source)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceRaw, diags := self.ManageIt(ctx, &source, "create")
	resp.Diagnostics.Append(diags...)

	// Save catalog source into Terraform state
	resp.Diagnostics.Append(source.FromAPI(ctx, sourceRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &source)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", source.String()))
}

func (self *CatalogSourceResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var source CatalogSourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &source)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var sourceRaw CatalogSourceAPIModel
	found, _, readDiags := self.client.ReadIt(ctx, &source, &sourceRaw)
	resp.Diagnostics.Append(readDiags...)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	if !resp.Diagnostics.HasError() {
		// Save updated catalog source into Terraform state
		resp.Diagnostics.Append(source.FromAPI(ctx, sourceRaw)...)
		resp.Diagnostics.Append(resp.State.Set(ctx, &source)...)
	}
}

func (self *CatalogSourceResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var source CatalogSourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &source)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceRaw, diags := self.ManageIt(ctx, &source, "update")
	resp.Diagnostics.Append(diags...)

	// Save catalog source into Terraform state
	resp.Diagnostics.Append(source.FromAPI(ctx, sourceRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &source)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", source.String()))
}

func (self *CatalogSourceResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var source CatalogSourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &source)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(self.client.DeleteIt(ctx, &source)...)
	}
}

// -------------------------------------------------------------------------------------------------

// Implement the magic behind the create and update methods.
func (self *CatalogSourceResource) ManageIt(
	ctx context.Context,
	source *CatalogSourceModel,
	method string,
) (CatalogSourceAPIModel, diag.Diagnostics) {

	var sourceRaw CatalogSourceAPIModel
	diags := diag.Diagnostics{}

	// Check method is valid
	if !slices.Contains([]string{"create", "update"}, method) {
		diags.AddError("Client error", fmt.Sprintf("BUG: Wrong method %s", method))
		return sourceRaw, diags
	}

	sourceRaw, someDiags := source.ToAPI(ctx)
	diags.Append(someDiags...)
	if diags.HasError() {
		return sourceRaw, diags
	}

	var path string
	if method == "create" {
		path = source.CreatePath()
	} else {
		path = source.UpdatePath()
	}

	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", FORM_API_VERSION).
		SetBody(sourceRaw).
		SetResult(&sourceRaw).
		Post(path)
	err = handleAPIResponse(ctx, response, err, []int{200, 201})
	if err != nil {
		diags.AddError(
			"Client error",
			fmt.Sprintf("Unable to %s %s, got error: %s", method, source.String(), err))
		return sourceRaw, diags
	}

	// Do not wait for catalog items to be imported
	diags.Append(source.FromAPI(ctx, sourceRaw)...)
	if diags.HasError() || !source.WaitImported.ValueBool() || !source.IsImporting(ctx) {
		return sourceRaw, diags
	}

	name := source.String()
	tflog.Debug(ctx, fmt.Sprintf("Wait %s to be imported...", name))

	// Poll for catalog items to be imported up to 10 minutes (60 x 10 seconds)
	maxAttempts := 60
	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Poll resource until imported
		time.Sleep(time.Duration(10 * time.Second))
		tflog.Debug(
			ctx,
			fmt.Sprintf("Poll %d of %d - Check %s is imported...", attempt+1, maxAttempts, name))

		found, _, someDiags := self.client.ReadIt(ctx, source, &sourceRaw)
		diags.Append(someDiags...)
		if !found {
			diags.AddError(
				"Client error",
				fmt.Sprintf("Resource %s has vanished while waiting to be imported.", name))
			return sourceRaw, diags
		}

		diags.Append(source.FromAPI(ctx, sourceRaw)...)
		if diags.HasError() || !source.IsImporting(ctx) {
			return sourceRaw, diags
		}

		// FIXME Ensure the catalog item is available before creating the catalog source?
		// FIXME Chicken & Egg Problem ...

		// Error downloading catalog item '/workflow/98789b67-8813-4fb1-b093-677d85d1017b'
		// (Error: Content provider error).
	}

	diags.AddError(
		"Client error",
		fmt.Sprintf("Timeout while waiting for %s to be imported.", name))
	return sourceRaw, diags
}

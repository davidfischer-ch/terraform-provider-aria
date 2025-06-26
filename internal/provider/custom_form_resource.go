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
var _ resource.Resource = &CustomFormResource{}
var _ resource.ResourceWithImportState = &CustomFormResource{}

func NewCustomFormResource() resource.Resource {
	return &CustomFormResource{}
}

// CustomFormResource defines the resource implementation.
type CustomFormResource struct {
	client *AriaClient
}

func (self *CustomFormResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_custom_form"
}

func (self *CustomFormResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = CustomFormSchema()
}

func (self *CustomFormResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *CustomFormResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var form CustomFormModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &form)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// First, try to fetch (existing form)
	var formFromFetchAPI CustomFormAPIModel
	path := form.FetchPath()
	response, err := self.client.R(path).
		SetQueryParam("formFormat", "JSON").
		SetQueryParam("formType", form.Type.ValueString()).
		SetQueryParam("sourceId", form.SourceId.ValueString()).
		SetQueryParam("sourceType", form.SourceType.ValueString()).
		SetResult(&formFromFetchAPI).
		Get(path)
	err = self.client.HandleAPIResponse(response, err, []int{200, 404})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to fetch %s, got error: %s", form.String(), err))
		return
	}

	// It may be missing, so generate the identifier in such as case...
	form.GenerateId(formFromFetchAPI.Id)

	// Then create (or update) it
	path = form.CreatePath()
	response, err = self.client.R(path).SetBody(form.ToAPI()).Post(path)
	err = self.client.HandleAPIResponse(response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", form.String(), err))
		return
	}

	// Read (using API) to retrieve the custom form content (and not empty stuff)
	var formFromAPI CustomFormAPIModel
	path = form.ReadPath()
	response, err = self.client.R(path).SetResult(&formFromAPI).Get(path)
	err = self.client.HandleAPIResponse(response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read %s, got error: %s", form.String(), err))
		return
	}

	// Save custom form into Terraform state
	form.FromAPI(formFromAPI)
	resp.Diagnostics.Append(resp.State.Set(ctx, &form)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", form.String()))
}

func (self *CustomFormResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var form CustomFormModel
	resp.Diagnostics.Append(req.State.Get(ctx, &form)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var formFromAPI CustomFormAPIModel
	found, _, readDiags := self.client.ReadIt(&form, &formFromAPI)
	resp.Diagnostics.Append(readDiags...)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated custom form into Terraform state
	form.FromAPI(formFromAPI)
	resp.Diagnostics.Append(resp.State.Set(ctx, &form)...)
}

func (self *CustomFormResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var form CustomFormModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &form)...)
	if resp.Diagnostics.HasError() {
		return
	}

	path := form.UpdatePath()
	response, err := self.client.R(path).SetBody(form.ToAPI()).Post(path)
	err = self.client.HandleAPIResponse(response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", form.String(), err))
		return
	}

	// Read (using API) to retrieve the custom form content (and not empty stuff)
	var formFromAPI CustomFormAPIModel
	path = form.ReadPath()
	response, err = self.client.R(path).SetResult(&formFromAPI).Get(path)
	err = self.client.HandleAPIResponse(response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read %s, got error: %s", form.String(), err))
		return
	}

	// Save custom form into Terraform state
	form.FromAPI(formFromAPI)
	resp.Diagnostics.Append(resp.State.Set(ctx, &form)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", form.String()))
}

func (self *CustomFormResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var form CustomFormModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &form)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(self.client.DeleteIt(&form)...)
	}
}

func (self *CustomFormResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

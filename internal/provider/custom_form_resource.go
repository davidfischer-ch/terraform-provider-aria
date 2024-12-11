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

	form.GenerateId()
	formRaw, diags := form.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", FORM_API_VERSION).
		SetBody(formRaw).
		Post(form.CreatePath())
	err = handleAPIResponse(ctx, response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", form.String(), err))
		return
	}

	// Save custom form into Terraform state
	resp.Diagnostics.Append(form.FromAPI(ctx, formRaw)...)
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

	var formRaw CustomFormAPIModel
	found, _, readDiags := self.client.ReadIt(ctx, &form, &formRaw)
	resp.Diagnostics.Append(readDiags...)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	if !resp.Diagnostics.HasError() {
		// Save updated custom form into Terraform state
		resp.Diagnostics.Append(form.FromAPI(ctx, formRaw)...)
		resp.Diagnostics.Append(resp.State.Set(ctx, &form)...)
	}
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

	formRaw, diags := form.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", FORM_API_VERSION).
		SetBody(formRaw).
		Post(form.UpdatePath())
	err = handleAPIResponse(ctx, response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", form.String(), err))
		return
	}

	// Save custom form into Terraform state
	resp.Diagnostics.Append(form.FromAPI(ctx, formRaw)...)
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
		resp.Diagnostics.Append(self.client.DeleteIt(ctx, &form)...)
	}
}

func (self *CustomFormResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

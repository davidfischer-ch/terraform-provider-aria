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
var _ resource.Resource = &CloudTemplateV1Resource{}
var _ resource.ResourceWithImportState = &CloudTemplateV1Resource{}

func NewCloudTemplateV1Resource() resource.Resource {
	return &CloudTemplateV1Resource{}
}

// CloudTemplateV1Resource defines the resource implementation.
type CloudTemplateV1Resource struct {
	client *AriaClient
}

func (self *CloudTemplateV1Resource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_cloud_template_v1"
}

func (self *CloudTemplateV1Resource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = CloudTemplateV1Schema()
}

func (self *CloudTemplateV1Resource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *CloudTemplateV1Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var template CloudTemplateV1Model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &template)...)
	if resp.Diagnostics.HasError() {
		return
	}

	templateRaw, diags := template.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", BLUEPRINT_API_VERSION).
		SetBody(templateRaw).
		SetResult(&templateRaw).
		Post(template.CreatePath())
	err = handleAPIResponse(ctx, response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", template.String(), err))
		return
	}

	// Save cloud template into Terraform state
	resp.Diagnostics.Append(template.FromAPI(ctx, templateRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &template)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", template.String()))
}

func (self *CloudTemplateV1Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var template CloudTemplateV1Model
	resp.Diagnostics.Append(req.State.Get(ctx, &template)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var templateRaw CloudTemplateV1APIModel
	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", BLUEPRINT_API_VERSION).
		SetResult(&templateRaw).
		Get(template.ReadPath())

	// Handle gracefully a resource that has vanished on the platform
	// Beware that some APIs respond with HTTP 404 instead of 403 ...
	if response.StatusCode() == 404 {
		tflog.Debug(ctx, fmt.Sprintf("%s not found", template.String()))
		resp.State.RemoveResource(ctx)
		return
	}

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read %s, got error: %s", template.String(), err))
		return
	}

	// Save updated cloud template into Terraform state
	resp.Diagnostics.Append(template.FromAPI(ctx, templateRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &template)...)
}

func (self *CloudTemplateV1Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var template CloudTemplateV1Model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &template)...)
	if resp.Diagnostics.HasError() {
		return
	}

	templateRaw, diags := template.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", BLUEPRINT_API_VERSION).
		SetBody(templateRaw).
		SetResult(&templateRaw).
		Put(template.UpdatePath())

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", template.String(), err))
		return
	}

	// Save updated cloud template into Terraform state
	resp.Diagnostics.Append(template.FromAPI(ctx, templateRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &template)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", template.String()))
}

func (self *CloudTemplateV1Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var template CloudTemplateV1Model
	resp.Diagnostics.Append(req.State.Get(ctx, &template)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(self.client.DeleteIt(ctx, &template)...)
	}
}

func (self *CloudTemplateV1Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// FIXME must be filtered by id and projectId
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

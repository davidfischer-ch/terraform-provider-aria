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
var _ resource.Resource = &TagResource{}
var _ resource.ResourceWithImportState = &TagResource{}

func NewTagResource() resource.Resource {
	return &TagResource{}
}

// TagResource defines the resource implementation.
type TagResource struct {
	client *AriaClient
}

func (self *TagResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_tag"
}

func (self *TagResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = TagSchema()
}

func (self *TagResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *TagResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var tag TagModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &tag)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var tagRaw TagAPIModel
	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", IAAS_API_VERSION).
		SetBody(tag.ToAPI()).
		SetResult(&tagRaw).
		Post(tag.CreatePath())

	err = handleAPIResponse(ctx, response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", tag.String(), err))
		return
	}

	// Save tag into Terraform state
	tag.FromAPI(tagRaw)
	resp.Diagnostics.Append(resp.State.Set(ctx, &tag)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", tag.String()))
}

func (self *TagResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var tag TagModel
	resp.Diagnostics.Append(req.State.Get(ctx, &tag)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO Read by filtering tag list by ID
	var listRaw TagListAPIModel
	listPath := tag.ListPath()
	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", GetVersionFromPath(listPath)).
		SetQueryParam("$filter", fmt.Sprintf("id eq %s", tag.Id.ValueString())).
		SetQueryParam("$top", "2"). // Make it possible to know if filter works properly
		SetResult(&listRaw).
		Get(listPath)
	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to list tags to get %s, got error: %s", tag.String(), err))
		return
	}

	// Do not rely on NumberOfElements
	if len(listRaw.Content) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}

	// Neither TotalElements
	if len(listRaw.Content) > 1 {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf(
				"Expected one and only one tag matching %s ID, found: %d",
				tag.String(), listRaw.TotalElements,
			),
		)
		return
	}

	// Save updated tag into Terraform state
	tag.FromAPI(listRaw.Content[0])
	resp.Diagnostics.Append(resp.State.Set(ctx, &tag)...)
}

func (self *TagResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var tag TagModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &tag)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated tag into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &tag)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", tag.String()))
}

func (self *TagResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var tag TagModel
	resp.Diagnostics.Append(req.State.Get(ctx, &tag)...)
	if !resp.Diagnostics.HasError() && !tag.KeepOnDestroy.ValueBool() {
		resp.Diagnostics.Append(self.client.DeleteIt(ctx, &tag)...)
	}
}

func (self *TagResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("force_delete"), false)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("keep_on_destroy"), false)...)
}

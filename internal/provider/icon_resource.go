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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
	resp.Schema = schema.Schema{
		MarkdownDescription: "Icon resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Icon identifier (Aria seem to compute it from content)",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"content": schema.StringAttribute{
				MarkdownDescription: "Icon content (force recreation on change)",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
		},
	}
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
		SetFileReader("file", "file", strings.NewReader(icon.Content.ValueString())).
		Post("icon/api/icons")

	err = handleAPIResponse(ctx, response, err, 201)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create icon, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("icon created %s", response.Body()))
	iconId, err := GetIdFromLocation(response)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to parse icon id, got error: %s", err))
		return
	}

	// Save icon into Terraform state
	icon.Id = types.StringValue(iconId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &icon)...)
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
	response, err := self.client.R().Get("icon/api/icons/" + iconId)

	// Handle gracefully a resource that has vanished on the platform
	// Beware that some APIs respond with HTTP 404 instead of 403 ...
	if response.StatusCode() == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	err = handleAPIResponse(ctx, response, err, 200)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read icon %s, got error: %s", iconId, err))
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
		fmt.Sprintf("Unable to update icon %s, method is not implement.", icon.Id.ValueString()))
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

	response, err := self.client.R().Delete("icon/api/icons/" + iconId)

	err = handleAPIResponse(ctx, response, err, 204)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to delete icon %s, got error: %s", iconId, err))
	}

	tflog.Debug(ctx, fmt.Sprintf("Icon %s deleted", iconId))
}

func (self *IconResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

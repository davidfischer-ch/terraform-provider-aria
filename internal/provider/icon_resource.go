// Copyright (c) HashiCorp, Inc.
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

// IconResourceModel describes the resource data model.
type IconResourceModel struct {
	Id      types.String `tfsdk:"id"`
	Content types.String `tfsdk:"content"`
}

func (r *IconResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_icon"
}

func (r *IconResource) Schema(
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

func (r *IconResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = GetResourceClient(ctx, req, resp)
}

func (r *IconResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var icon IconResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &icon)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.R().
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
	icon.Id = types.StringValue(iconId)

	// Save icon into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &icon)...)
}

func (r *IconResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var icon IconResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &icon)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.R().Get("icon/api/icons/" + icon.Id.ValueString())

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
			fmt.Sprintf("Unable to read icon %s, got error: %s", icon.Id.ValueString(), err))
		return
	}

	icon.Content = types.StringValue(response.String())

	// Save updated icon into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &icon)...)
}

func (r *IconResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var icon IconResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &icon)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddError(
		"Client error",
		fmt.Sprintf("Unable to update icon %s, method is not implement.", icon.Id.ValueString()))
}

func (r *IconResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var icon IconResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &icon)...)
	if resp.Diagnostics.HasError() {
		return
	}

	iconId := icon.Id.ValueString()
	if len(iconId) == 0 {
		return
	}

	response, err := r.client.R().Delete("icon/api/icons/" + icon.Id.ValueString())
	err = handleAPIResponse(ctx, response, err, 204)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to delete icon %s, got error: %s", icon.Id.ValueString(), err))
	}
}

func (r *IconResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

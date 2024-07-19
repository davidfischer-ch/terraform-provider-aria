// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ABXConstantResource{}
var _ resource.ResourceWithImportState = &ABXConstantResource{}

func NewABXConstantResource() resource.Resource {
	return &ABXConstantResource{}
}

// ABXConstantResource defines the resource implementation.
type ABXConstantResource struct {
	client *resty.Client
}

func (self *ABXConstantResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_abx_constant"
}

func (self *ABXConstantResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "ABX constant resource",
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "Name",
				Required:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "Value",
				Required:            true,
			},
			"encrypted": schema.BoolAttribute{
				MarkdownDescription: "Should be always unencrypted!",
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"org_id": ComputedOrganizationIdSchema(),
		},
	}
}

func (self *ABXConstantResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *ABXConstantResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var constant ABXConstantModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &constant)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var constantRaw ABXConstantAPIModel
	response, err := self.client.R().
		SetQueryParam("apiVersion", ABX_API_VERSION).
		SetBody(constant.ToAPI()).
		SetResult(&constantRaw).
		Post("abx/api/resources/action-secrets")

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", constant.String(), err))
		return
	}

	// Save constant into Terraform state
	resp.Diagnostics.Append(constant.FromAPI(ctx, constantRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &constant)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", constant.String()))
}

func (self *ABXConstantResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var constant ABXConstantModel
	resp.Diagnostics.Append(req.State.Get(ctx, &constant)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var constantRaw ABXConstantAPIModel
	constantId := constant.Id.ValueString()
	response, err := self.client.R().
		SetQueryParam("apiVersion", ABX_API_VERSION).
		SetResult(&constantRaw).
		Get("abx/api/resources/action-secrets/" + constantId)

	// Handle gracefully a resource that has vanished on the platform
	// Beware that some APIs respond with HTTP 404 instead of 403 ...
	if response.StatusCode() == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read %s, got error: %s", constant.String(), err))
		return
	}

	// Save updated constant into Terraform state
	resp.Diagnostics.Append(constant.FromAPI(ctx, constantRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &constant)...)
}

func (self *ABXConstantResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var constant ABXConstantModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &constant)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var constantRaw ABXConstantAPIModel
	constantId := constant.Id.ValueString()
	response, err := self.client.R().
		SetQueryParam("apiVersion", ABX_API_VERSION).
		SetBody(constant.ToAPI()).
		SetResult(&constantRaw).
		Put("abx/api/resources/action-secrets/" + constantId)

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", constant.String(), err))
		return
	}

	// Save constant into Terraform state
	resp.Diagnostics.Append(constant.FromAPI(ctx, constantRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &constant)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", constant.String()))
}

func (self *ABXConstantResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var constant ABXConstantModel
	resp.Diagnostics.Append(req.State.Get(ctx, &constant)...)
	if resp.Diagnostics.HasError() {
		return
	}

	constantId := constant.Id.ValueString()
	if len(constantId) == 0 {
		return
	}

	resp.Diagnostics.Append(
		DeleteIt(
			self.client,
			ctx,
			constant.String(),
			"abx/api/resources/action-secrets/"+constantId,
			ABX_API_VERSION,
		)...,
	)
}

func (self *ABXConstantResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

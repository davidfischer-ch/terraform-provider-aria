// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ABXSensitiveConstantResource{}

func NewABXSensitiveConstantResource() resource.Resource {
	return &ABXSensitiveConstantResource{}
}

// ABXSensitiveConstantResource defines the resource implementation.
type ABXSensitiveConstantResource struct {
	client *resty.Client
}

func (self *ABXSensitiveConstantResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_abx_sensitive_constant"
}

func (self *ABXSensitiveConstantResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "ABX sensitive constant resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Constant identifier",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Constant name",
				Required:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "Constant value (cannot be enforced since API don't return it)",
				Required:            true,
				Sensitive:           true,
			},
			"encrypted": schema.BoolAttribute{
				MarkdownDescription: "Constant should be always encrypted!",
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"org_id": schema.StringAttribute{
				MarkdownDescription: "Constant organisation identifier",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (self *ABXSensitiveConstantResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *ABXSensitiveConstantResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var secret ABXSensitiveConstantModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &secret)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var secretRaw ABXSensitiveConstantAPIModel
	response, err := self.client.R().
		SetBody(secret.ToAPI()).
		SetResult(&secretRaw).
		Post("abx/api/resources/action-secrets")

	err = handleAPIResponse(ctx, response, err, 200)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create ABX sensitive constant, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("ABX sensitive constant %s created", secretRaw.Id))

	// Save secret into Terraform state
	resp.Diagnostics.Append(secret.FromAPI(ctx, secretRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &secret)...)
}

func (self *ABXSensitiveConstantResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var secret ABXSensitiveConstantModel
	resp.Diagnostics.Append(req.State.Get(ctx, &secret)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var secretRaw ABXSensitiveConstantAPIModel
	secretId := secret.Id.ValueString()
	response, err := self.client.R().
		SetResult(&secretRaw).
		Get("abx/api/resources/action-secrets/" + secretId)

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
			fmt.Sprintf("Unable to read ABX sensitive constant %s, got error: %s", secretId, err))
		return
	}

	// Save updated secret into Terraform state
	resp.Diagnostics.Append(secret.FromAPI(ctx, secretRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &secret)...)
}

func (self *ABXSensitiveConstantResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var secret ABXSensitiveConstantModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &secret)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var secretRaw ABXSensitiveConstantAPIModel
	secretId := secret.Id.ValueString()
	response, err := self.client.R().
		SetBody(secret.ToAPI()).
		SetResult(&secretRaw).
		Put("abx/api/resources/action-secrets/" + secretId)

	err = handleAPIResponse(ctx, response, err, 200)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update ABX sensitive constant %s, got error: %s", secretId, err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Secret %s updated", secretId))

	// Save secret into Terraform state
	resp.Diagnostics.Append(secret.FromAPI(ctx, secretRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &secret)...)
}

func (self *ABXSensitiveConstantResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var secret ABXSensitiveConstantModel
	resp.Diagnostics.Append(req.State.Get(ctx, &secret)...)
	if resp.Diagnostics.HasError() {
		return
	}

	secretId := secret.Id.ValueString()
	if len(secretId) == 0 {
		return
	}

	response, err := self.client.R().Delete("abx/api/resources/action-secrets/" + secretId)

	err = handleAPIResponse(ctx, response, err, 200)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to delete ABX sensitive constant %s, got error: %s", secretId, err))
	}

	tflog.Debug(ctx, fmt.Sprintf("ABX sensitive constant %s deleted", secretId))
}

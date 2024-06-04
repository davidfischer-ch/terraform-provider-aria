// Copyright (c) HashiCorp, Inc.
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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ABXConstantResource{}
var _ resource.ResourceWithImportState = &IconResource{}

func NewABXConstantResource() resource.Resource {
	return &ABXConstantResource{}
}

// ABXConstantResource defines the resource implementation.
type ABXConstantResource struct {
	client *resty.Client
}

// ABXConstantResourceModel describes the resource data model.
type ABXConstantResourceModel struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Value     types.String `tfsdk:"value"`
	Encrypted types.Bool   `tfsdk:"encrypted"`
	OrgId     types.String `tfsdk:"org_id"`
}

// ABXConstantResourceAPIModel describes the resource API model.
type ABXConstantResourceAPIModel struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Value         string `json:"value"`
	Encrypted     bool   `json:"encrypted"`
	OrgId         string `json:"orgId"`
	CreatedMillis int64  `json:"createdMillis"`
}

func (r *ABXConstantResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_abx_constant"
}

func (r *ABXConstantResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "ABX constant resource",
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
				MarkdownDescription: "Constant value",
				Required:            true,
				Sensitive:           true,
			},
			"encrypted": schema.BoolAttribute{
				MarkdownDescription: "Constant should be always unencrypted!",
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"org_id": schema.StringAttribute{
				MarkdownDescription: "Constant organisation ID",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (r *ABXConstantResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = GetResourceClient(ctx, req, resp)
}

func (r *ABXConstantResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var constant ABXConstantResourceModel
	var constantRaw ABXConstantResourceAPIModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &constant)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.R().
		SetBody(ABXConstantResourceAPIModel{
			Name:      constant.Name.ValueString(),
			Value:     constant.Value.ValueString(),
			Encrypted: constant.Encrypted.ValueBool(),
		}).
		SetResult(&constantRaw).
		Post("abx/api/resources/action-secrets")

	err = handleAPIResponse(ctx, response, err, 200)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create ABX constant, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("ABX constant %s created", constantRaw.Id))

	constant.Id = types.StringValue(constantRaw.Id)
	constant.Name = types.StringValue(constantRaw.Name)
	constant.Value = types.StringValue(constantRaw.Value)
	constant.Encrypted = types.BoolValue(constantRaw.Encrypted)
	constant.OrgId = types.StringValue(constantRaw.OrgId)

	// Save constant into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &constant)...)
}

func (r *ABXConstantResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var constant ABXConstantResourceModel
	var constantRaw ABXConstantResourceAPIModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &constant)...)
	if resp.Diagnostics.HasError() {
		return
	}

	constantId := constant.Id.ValueString()
	response, err := r.client.R().
		SetResult(&constantRaw).
		Get("abx/api/resources/action-secrets/" + constantId)

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
			fmt.Sprintf("Unable to read ABX constant %s, got error: %s", constantId, err))
		return
	}

	constant.Name = types.StringValue(constantRaw.Name)
	constant.Value = types.StringValue(constantRaw.Value)
	constant.Encrypted = types.BoolValue(constantRaw.Encrypted)
	constant.OrgId = types.StringValue(constantRaw.OrgId)

	// Save updated constant into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &constant)...)
}

func (r *ABXConstantResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var constant ABXConstantResourceModel
	var constantRaw ABXConstantResourceAPIModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &constant)...)
	if resp.Diagnostics.HasError() {
		return
	}

	constantId := constant.Id.ValueString()
	response, err := r.client.R().
		SetBody(ABXConstantResourceAPIModel{
			Name:      constant.Name.ValueString(),
			Value:     constant.Value.ValueString(),
			Encrypted: constant.Encrypted.ValueBool(),
		}).
		SetResult(&constantRaw).
		Put("abx/api/resources/action-secrets/" + constantId)

	err = handleAPIResponse(ctx, response, err, 200)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update ABX constant %s, got error: %s", constantId, err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("ABX constant %s updated", constantId))

	constant.Name = types.StringValue(constantRaw.Name)
	constant.Value = types.StringValue(constantRaw.Value)
	constant.Encrypted = types.BoolValue(constantRaw.Encrypted)
	constant.OrgId = types.StringValue(constantRaw.OrgId)

	// Save constant into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &constant)...)
}

func (r *ABXConstantResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var constant ABXConstantResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &constant)...)
	if resp.Diagnostics.HasError() {
		return
	}

	constantId := constant.Id.ValueString()
	if len(constantId) == 0 {
		return
	}

	response, err := r.client.R().Delete("abx/api/resources/action-secrets/" + constantId)
	err = handleAPIResponse(ctx, response, err, 200)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to delete ABX constant %s, got error: %s", constantId, err))
	}

	tflog.Debug(ctx, fmt.Sprintf("ABX constant %s deleted", constantId))
}

func (r *ABXConstantResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

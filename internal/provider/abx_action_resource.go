// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ABXActionResource{}
var _ resource.ResourceWithImportState = &ABXActionResource{}

func NewABXActionResource() resource.Resource {
	return &ABXActionResource{}
}

// ABXActionResource defines the resource implementation.
type ABXActionResource struct {
	client *resty.Client
}

func (self *ABXActionResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_abx_action"
}

func (self *ABXActionResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "ABX action resource",
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "A name (must be unique)",
				Required:            true,
			},
			"description": RequiredDescriptionSchema(),
			"faas_provider": schema.StringAttribute{
				MarkdownDescription: "FaaS provider used for code execution, one of auto (default), on-prem, aws and azure (automatically set by the platform if unset)",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("auto"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"auto", "on-prem", "aws", "azure"}...),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of action, one of SCRIPT (default), REST_CALL, REST_POLL, FLOW, VAULT and CYBERARK",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("SCRIPT"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"SCRIPT", "REST_CALL", "REST_POLL", "FLOW", "VAULT", "CYBERARK"}...),
				},
			},
			"runtime_name": schema.StringAttribute{
				MarkdownDescription: "Runtime name (python, nodejs, ...)",
				Required:            true,
			},
			"runtime_version": schema.StringAttribute{
				MarkdownDescription: "Runtime version (3.10, ...)",
				Computed:            true,
				Optional:            true,
			},
			"cpu_shares": schema.Int32Attribute{
				MarkdownDescription: "Runtime CPU shares",
				Computed:            true,
				Optional:            true,
				Default:             int32default.StaticInt32(1024),
			},
			"memory_in_mb": schema.Int32Attribute{
				MarkdownDescription: "Runtime memory constraint in MB",
				Required:            true,
			},
			"timeout_seconds": schema.Int32Attribute{
				MarkdownDescription: "How long an action can run (default to 600)",
				Computed:            true,
				Optional:            true,
				Default:             int32default.StaticInt32(600),
			},
			"deployment_timeout_seconds": schema.Int32Attribute{
				MarkdownDescription: "How long ??",
				Computed:            true,
				Optional:            true,
				Default:             int32default.StaticInt32(900),
			},
			"entrypoint": schema.StringAttribute{
				MarkdownDescription: "Main function's name",
				Required:            true,
			},
			"dependencies": schema.ListAttribute{
				MarkdownDescription: "Dependencies (python packages, ...)",
				ElementType:         types.StringType,
				Required:            true,
			},
			"constants": schema.SetAttribute{
				MarkdownDescription: "ABX Constants to expose to the action",
				ElementType:         types.StringType,
				Required:            true,
			},
			"secrets": schema.SetAttribute{
				MarkdownDescription: "Secrets to expose to the action",
				ElementType:         types.StringType,
				Required:            true,
			},
			/* "inputs": schema.MapAttribute{
				// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/dynamic-data
				MarkdownDescription: "Inputs to expose to the action",
				ElementType: types.DynamicType,
				Required: true,
			}, */
			"source": schema.StringAttribute{
				MarkdownDescription: "Action source code",
				Required:            true,
			},
			"project_id": OptionalImmutableProjectIdSchema(),
			"shared": schema.BoolAttribute{
				MarkdownDescription: "Flag indicating if the action can be shared across projects",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"system": schema.BoolAttribute{
				MarkdownDescription: "Flag indicating if the action is a system action",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"async_deployed": schema.BoolAttribute{
				MarkdownDescription: "TODO",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"org_id": ComputedOrganizationIdSchema(),
			/* created_millis int64 */
			/* updated_millis int64 */
			/* metadata {} */
			/* configuration {} */
			/* "self_link": schema.StringAttribute{
				MarkdownDescription: "URL to the action",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			}, */
		},
	}
}

func (self *ABXActionResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *ABXActionResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var action ABXActionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &action)...)
	if resp.Diagnostics.HasError() {
		return
	}

	actionRaw, diags := action.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.R().
		SetQueryParam("apiVersion", ABX_API_VERSION).
		SetBody(actionRaw).
		SetResult(&actionRaw).
		Post("abx/api/resources/actions")
	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", action.String(), err))
		return
	}

	// Save action into Terraform state
	resp.Diagnostics.Append(action.FromAPI(ctx, actionRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &action)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", action.String()))
}

func (self *ABXActionResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var action ABXActionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &action)...)
	if resp.Diagnostics.HasError() {
		return
	}

	actionId := action.Id.ValueString()
	projectId := action.ProjectId.ValueString()
	var actionRaw ABXActionAPIModel
	response, err := self.client.R().
		SetQueryParam("apiVersion", ABX_API_VERSION).
		SetResult(&actionRaw).
		Get(fmt.Sprintf("abx/api/resources/actions/%s?projectId=%s", actionId, projectId))

	// Handle gracefully a resource that has vanished on the platform
	// Beware that some APIs respond with HTTP 404 instead of 403 ...
	if response.StatusCode() == 404 {
		tflog.Debug(ctx, fmt.Sprintf("%s not found", action.String()))
		resp.State.RemoveResource(ctx)
		return
	}

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read %s, got error: %s", action.String(), err))
		return
	}

	// Save updated action into Terraform state
	resp.Diagnostics.Append(action.FromAPI(ctx, actionRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &action)...)
}

func (self *ABXActionResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var action ABXActionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &action)...)
	if resp.Diagnostics.HasError() {
		return
	}

	actionId := action.Id.ValueString()
	projectId := action.ProjectId.ValueString()
	actionRaw, diags := action.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.R().
		SetQueryParam("apiVersion", ABX_API_VERSION).
		SetBody(actionRaw).
		SetResult(&actionRaw).
		Put(fmt.Sprintf("abx/api/resources/actions/%s?projectId=%s", actionId, projectId))

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", action.String(), err))
		return
	}

	// Save updated action into Terraform state
	resp.Diagnostics.Append(action.FromAPI(ctx, actionRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &action)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", action.String()))
}

func (self *ABXActionResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var action ABXActionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &action)...)
	if resp.Diagnostics.HasError() {
		return
	}

	actionId := action.Id.ValueString()
	projectId := action.ProjectId.ValueString()
	if len(actionId) == 0 {
		return
	}

	resp.Diagnostics.Append(
		DeleteIt(
			self.client,
			ctx,
			action.String(),
			fmt.Sprintf("abx/api/resources/actions/%s?projectId=%s", actionId, projectId),
			ABX_API_VERSION,
		)...,
	)
}

func (self *ABXActionResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// FIXME must be filtered by id and projectId
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

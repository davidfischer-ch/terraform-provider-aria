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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ResourceActionResource{}
var _ resource.ResourceWithImportState = &ResourceActionResource{}

func NewResourceActionResource() resource.Resource {
	return &ResourceActionResource{}
}

// ResourceActionResource defines the resource implementation.
type ResourceActionResource struct {
	client *resty.Client
}

func (self *ResourceActionResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_resource_action"
}

func (self *ResourceActionResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Native resource's action resource",
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "Action name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Action display name",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description",
				Required:            true,
			},
			"provider_name": schema.StringAttribute{
				MarkdownDescription: "Provider name, one of xaas (and that's all, maybe)",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("xaas"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_type": schema.StringAttribute{
				MarkdownDescription: "Native resource type",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			//"criteria": ResourceActionCriteriaSchema(),
			"form_definition": CustomFormSchema(),
			"runnable_item":   ResourceActionRunnableSchema("Action's runnable"),
			"status": schema.StringAttribute{
				MarkdownDescription: "Action status, either DRAFT or RELEASED",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("RELEASED"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"DRAFT", "RELEASED"}...),
				},
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "Project ID. Set it to \"\" to make the action " +
					"available for all projects.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"org_id": ComputedOrganizationIdSchema(),
		},
	}
}

func (self *ResourceActionResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *ResourceActionResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var action ResourceActionModel
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
		SetQueryParam("apiVersion", FORM_API_VERSION).
		SetBody(actionRaw).
		SetResult(&actionRaw).
		Post("form-service/api/custom/resource-actions")
	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", action.String(), err))
		return
	}

	// Save additional action into Terraform state
	resp.Diagnostics.Append(action.FromAPI(ctx, actionRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &action)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", action.String()))
}

func (self *ResourceActionResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var action ResourceActionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &action)...)
	if resp.Diagnostics.HasError() {
		return
	}

	actionId := action.Id.ValueString()
	var actionRaw ResourceActionAPIModel
	response, err := self.client.R().
		SetQueryParam("apiVersion", FORM_API_VERSION).
		SetResult(&actionRaw).
		Get("form-service/api/custom/resource-actions/" + actionId)

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

	// Save updated additional action into Terraform state
	resp.Diagnostics.Append(action.FromAPI(ctx, actionRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &action)...)
}

func (self *ResourceActionResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var action ResourceActionModel
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
		SetQueryParam("apiVersion", FORM_API_VERSION).
		SetBody(actionRaw).
		SetResult(&actionRaw).
		Post("form-service/api/custom/resource-actions") // Its not a mistake...

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", action.String(), err))
		return
	}

	// Save updated additional action into Terraform state
	resp.Diagnostics.Append(action.FromAPI(ctx, actionRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &action)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", action.String()))
}

func (self *ResourceActionResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var action ResourceActionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &action)...)
	if resp.Diagnostics.HasError() {
		return
	}

	actionId := action.Id.ValueString()
	if len(actionId) == 0 {
		return
	}

	resp.Diagnostics.Append(
		DeleteIt(
			self.client,
			ctx,
			action.String(),
			"form-service/api/custom/resource-actions/"+actionId,
			FORM_API_VERSION,
		)...,
	)
}

func (self *ResourceActionResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// FIXME must be filtered by id and projectId
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

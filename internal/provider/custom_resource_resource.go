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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &CustomResourceResource{}
var _ resource.ResourceWithImportState = &CustomResourceResource{}

func NewCustomResourceResource() resource.Resource {
	return &CustomResourceResource{}
}

// CustomResourceResource defines the resource implementation.
type CustomResourceResource struct {
	client *resty.Client
}

func (self *CustomResourceResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_custom_resource"
}

func (self *CustomResourceResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Custom Resource resource",
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"display_name": schema.StringAttribute{
				MarkdownDescription: "A friendly name",
				Required:            true,
			},
			"description": RequiredDescriptionSchema(),
			"resource_type": schema.StringAttribute{
				MarkdownDescription: "Define the type (must be unique, e.g. Custom.DB.PostgreSQL)",
				Required:            true,
			},
			"schema_type": schema.StringAttribute{
				MarkdownDescription: "Type of resource, one of ABX_USER_DEFINED (and that's all, maybe)",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("ABX_USER_DEFINED"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"ABX_USER_DEFINED"}...),
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Resource status, one of DRAFT, ON, or RELEASED",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("RELEASED"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"DRAFT", "ON", "RELEASED"}...),
				},
			},
			"properties": UnorderedPropertiesSchema("Resource's properties"),
			/* "allocate:" TODO one of the optional main actions */
			"create":     ResourceActionRunnableSchema("Create action"),
			"read":       ResourceActionRunnableSchema("Read action"),
			"update":     ResourceActionRunnableSchema("Update action"),
			"delete":     ResourceActionRunnableSchema("Delete action"),
			"project_id": OptionalImmutableProjectIdSchema(),
			"org_id":     ComputedOrganizationIdSchema(),
		},
	}
}

func (self *CustomResourceResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *CustomResourceResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var resource CustomResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resource)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceRaw, diags := resource.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.R().
		SetQueryParam("apiVersion", FORM_API_VERSION).
		SetBody(resourceRaw).
		SetResult(&resourceRaw).
		Post("form-service/api/custom/resource-types")
	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", resource.String(), err))
		return
	}

	// Save custom resource into Terraform state
	resp.Diagnostics.Append(resource.FromAPI(ctx, resourceRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &resource)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", resource.String()))
}

func (self *CustomResourceResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var resource CustomResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &resource)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceId := resource.Id.ValueString()
	var resourceRaw CustomResourceAPIModel
	response, err := self.client.R().
		SetQueryParam("apiVersion", FORM_API_VERSION).
		SetResult(&resourceRaw).
		Get("form-service/api/custom/resource-types/" + resourceId)

	// Handle gracefully a resource that has vanished on the platform
	// Beware that some APIs respond with HTTP 404 instead of 403 ...
	if response.StatusCode() == 404 {
		tflog.Debug(ctx, fmt.Sprintf("%s not found", resource.String()))
		resp.State.RemoveResource(ctx)
		return
	}

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read %s, got error: %s", resource.String(), err))
		return
	}

	// Save updated custom resource into Terraform state
	resp.Diagnostics.Append(resource.FromAPI(ctx, resourceRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &resource)...)
}

func (self *CustomResourceResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var resource CustomResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resource)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceRaw, diags := resource.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.R().
		SetQueryParam("apiVersion", FORM_API_VERSION).
		SetBody(resourceRaw).
		SetResult(&resourceRaw).
		Post("form-service/api/custom/resource-types") // Its not a mistake...

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", resource.String(), err))
		return
	}

	// Save updated custom resource into Terraform state
	resp.Diagnostics.Append(resource.FromAPI(ctx, resourceRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &resource)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", resource.String()))
}

func (self *CustomResourceResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var resource CustomResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &resource)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceId := resource.Id.ValueString()
	if len(resourceId) == 0 {
		return
	}

	resp.Diagnostics.Append(
		DeleteIt(
			self.client,
			ctx,
			resource.String(),
			"form-service/api/custom/resource-types/"+resourceId,
			FORM_API_VERSION,
		)...,
	)
}

func (self *CustomResourceResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// FIXME must be filtered by id and projectId
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

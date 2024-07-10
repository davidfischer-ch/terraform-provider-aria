// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &CustomNamingResource{}
var _ resource.ResourceWithImportState = &CustomNamingResource{}

func NewCustomNamingResource() resource.Resource {
	return &CustomNamingResource{}
}

// CustomNamingResource defines the resource implementation.
type CustomNamingResource struct {
	client *resty.Client
}

func (self *CustomNamingResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_custom_naming"
}

func (self *CustomNamingResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Custom Naming resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Resource identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "A friendly name",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Describe the resource in few sentences",
				Required:            true,
			},
			"projects": schema.ListNestedAttribute{
				MarkdownDescription: "Restrict the naming template to given projects (by filters).",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Resource identifier",
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"active": schema.BoolAttribute{
							MarkdownDescription: "TODO",
							Required:            true,
						},
						"org_default": schema.BoolAttribute{
							MarkdownDescription: "Default for the organization?",
							Required:            true,
						},
						"org_id": schema.StringAttribute{
							MarkdownDescription: "Organization identifier",
							Required:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"project_id": schema.StringAttribute{
							MarkdownDescription: "Projects identifier pattern (e.g. *).",
							Required:            true,
						},
						"project_name": schema.StringAttribute{
							MarkdownDescription: "Projects name pattern (e.g. *).",
							Required:            true,
						},
					},
				},
			},
			"templates": schema.MapNestedAttribute{
				MarkdownDescription: strings.Join([]string{
					"Resource naming patterns.",
					"Map key must be set to " +
						"\"resource_type.resource_type_name > static_pattern\"" +
						" for the provider to correlate API with state data.",
					" See example in documentation for details.",
					"Inspired by https://discuss.hashicorp.com/t/terraform-framework-optional-" +
						"inside-a-setnestedattribute-produces-a-does-not-correlate-with-any-element" +
						"-in-actual/62974/2.",
				}, "\n"),
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Resource identifier",
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Template name (valid for types that supports named templates)",
							Required:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"resource_type": schema.StringAttribute{
							MarkdownDescription: "Resource type, one of COMPUTE, COMPUTE_STORAGE, NETWORK, LOAD_BALANCER, RESOURCE_GROUP, GATEWAY, NAT, SECURITY_GROUP, GENERIC",
							Required:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Validators: []validator.String{
								stringvalidator.OneOf([]string{
									"COMPUTE",
									"COMPUTE_STORAGE",
									"NETWORK",
									"LOAD_BALANCER",
									"RESOURCE_GROUP",
									"GATEWAY",
									"NAT",
									"SECURITY_GROUP",
									"GENERIC",
								}...),
							},
						},
						"resource_type_name": schema.StringAttribute{
							MarkdownDescription: "Resource type name (e.g. Machine)",
							Required:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"resource_default": schema.BoolAttribute{
							MarkdownDescription: "True when static pattern is empty (automatically inferred by the provider)",
							Computed:            true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.RequiresReplace(),
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						"unique_name": schema.BoolAttribute{
							MarkdownDescription: "TODO",
							Required:            true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.RequiresReplace(),
							},
						},
						"pattern": schema.StringAttribute{
							MarkdownDescription: "TODO",
							Required:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"static_pattern": schema.StringAttribute{
							MarkdownDescription: "TODO",
							Required:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"start_counter": schema.Int32Attribute{
							MarkdownDescription: "TODO",
							Computed:            true,
							Optional:            true,
							Default:             int32default.StaticInt32(1),
							PlanModifiers: []planmodifier.Int32{
								int32planmodifier.RequiresReplace(),
							},
						},
						"increment_step": schema.Int32Attribute{
							MarkdownDescription: "TODO",
							Computed:            true,
							Optional:            true,
							Default:             int32default.StaticInt32(1),
							PlanModifiers: []planmodifier.Int32{
								int32planmodifier.RequiresReplace(),
							},
						},
					},
				},
			},
		},
	}
}

func (self *CustomNamingResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *CustomNamingResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var naming CustomNamingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &naming)...)
	if resp.Diagnostics.HasError() {
		return
	}

	namingRaw, diags := naming.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.R().
		SetQueryParam("apiVersion", IAAS_API_VERSION).
		SetBody(namingRaw).
		SetResult(&namingRaw).
		Post("iaas/api/naming")
	err = handleAPIResponse(ctx, response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", naming.String(), err))
		return
	}

	// Save custom naming into Terraform state
	resp.Diagnostics.Append(naming.FromAPI(ctx, namingRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &naming)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", naming.String()))
}

func (self *CustomNamingResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var naming CustomNamingModel
	resp.Diagnostics.Append(req.State.Get(ctx, &naming)...)
	if resp.Diagnostics.HasError() {
		return
	}

	namingId := naming.Id.ValueString()
	var namingRaw CustomNamingAPIModel
	response, err := self.client.R().
		SetQueryParam("apiVersion", IAAS_API_VERSION).
		SetResult(&namingRaw).
		Get("iaas/api/naming/" + namingId)

	// Handle gracefully a resource that has vanished on the platform
	// Beware that some APIs respond with HTTP 404 instead of 403 ...
	if response.StatusCode() == 404 {
		tflog.Debug(ctx, fmt.Sprintf("%s not found", naming.String()))
		resp.State.RemoveResource(ctx)
		return
	}

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read %s, got error: %s", naming.String(), err))
		return
	}

	// Save updated custom naming into Terraform state
	resp.Diagnostics.Append(naming.FromAPI(ctx, namingRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &naming)...)
}

func (self *CustomNamingResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var naming CustomNamingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &naming)...)
	if resp.Diagnostics.HasError() {
		return
	}

	namingRaw, diags := naming.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.R().
		SetQueryParam("apiVersion", IAAS_API_VERSION).
		SetBody(namingRaw).
		SetResult(&namingRaw).
		Put("iaas/api/naming") // Its not a mistake...

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", naming.String(), err))
		return
	}

	// Save updated custom naming into Terraform state
	resp.Diagnostics.Append(naming.FromAPI(ctx, namingRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &naming)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", naming.String()))
}

func (self *CustomNamingResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var naming CustomNamingModel
	resp.Diagnostics.Append(req.State.Get(ctx, &naming)...)
	if resp.Diagnostics.HasError() {
		return
	}

	namingId := naming.Id.ValueString()
	if len(namingId) == 0 {
		return
	}

	resp.Diagnostics.Append(
		DeleteIt(
			self.client,
			ctx,
			naming.String(),
			"iaas/api/naming/"+namingId,
			IAAS_API_VERSION,
		)...,
	)
}

func (self *CustomNamingResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// FIXME must be filtered by id and projectId
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

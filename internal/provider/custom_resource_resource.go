// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
			"id": schema.StringAttribute{
				MarkdownDescription: "Resource identifier",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "A friendly name",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Describe the resource in few sentences",
				Required:            true,
			},
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
			"properties": schema.ListNestedAttribute{
				MarkdownDescription: "Resource's properties",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"title": schema.StringAttribute{
							MarkdownDescription: "Title",
							Required:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Description",
							Required:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type, one of string, integer, number, boolean, object, array",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf([]string{"array", "boolean", "integer", "number", "object", "string"}...),
							},
						},
						"default": schema.StringAttribute{
							MarkdownDescription: strings.Join([]string{
								"Default value (JSON encoded default value).",
								"Should be a dynamic type, but Terraform SDK returns this issue:",
								"Dynamic types inside of collections are not currently supported in terraform-plugin-framework.",
								"If underlying dynamic values are required, replace the 'properties' attribute definition with DynamicAttribute instead.",
							}, "\n"),
							Computed: true,
							Optional: true,
						},
						"encrypted": schema.BoolAttribute{
							MarkdownDescription: "Encrypted?",
							Required:            true,
						},
						"read_only": schema.BoolAttribute{
							MarkdownDescription: "Make the field read-only (in the form)",
							Required:            true,
						},
						"recreate_on_update": schema.BoolAttribute{
							MarkdownDescription: "Mark this field as writable once (resource will be recreated on change)",
							Required:            true,
						},
						"minimum": schema.Int64Attribute{
							MarkdownDescription: "Minimum value (incluse, valid for an integer)",
							Computed:            true,
							Optional:            true,
						},
						"maximum": schema.Int64Attribute{
							MarkdownDescription: "Maximum value (incluse, valid for an integer)",
							Computed:            true,
							Optional:            true,
						},
						"min_length": schema.Int64Attribute{
							MarkdownDescription: "Minimum length (valid for a string)",
							Computed:            true,
							Optional:            true,
						},
						"max_length": schema.Int64Attribute{
							MarkdownDescription: "Maximum length (valid for a string)",
							Computed:            true,
							Optional:            true,
						},
						"pattern": schema.StringAttribute{
							MarkdownDescription: "Pattern (valid for a string)",
							Computed:            true,
							Optional:            true,
							Default:             stringdefault.StaticString(""),
						},
						"one_of": schema.ListNestedAttribute{
							Computed: true,
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"const": schema.StringAttribute{
										MarkdownDescription: "Technical value",
										Required:            true,
									},
									"title": schema.StringAttribute{
										MarkdownDescription: "Display value",
										Required:            true,
									},
									"encrypted": schema.BoolAttribute{
										MarkdownDescription: "Encrypted?",
										Required:            true,
									},
								},
							},
						},
					},
				},
			},
			"create": schema.SingleNestedAttribute{
				MarkdownDescription: "Resource's create action",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "Runnable identifier",
						Required:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "Runnable name",
						Computed:            true,
					},
					"type": schema.StringAttribute{
						MarkdownDescription: "Runnable type, either abx.action or vro.workflow",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf([]string{"abx.action", "vro.workflow"}...),
						},
					},
					"project_id": schema.StringAttribute{
						MarkdownDescription: "Runnable's project identifier",
						Required:            true,
					},
					"input_parameters": schema.ListAttribute{
						MarkdownDescription: "TODO",
						ElementType:         types.StringType,
						Computed:            true,
						Optional:            true,
						Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
					},
				},
			},
			"read": schema.SingleNestedAttribute{
				MarkdownDescription: "Resource's read action",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "Runnable identifier",
						Required:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "Runnable name",
						Computed:            true,
					},
					"type": schema.StringAttribute{
						MarkdownDescription: "Runnable type, either abx.action or vro.workflow",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf([]string{"abx.action", "vro.workflow"}...),
						},
					},
					"project_id": schema.StringAttribute{
						MarkdownDescription: "Runnable's project identifier",
						Required:            true,
					},
					"input_parameters": schema.ListAttribute{
						MarkdownDescription: "TODO",
						ElementType:         types.StringType,
						Computed:            true,
						Optional:            true,
						Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
					},
				},
			},
			"update": schema.SingleNestedAttribute{
				MarkdownDescription: "Resource's update action",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "Runnable identifier",
						Required:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "Runnable name",
						Computed:            true,
					},
					"type": schema.StringAttribute{
						MarkdownDescription: "Runnable type, either abx.action or vro.workflow",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf([]string{"abx.action", "vro.workflow"}...),
						},
					},
					"project_id": schema.StringAttribute{
						MarkdownDescription: "Runnable's project identifier",
						Required:            true,
					},
					"input_parameters": schema.ListAttribute{
						MarkdownDescription: "TODO",
						ElementType:         types.StringType,
						Computed:            true,
						Optional:            true,
						Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
					},
				},
			},
			"delete": schema.SingleNestedAttribute{
				MarkdownDescription: "Resource's delete action",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "Runnable identifier",
						Required:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "Runnable name",
						Computed:            true,
					},
					"type": schema.StringAttribute{
						MarkdownDescription: "Runnable type, either abx.action or vro.workflow",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf([]string{"abx.action", "vro.workflow"}...),
						},
					},
					"project_id": schema.StringAttribute{
						MarkdownDescription: "Runnable's project identifier",
						Required:            true,
					},
					"input_parameters": schema.ListAttribute{
						MarkdownDescription: "TODO",
						ElementType:         types.StringType,
						Computed:            true,
						Optional:            true,
						Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
					},
				},
			},
			/* "additional_actions": schema.ListNestedAttribute{
				MarkdownDescription: "Additional actions (aka Day 2)",
				Computed: true,
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Action identifier",
							Computed:            true,
							PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Action name",
							Required: true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "Action display name",
							Required: true,
						},
						"provider_name": schema.StringAttribute{
							MarkdownDescription: "Provider name, one of xaas (and that's all, maybe)",
							Computed: true,
							Optional: true,
							Default:             stringdefault.StaticString("xaas"),
							PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						},
						"project_id": schema.StringAttribute{
							MarkdownDescription: "Action's project identifier",
							Computed: true,
							PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						},
						"resource_type": schema.StringAttribute{
							MarkdownDescription: "Resource type (e.g. Custom.DB.PostgreSQL)",
							Computed:            true,
							PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						},
						"runnable_item": schema.SingleNestedAttribute{
							MarkdownDescription: "Additional action's runnable",
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									MarkdownDescription: "Runnable identifier",
									Required:            true,
								},
								"name": schema.StringAttribute{
									MarkdownDescription: "Runnable name",
									Computed: true,
								},
								"type": schema.StringAttribute{
									MarkdownDescription: "Runnable type, either abx.action or vro.workflow",
									Required:            true,
									Validators: []validator.String{
										stringvalidator.OneOf([]string{"abx.action", "vro.workflow"}...),
									},
								},
								"project_id": schema.StringAttribute{
									MarkdownDescription: "Runnable's project identifier",
									Required:            true,
								},
								"input_parameters": schema.ListAttribute{
									MarkdownDescription: "TODO",
									ElementType:         types.StringType,
									Computed:            true,
									Optional:            true,
									Default:             listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
								},
							},
							Required: true,
						},
						"status": schema.StringAttribute{
							MarkdownDescription: "Resource status",
							Computed:            true,
						},
						"org_id": schema.StringAttribute{
							MarkdownDescription: "Organisation identifier",
							Computed:            true,
							PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						},
						"form_definition": schema.SingleNestedAttribute{
							MarkdownDescription: "Additional action's custom form",
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									MarkdownDescription: "Form identifier",
									Computed: true,
								},
								"name": schema.StringAttribute{
									MarkdownDescription: "Form name",
									Required:            true,
								},
								"type": schema.StringAttribute{
									MarkdownDescription: "Form type, requestForm",
									Computed:            true,
									Optional:            true,
									Default:             stringdefault.StaticString("requestForm"),
									PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
									Validators: []validator.String{
										stringvalidator.OneOf([]string{"requestForm"}...),
									},
								},
								"form": schema.StringAttribute{
									MarkdownDescription: "Form content in JSON (TODO nested attribute to define this instead of messing with JSON)",
									Required:            true,
								},
								"form_format": schema.StringAttribute{
									MarkdownDescription: "Form format either JSON or YAML, will be forced to JSON by Aria ...",
									Computed:            true,
									Default:             stringdefault.StaticString("JSON"),
								},
								"styles": schema.StringAttribute{
									MarkdownDescription: "Form stylesheet",
									Computed: true,
									Optional: true,
									Default:             stringdefault.StaticString(""),
									PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
								}
								"source_type": schema.StringAttribute{
									MarkdownDescription: "Form source type",
									Computed: true,
									Default: stringdefault.StaticString("resource.action"),
									PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
								}
								"status": schema.StringAttribute{
									MarkdownDescription: "Resource status, one of DRAFT, ON, or RELEASED",
									Computed:            true,
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.OneOf([]string{"DRAFT", "ON", "RELEASED"}...),
									},
								},
							},
							Computed: true,
							Optional: true,
							PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						},
					},
				},
			}, */
			"project_id": schema.StringAttribute{
				MarkdownDescription: "Project ID",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"org_id": schema.StringAttribute{
				MarkdownDescription: "Constant organisation identifier",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
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
		SetBody(resourceRaw).
		SetResult(&resourceRaw).
		Post("form-service/api/custom/resource-types")
	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create Custom Resource, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Custom Resource %s created", resourceRaw.Id))

	// Save custom resource into Terraform state
	resp.Diagnostics.Append(resource.FromAPI(ctx, resourceRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &resource)...)
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
		SetResult(&resourceRaw).
		Get("form-service/api/custom/resource-types/" + resourceId)

	// Handle gracefully a resource that has vanished on the platform
	// Beware that some APIs respond with HTTP 404 instead of 403 ...
	if response.StatusCode() == 404 {
		tflog.Debug(ctx, fmt.Sprintf("Custom Resource %s not found", resourceId))
		resp.State.RemoveResource(ctx)
		return
	}

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Custom Resource %s, got error: %s", resourceId, err))
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

	resourceId := resource.Id.ValueString()
	resourceRaw, diags := resource.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.R().
		SetBody(resourceRaw).
		SetResult(&resourceRaw).
		Post("form-service/api/custom/resource-types") // Its not a mistake...

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update Custom Resource %s, got error: %s", resourceId, err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Custom Resource %s updated", resourceId))

	// Save updated custom resource into Terraform state
	resp.Diagnostics.Append(resource.FromAPI(ctx, resourceRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &resource)...)
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
			"Custom Resource "+resourceId,
			"form-service/api/custom/resource-types/"+resourceId,
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

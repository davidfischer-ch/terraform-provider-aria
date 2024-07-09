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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
			"id": schema.StringAttribute{
				MarkdownDescription: "Resource identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
						"name": schema.StringAttribute{
							MarkdownDescription: "Name",
							Required:            true,
						},
						"title": schema.StringAttribute{
							MarkdownDescription: "Title",
							Required:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Description",
							Required:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type, one of string, integer, number or boolean. (handling object and array is not yet implemented)",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf([]string{"boolean", "integer", "number", "string"}...),
							},
						},
						"default": schema.StringAttribute{
							MarkdownDescription: strings.Join([]string{
								"Default value as string (will be seamlessly converted to appropriate type).",
								"This attribute should be a dynamic type, but Terraform SDK returns this issue:",
								"Dynamic types inside of collections are not currently supported in terraform-plugin-framework.",
								"If underlying dynamic values are required, replace the 'properties' attribute definition with DynamicAttribute instead.",
							}, "\n"),
							Computed: true,
							Optional: true,
						},
						"encrypted": schema.BoolAttribute{
							MarkdownDescription: "Encrypted?",
							Computed:            true,
							Optional:            true,
							Default:             booldefault.StaticBool(false),
						},
						"read_only": schema.BoolAttribute{
							MarkdownDescription: "Make the field read-only (in the form)",
							Computed:            true,
							Optional:            true,
							Default:             booldefault.StaticBool(false),
						},
						"recreate_on_update": schema.BoolAttribute{
							MarkdownDescription: "Mark this field as writable once (resource will be recreated on change)",
							Computed:            true,
							Optional:            true,
							Default:             booldefault.StaticBool(false),
						},
						"minimum": schema.Int64Attribute{
							MarkdownDescription: "Minimum value (inclusive, valid for an integer)",
							Computed:            true,
							Optional:            true,
						},
						"maximum": schema.Int64Attribute{
							MarkdownDescription: "Maximum value (inclusive, valid for an integer)",
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
							Required: true,
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
			/* "allocate:" TODO one of the optional main actions */
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
					"input_parameters": schema.ListNestedAttribute{
						Required: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									MarkdownDescription: "Type",
									Required:            true,
								},
								"name": schema.StringAttribute{
									MarkdownDescription: "Name",
									Required:            true,
								},
								"description": schema.StringAttribute{
									MarkdownDescription: "Description",
									Required:            true,
								},
							},
						},
					},
					"output_parameters": schema.ListNestedAttribute{
						Required: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									MarkdownDescription: "Type",
									Required:            true,
								},
								"name": schema.StringAttribute{
									MarkdownDescription: "Name",
									Required:            true,
								},
								"description": schema.StringAttribute{
									MarkdownDescription: "Description",
									Required:            true,
								},
							},
						},
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
					"input_parameters": schema.ListNestedAttribute{
						Required: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									MarkdownDescription: "Type",
									Required:            true,
								},
								"name": schema.StringAttribute{
									MarkdownDescription: "Name",
									Required:            true,
								},
								"description": schema.StringAttribute{
									MarkdownDescription: "Description",
									Required:            true,
								},
							},
						},
					},
					"output_parameters": schema.ListNestedAttribute{
						Required: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									MarkdownDescription: "Type",
									Required:            true,
								},
								"name": schema.StringAttribute{
									MarkdownDescription: "Name",
									Required:            true,
								},
								"description": schema.StringAttribute{
									MarkdownDescription: "Description",
									Required:            true,
								},
							},
						},
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
					"input_parameters": schema.ListNestedAttribute{
						Required: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									MarkdownDescription: "Type",
									Required:            true,
								},
								"name": schema.StringAttribute{
									MarkdownDescription: "Name",
									Required:            true,
								},
								"description": schema.StringAttribute{
									MarkdownDescription: "Description",
									Required:            true,
								},
							},
						},
					},
					"output_parameters": schema.ListNestedAttribute{
						Required: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									MarkdownDescription: "Type",
									Required:            true,
								},
								"name": schema.StringAttribute{
									MarkdownDescription: "Name",
									Required:            true,
								},
								"description": schema.StringAttribute{
									MarkdownDescription: "Description",
									Required:            true,
								},
							},
						},
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
					"input_parameters": schema.ListNestedAttribute{
						Required: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									MarkdownDescription: "Type",
									Required:            true,
								},
								"name": schema.StringAttribute{
									MarkdownDescription: "Name",
									Required:            true,
								},
								"description": schema.StringAttribute{
									MarkdownDescription: "Description",
									Required:            true,
								},
							},
						},
					},
					"output_parameters": schema.ListNestedAttribute{
						Required: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									MarkdownDescription: "Type",
									Required:            true,
								},
								"name": schema.StringAttribute{
									MarkdownDescription: "Name",
									Required:            true,
								},
								"description": schema.StringAttribute{
									MarkdownDescription: "Description",
									Required:            true,
								},
							},
						},
					},
				},
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "Project ID",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"org_id": schema.StringAttribute{
				MarkdownDescription: "Resource organisation identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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

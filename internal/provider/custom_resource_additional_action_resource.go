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
var _ resource.Resource = &CustomResourceAdditionalActionResource{}
var _ resource.ResourceWithImportState = &CustomResourceAdditionalActionResource{}

func NewCustomResourceAdditionalActionResource() resource.Resource {
	return &CustomResourceAdditionalActionResource{}
}

// CustomResourceAdditionalActionResource defines the resource implementation.
type CustomResourceAdditionalActionResource struct {
	client *resty.Client
}

func (self *CustomResourceAdditionalActionResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_custom_resource_additional_action"
}

func (self *CustomResourceAdditionalActionResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Custom Resource's additional action resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Additional action identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Action name",
				Required:            true,
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
				MarkdownDescription: "Custom resource type",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"runnable_item": schema.SingleNestedAttribute{
				MarkdownDescription: "Additional action's runnable",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "Runnable identifier",
						Required:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "Runnable name",
						Computed:            true,
					},
					"type": schema.StringAttribute{
						MarkdownDescription: "Runnable type, either abx.action or vro.workflow",
						Required:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Validators: []validator.String{
							stringvalidator.OneOf([]string{"abx.action", "vro.workflow"}...),
						},
					},
					"project_id": schema.StringAttribute{
						MarkdownDescription: "Runnable's project identifier",
						Required:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
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
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"org_id": schema.StringAttribute{
				MarkdownDescription: "Action organisation identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			/*"form_definition": schema.SingleNestedAttribute{
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
			},*/
		},
	}
}

func (self *CustomResourceAdditionalActionResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *CustomResourceAdditionalActionResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var action CustomResourceAdditionalActionModel
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

func (self *CustomResourceAdditionalActionResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var action CustomResourceAdditionalActionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &action)...)
	if resp.Diagnostics.HasError() {
		return
	}

	actionId := action.Id.ValueString()
	var actionRaw CustomResourceAdditionalActionAPIModel
	response, err := self.client.R().
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

func (self *CustomResourceAdditionalActionResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var action CustomResourceAdditionalActionModel
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

func (self *CustomResourceAdditionalActionResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var action CustomResourceAdditionalActionModel
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
		)...,
	)
}

func (self *CustomResourceAdditionalActionResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// FIXME must be filtered by id and projectId
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

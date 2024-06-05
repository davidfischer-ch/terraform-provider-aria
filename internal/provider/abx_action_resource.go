// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	// "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ABXActionResource{}
var _ resource.ResourceWithImportState = &IconResource{}

func NewABXActionResource() resource.Resource {
	return &ABXActionResource{}
}

// ABXActionResource defines the resource implementation.
type ABXActionResource struct {
	client *resty.Client
}

// ABXActionResourceModel describes the resource data model.
type ABXActionResourceModel struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	FAASProvider types.String `tfsdk:"faas_provider"`
	Type         types.String `tfsdk:"type"`

	RuntimeName    types.String `tfsdk:"runtime_name"`
	RuntimeVersion types.String `tfsdk:"runtime_version"`
	MemoryInMB     types.Int64  `tfsdk:"memory_in_mb"`
	TimeoutSeconds types.Int64  `tfsdk:"timeout_seconds"`
	Entrypoint     types.String `tfsdk:"entrypoint"`
	Dependencies   types.List   `tfsdk:"dependencies"`
	// Constants types.List[String] `tfsdk:"constants"`
	// Secrets types.List[String] `tfsdk:"secrets"`

	Source types.String `tfsdk:"source"`

	ProjectId types.String `tfsdk:"project_id"`
	OrgId     types.String `tfsdk:"org_id"`
}

func (self *ABXActionResourceModel) FromAPI(
	ctx context.Context,
	raw ABXActionResourceAPIModel,
) diag.Diagnostics {

	// https://go.dev/blog/maps
	/* inputs := map[string]string{}
	for key, value := range self.Constants.Elemets() {
		inputs["secret:"+key] = value
	}
	for key, value := range self.Secrets.Elements() {
		inputs["psecret:"+key] = value
	} */

	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Description = types.StringValue(raw.Description)
	self.FAASProvider = types.StringValue(strings.Replace(raw.FAASProvider, "", "auto", 1))
	self.RuntimeName = types.StringValue(raw.RuntimeName)
	self.RuntimeVersion = types.StringValue(raw.RuntimeVersion)
	self.MemoryInMB = types.Int64Value(raw.MemoryInMB)
	self.TimeoutSeconds = types.Int64Value(raw.TimeoutSeconds)
	self.Entrypoint = types.StringValue(raw.Entrypoint)
	self.Source = types.StringValue(raw.Source)
	self.ProjectId = types.StringValue(raw.ProjectId)
	self.OrgId = types.StringValue(raw.OrgId)

	dependencies, diags := types.ListValueFrom(
		ctx,
		types.StringType,
		strings.Split(raw.Dependencies, ","),
	)

	if !diags.HasError() {
		self.Dependencies = dependencies
	}

	return diags
}
func (self *ABXActionResourceModel) ToAPI() ABXActionResourceAPIModel {

	// https://go.dev/blog/maps
	/* inputs := map[string]string{}
	for key, value := range self.Constants.Elemets() {
		inputs["secret:"+key] = value
	}
	for key, value := range self.Secrets.Elements() {
		inputs["psecret:"+key] = value
	} */

	return ABXActionResourceAPIModel{
		Name:           self.Name.ValueString(),
		Description:    self.Description.ValueString(),
		FAASProvider:   strings.Replace(self.FAASProvider.ValueString(), "auto", "", 1),
		Type:           self.Type.ValueString(),
		RuntimeName:    self.RuntimeName.ValueString(),
		RuntimeVersion: self.RuntimeVersion.ValueString(),
		MemoryInMB:     self.MemoryInMB.ValueInt64(),
		TimeoutSeconds: self.TimeoutSeconds.ValueInt64(),
		Entrypoint:     self.Entrypoint.ValueString(),
		Dependencies:   "",                  // FIXME
		Inputs:         map[string]string{}, // FIXME
		Source:         self.Source.ValueString(),
		ProjectId:      self.ProjectId.ValueString(),
	}
}

// ABXActionResourceAPIModel describes the resource API model.
type ABXActionResourceAPIModel struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	FAASProvider string `json:"provider"`
	Type         string `json:"actionType"`

	RuntimeName    string            `json:"runtime"`
	RuntimeVersion string            `json:"runtimeVersion"`
	MemoryInMB     int64             `json:"memoryInMB"`
	TimeoutSeconds int64             `json:"timeoutSeconds"`
	Entrypoint     string            `json:"entrypoint"`
	Dependencies   string            `json:"dependencies"`
	Inputs         map[string]string `json:"inputs"`

	Source string `json:"source"`

	ProjectId string `json:"projectId"`
	OrgId     string `json:"orgId"`
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
			"id": schema.StringAttribute{
				MarkdownDescription: "Action identifier",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "A name (must be unique)",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Describe the action in few sentences",
				Required:            true,
			},
			"faas_provider": schema.StringAttribute{
				MarkdownDescription: "FaaS provider used for code execution",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("auto"),
				// "auto" "on-prem" "aws" "azure"
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of action",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("SCRIPT"),
				// SCRIPT, REST_CALL, REST_POLL, FLOW, VAULT, CYBERARK
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
			"memory_in_mb": schema.Int64Attribute{
				MarkdownDescription: "Runtime memory constraint in MB",
				Required:            true,
			},
			"timeout_seconds": schema.Int64Attribute{
				MarkdownDescription: "How long an action can run (default 600)",
				Computed:            true,
				Optional:            true,
				Default:             int64default.StaticInt64(600),
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

			"source": schema.StringAttribute{
				MarkdownDescription: "Action source code",
				Required:            true,
			},

			/* "constants": schema.MapAttribute{
				ElementType: types.StringType,
				Required: true,
			},

			"secrets": schema.MapAttribute{
				ElementType: types.StringType,
				Required: true,
			}, */

			// TODO Validate is set to prevent 400!
			"project_id": schema.StringAttribute{
				MarkdownDescription: "Required for non-system actions",
				Optional:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},

			"org_id": schema.StringAttribute{
				MarkdownDescription: "Organisation ID",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},

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
	var action ABXActionResourceModel
	var actionRaw ABXActionResourceAPIModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &action)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.R().
		SetBody(action.ToAPI()).
		SetResult(&actionRaw).
		Post("abx/api/resources/actions")

	err = handleAPIResponse(ctx, response, err, 200)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create ABX action, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("ABX action %s created", actionRaw.Id))

	// Save action into Terraform state
	resp.Diagnostics.Append(action.FromAPI(ctx, actionRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &action)...)
}

func (self *ABXActionResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var action ABXActionResourceModel
	var actionRaw ABXActionResourceAPIModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &action)...)
	if resp.Diagnostics.HasError() {
		return
	}

	actionId := action.Id.ValueString()
	response, err := self.client.R().
		SetResult(&actionRaw).
		Get("abx/api/resources/actions/" + actionId)

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
			fmt.Sprintf("Unable to read ABX action %s, got error: %s", actionId, err))
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
	var action ABXActionResourceModel
	var actionRaw ABXActionResourceAPIModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &action)...)
	if resp.Diagnostics.HasError() {
		return
	}

	actionId := action.Id.ValueString()
	response, err := self.client.R().
		SetBody(action.ToAPI()).
		SetResult(&actionRaw).
		Put("abx/api/resources/actions/" + actionId)

	err = handleAPIResponse(ctx, response, err, 200)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update ABX action %s, got error: %s", actionId, err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("ABX action %s updated", actionId))

	// Save action into Terraform state
	resp.Diagnostics.Append(action.FromAPI(ctx, actionRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &action)...)
}

func (self *ABXActionResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var action ABXActionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &action)...)
	if resp.Diagnostics.HasError() {
		return
	}

	actionId := action.Id.ValueString()
	if len(actionId) == 0 {
		return
	}

	response, err := self.client.R().Delete("abx/api/resources/actions/" + actionId)
	err = handleAPIResponse(ctx, response, err, 200)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to delete ABX action %s, got error: %s", actionId, err))
	}

	tflog.Debug(ctx, fmt.Sprintf("ABX action %s deleted", actionId))
}

func (self *ABXActionResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

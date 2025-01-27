// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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
	client *AriaClient
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
	resp.Schema = ResourceActionSchema()
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

	self.client.Mutex.Lock(ctx, action.LockKey())
	actionRaw, diags := self.ManageIt(ctx, &action, "create")
	self.client.Mutex.Unlock(ctx, action.LockKey())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save resource action into Terraform state
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

	var actionRaw ResourceActionAPIModel
	self.client.Mutex.RLock(ctx, action.LockKey())
	found, _, diags := self.client.ReadIt(ctx, &action, &actionRaw)
	self.client.Mutex.RUnlock(ctx, action.LockKey())
	resp.Diagnostics.Append(diags...)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	if !resp.Diagnostics.HasError() {
		// Save updated resource action into Terraform state
		resp.Diagnostics.Append(action.FromAPI(ctx, actionRaw)...)
		resp.Diagnostics.Append(resp.State.Set(ctx, &action)...)
	}
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

	self.client.Mutex.Lock(ctx, action.LockKey())
	actionRaw, diags := self.ManageIt(ctx, &action, "update")
	self.client.Mutex.Unlock(ctx, action.LockKey())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated resource action into Terraform state
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

	self.client.Mutex.Lock(ctx, action.LockKey())
	_, diags := self.ManageIt(ctx, &action, "delete")
	self.client.Mutex.Unlock(ctx, action.LockKey())
	resp.Diagnostics.Append(diags...)
}

func (self *ResourceActionResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// FIXME must be filtered by id and projectId
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// -------------------------------------------------------------------------------------------------

// Implement the magic behind the create, update and delete methods.
func (self *ResourceActionResource) ManageIt(
	ctx context.Context,
	action *ResourceActionModel,
	method string,
) (ResourceActionAPIModel, diag.Diagnostics) {

	var actionRaw ResourceActionAPIModel
	var someDiags diag.Diagnostics
	diags := diag.Diagnostics{}

	// Check method is valid
	if !slices.Contains([]string{"create", "update", "delete"}, method) {
		diags.AddError("Client error", fmt.Sprintf("BUG: Wrong method %s", method))
		return actionRaw, diags
	}

	/* Custom resource ... */
	if action.ForCustom() {
		var resourceRaw CustomResourceAPIModel
		resource := CustomResourceModel{Id: action.ResourceId}

		// Retrieve the custom resource
		tflog.Debug(ctx, fmt.Sprintf("Retrieve %s", resource.String()))
		found, _, someDiags := self.client.ReadIt(ctx, &resource, &resourceRaw)
		diags.Append(someDiags...)
		diags.Append(resource.FromAPI(ctx, resourceRaw)...)

		if !found || diags.HasError() {
			if !found {
				diags.AddError(
					"Client error",
					fmt.Sprintf(
						"Unable to %s %s: %s not found.",
						method, action.String(), resource.String()))
			}

			return actionRaw, diags
		}

		tflog.Debug(ctx, fmt.Sprintf("Analyze %s (status & additional actions)", resource.String()))

		/* Create & Update: Validate status */
		if method != "delete" && action.Status.ValueString() != resource.Status.ValueString() {
			diags.AddError(
				"Configuration error",
				fmt.Sprintf(
					"Unable to %s %s: Status %s must match %s status %s.",
					method,
					action.String(),
					action.Status.ValueString(),
					resource.String(),
					resource.Status.ValueString()))

			return actionRaw, diags
		}

		// Find matching action inside custom resource's additional actions and grab its index
		runnableId := action.RunnableItem.Id.ValueString()
		actionIndex := -1
		for index, additionalAction := range resource.AdditionalActions {
			if additionalAction.RunnableItem.Id.ValueString() == runnableId {
				actionIndex = index
				break
			}
		}

		/* Create: Insert the action to the custom resource's additional actions */
		if method == "create" {
			if actionIndex >= 0 {
				// Prevent registering the same action twice (match by runnable ID)
				diags.AddError(
					"Configuration error",
					fmt.Sprintf(
						"Unable to create %s: %s is already registered to %s.",
						action.String(), action.RunnableItem.String(), resource.String()))
			} else {
				resource.AdditionalActions = append(resource.AdditionalActions, *action)
			}

			/* Update: Overwrite matching action inside the custom resource's additional actions */
		} else if method == "update" {
			if actionIndex < 0 {
				diags.AddError(
					"Client error",
					fmt.Sprintf(
						"Unable to update %s: Unable to find action in %s additional actions.",
						action.String(), resource.String()))
			} else {
				resource.AdditionalActions[actionIndex] = *action
			}

			/* Delete: Remove the action from the custom resource's additional actions */
		} else {
			if actionIndex < 0 {
				// Nothing to do, action is already deleted
				return actionRaw, diags
			}
			// https://stackoverflow.com/questions/20545743
			resource.AdditionalActions = append(
				resource.AdditionalActions[:actionIndex],
				resource.AdditionalActions[actionIndex+1:]...)
		}

		if diags.HasError() {
			return actionRaw, diags
		}

		tflog.Debug(ctx, fmt.Sprintf("%s %s", strings.ToUpper(method), action.String()))

		// FIXME Deduplicate code by implementing UpdateIt
		// Copied from custom_resource_resource.go -> Update
		resourceRaw, someDiags = resource.ToAPI(ctx)
		diags.Append(someDiags...)
		if diags.HasError() {
			return actionRaw, diags
		}

		// Update the custom resource
		response, err := self.client.Client.R().
			SetQueryParam("apiVersion", FORM_API_VERSION).
			SetBody(resourceRaw).
			SetResult(&resourceRaw).
			Post(resource.UpdatePath())

		err = handleAPIResponse(ctx, response, err, []int{200})
		if err != nil {
			diags.AddError(
				"Client error",
				fmt.Sprintf("Unable to update %s, got error: %s", resource.String(), err))
			return actionRaw, diags
		}

		// Find the action inside the custom resource's additional actions (response from API)
		actionIndex = -1
		var index int
		for index, actionRaw = range resourceRaw.AdditionalActions {
			if actionRaw.RunnableItem.Id == runnableId {
				actionIndex = index
				break
			}
		}

		/* Delete: Ensure action was removed from additional actions */
		if method == "delete" {
			if actionIndex >= 0 {
				diags.AddError(
					"Client error",
					fmt.Sprintf(
						"Unable to %s %s: Found action in %s additional actions.",
						method, action.String(), resource.String()))
			}

			/* Create & Update: Ensure action is found in additional actions */
		} else {
			if actionIndex < 0 {
				diags.AddError(
					"Client error",
					fmt.Sprintf(
						"Unable to %s %s: Unable to find action in %s additional actions.",
						method, action.String(), resource.String()))
			}
		}

		/* Native resource ... */
	} else {
		actionRaw, someDiags = action.ToAPI(ctx)
		diags.Append(someDiags...)
		if diags.HasError() {
			return actionRaw, diags
		}

		/* Delete: Delete the resource action */
		if method == "delete" {
			diags.Append(self.client.DeleteIt(ctx, action)...)
			return actionRaw, diags
		}

		/* Create or update the resource action */
		var path string
		if method == "create" {
			path = action.CreatePath()
		} else {
			path = action.UpdatePath()
		}

		response, err := self.client.Client.R().
			SetQueryParam("apiVersion", FORM_API_VERSION).
			SetBody(actionRaw).
			SetResult(&actionRaw).
			Post(path)

		err = handleAPIResponse(ctx, response, err, []int{200})
		if err != nil {
			diags.AddError(
				"Client error",
				fmt.Sprintf("Unable to %s %s, got error: %s", method, action.String(), err))
		}
	}

	return actionRaw, diags
}

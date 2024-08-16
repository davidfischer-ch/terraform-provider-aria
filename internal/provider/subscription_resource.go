// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SubscriptionResource{}
var _ resource.ResourceWithImportState = &SubscriptionResource{}

func NewSubscriptionResource() resource.Resource {
	return &SubscriptionResource{}
}

// SubscriptionResource defines the resource implementation.
type SubscriptionResource struct {
	client *AriaClient
}

func (self *SubscriptionResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_subscription"
}

func (self *SubscriptionResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = SubscriptionSchema()
}

func (self *SubscriptionResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *SubscriptionResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var subscription SubscriptionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &subscription)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subscription.GenerateId()
	subscriptionId := subscription.Id.ValueString()
	subscriptionRaw, diags := subscription.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.Client.R().
		// TODO SetQueryParam("apiVersion", EVENT_BROKER_API_VERSION).
		SetBody(subscriptionRaw).
		Post("event-broker/api/subscriptions")
	err = handleAPIResponse(ctx, response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", subscription.String(), err))
		return
	}

	// Read (using API) to retrieve the subscription content (and not empty stuff)
	response, err = self.client.Client.R().
		// TODO SetQueryParam("apiVersion", EVENT_BROKER_API_VERSION).
		SetResult(&subscriptionRaw).
		Get("event-broker/api/subscriptions/" + subscriptionId)

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read %s, got error: %s", subscription.String(), err))
		return
	}

	// Save subscription into Terraform state
	resp.Diagnostics.Append(subscription.FromAPI(ctx, subscriptionRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &subscription)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", subscription.String()))
}

func (self *SubscriptionResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var subscription SubscriptionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &subscription)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subscriptionId := subscription.Id.ValueString()
	var subscriptionRaw SubscriptionAPIModel
	response, err := self.client.Client.R().
		// TODO SetQueryParam("apiVersion", EVENT_BROKER_API_VERSION).
		SetResult(&subscriptionRaw).
		Get("event-broker/api/subscriptions/" + subscriptionId)

	// Handle gracefully a resource that has vanished on the platform
	// Beware that some APIs respond with HTTP 404 instead of 403 ...
	if response.StatusCode() == 404 {
		tflog.Debug(ctx, fmt.Sprintf("%s not found", subscription.String()))
		resp.State.RemoveResource(ctx)
		return
	}

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read %s, got error: %s", subscription.String(), err))
		return
	}

	// Save updated subscription into Terraform state
	resp.Diagnostics.Append(subscription.FromAPI(ctx, subscriptionRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &subscription)...)
}

func (self *SubscriptionResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var subscription SubscriptionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &subscription)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subscriptionId := subscription.Id.ValueString()
	subscriptionRaw, diags := subscription.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.Client.R().
		// TODO SetQueryParam("apiVersion", EVENT_BROKER_API_VERSION).
		SetBody(subscriptionRaw).
		Post("event-broker/api/subscriptions")
	err = handleAPIResponse(ctx, response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", subscription.String(), err))
		return
	}

	// Read (using API) to retrieve the subscription content (and not empty stuff)
	response, err = self.client.Client.R().
		SetResult(&subscriptionRaw).
		Get("event-broker/api/subscriptions/" + subscriptionId)

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read %s, got error: %s", subscription.String(), err))
		return
	}

	// Save subscription into Terraform state
	resp.Diagnostics.Append(subscription.FromAPI(ctx, subscriptionRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &subscription)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", subscription.String()))
}

func (self *SubscriptionResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var subscription SubscriptionModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &subscription)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subscriptionId := subscription.Id.ValueString()
	if len(subscriptionId) == 0 {
		return
	}

	resp.Diagnostics.Append(
		self.client.DeleteIt(
			ctx,
			subscription.String(),
			"event-broker/api/subscriptions/"+subscriptionId,
			EVENT_BROKER_API_VERSION,
		)...,
	)
}

func (self *SubscriptionResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

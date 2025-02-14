// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &OrchestratorEnvironmentResource{}
var _ resource.ResourceWithImportState = &OrchestratorEnvironmentResource{}

func NewOrchestratorEnvironmentResource() resource.Resource {
	return &OrchestratorEnvironmentResource{}
}

// OrchestratorEnvironmentResource defines the resource implementation.
type OrchestratorEnvironmentResource struct {
	client *AriaClient
}

func (self *OrchestratorEnvironmentResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_orchestrator_environment"
}

func (self *OrchestratorEnvironmentResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = OrchestratorEnvironmentSchema()
}

func (self *OrchestratorEnvironmentResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *OrchestratorEnvironmentResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var environment OrchestratorEnvironmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &environment)...)
	if resp.Diagnostics.HasError() {
		return
	}

	environmentToAPI, diags := environment.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var environmentFromAPI OrchestratorEnvironmentAPIModel
	path := environment.CreatePath()
	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", GetVersionFromPath(path)).
		SetBody(environmentToAPI).
		SetResult(&environmentFromAPI).
		Post(path)
	err = handleAPIResponse(ctx, response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", environment.String(), err))
		return
	}

	// Save environment into Terraform state
	resp.Diagnostics.Append(environment.FromAPI(ctx, environmentFromAPI, response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &environment)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", environment.String()))

	// Optionally wait up-to-date then save updated environment into Terraform state
	resp.Diagnostics.Append(self.WaitUpToDate(ctx, &environment)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &environment)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", environment.String()))
}

func (self *OrchestratorEnvironmentResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var environment OrchestratorEnvironmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &environment)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var environmentFromAPI OrchestratorEnvironmentAPIModel
	found, response, someDiags := self.client.ReadIt(ctx, &environment, &environmentFromAPI)
	resp.Diagnostics.Append(someDiags...)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	if !resp.Diagnostics.HasError() {
		// Save updated environment into Terraform state
		resp.Diagnostics.Append(environment.FromAPI(ctx, environmentFromAPI, response)...)
		resp.Diagnostics.Append(resp.State.Set(ctx, &environment)...)
	}
}

func (self *OrchestratorEnvironmentResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var environment OrchestratorEnvironmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &environment)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform state data into the model
	var environmentFromState OrchestratorEnvironmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &environmentFromState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	environmentToAPI, diags := environment.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var environmentFromAPI OrchestratorEnvironmentAPIModel
	path := environment.UpdatePath()
	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", GetVersionFromPath(path)).
		SetHeader("x-vro-changeset-sha", environmentFromState.VersionId.ValueString()).
		SetBody(environmentToAPI).
		SetResult(&environmentFromAPI).
		Put(path)

	err = handleAPIResponse(ctx, response, err, []int{202})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", environment.String(), err))
		return
	}

	if !resp.Diagnostics.HasError() {
		// Save updated environment into Terraform state
		resp.Diagnostics.Append(environment.FromAPI(ctx, environmentFromAPI, response)...)
		resp.Diagnostics.Append(resp.State.Set(ctx, &environment)...)
	}

	// Optionally wait up-to-date then save updated environment into Terraform state
	resp.Diagnostics.Append(self.WaitUpToDate(ctx, &environment)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &environment)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", environment.String()))
}

func (self *OrchestratorEnvironmentResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var environment OrchestratorEnvironmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &environment)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(self.client.DeleteIt(ctx, &environment)...)
	}
}

func (self *OrchestratorEnvironmentResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("wait_up_to_date"), true)...)
}

// -------------------------------------------------------------------------------------------------

func (self *OrchestratorEnvironmentResource) WaitUpToDate(
	ctx context.Context,
	environment *OrchestratorEnvironmentModel,
) diag.Diagnostics {

	diags := diag.Diagnostics{}
	if !environment.WaitUpToDate.ValueBool() {
		return diags
	}

	name := environment.String()
	tflog.Debug(ctx, fmt.Sprintf("Wait %s to be up-to-date...", name))

	// Poll for environment to be up-to-date to 10 minutes (60 x 10 seconds)
	maxAttempts := 60
	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Poll resource until up-to-date
		time.Sleep(time.Duration(10) * time.Second)
		tflog.Debug(
			ctx,
			fmt.Sprintf("Poll %d of %d - Check %s is up-to-date...", attempt+1, maxAttempts, name))

		var environmentFromAPI OrchestratorEnvironmentAPIModel
		found, response, someDiags := self.client.ReadIt(ctx, environment, &environmentFromAPI)
		diags.Append(someDiags...)
		if !found {
			diags.AddError(
				"Client error",
				fmt.Sprintf("%s has vanished while waiting to be imported.", name))
			return diags
		}

		// Update environment from API
		diags.Append(environment.FromAPI(ctx, environmentFromAPI, response)...)
		if diags.HasError() || environment.IsUpToDate() {
			return diags
		}
	}

	diags.AddError(
		"Client error",
		fmt.Sprintf("Timeout while waiting for %s to be up-to-date.", name))
	return diags
}

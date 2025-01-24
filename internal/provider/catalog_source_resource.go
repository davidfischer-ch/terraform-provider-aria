// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &CatalogSourceResource{}

func NewCatalogSourceResource() resource.Resource {
	return &CatalogSourceResource{}
}

// CatalogSourceResource defines the resource implementation.
type CatalogSourceResource struct {
	client *AriaClient
}

func (self *CatalogSourceResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_catalog_source"
}

func (self *CatalogSourceResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = CatalogSourceSchema()
}

func (self *CatalogSourceResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *CatalogSourceResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var source CatalogSourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &source)...)
	if resp.Diagnostics.HasError() {
		return
	}

	source.CompleteWorkflows(ctx, self.client)
	sourceRaw, diags := source.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", FORM_API_VERSION).
		SetBody(sourceRaw).
		SetResult(&sourceRaw).
		Post(source.CreatePath())
	err = handleAPIResponse(ctx, response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", source.String(), err))
		return
	}

	// FIXME Add wait attribute and if enabled then wait until last_import_completed_at > last_import_started_at

	// Save catalog source into Terraform state
	resp.Diagnostics.Append(source.FromAPI(ctx, sourceRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &source)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", source.String()))
}

func (self *CatalogSourceResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var source CatalogSourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &source)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var sourceRaw CatalogSourceAPIModel
	found, _, readDiags := self.client.ReadIt(ctx, &source, &sourceRaw)
	resp.Diagnostics.Append(readDiags...)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	if !resp.Diagnostics.HasError() {
		source.CompleteWorkflows(ctx, self.client)

		// Save updated catalog source into Terraform state
		resp.Diagnostics.Append(source.FromAPI(ctx, sourceRaw)...)
		resp.Diagnostics.Append(resp.State.Set(ctx, &source)...)
	}
}

func (self *CatalogSourceResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var source CatalogSourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &source)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceRaw, diags := source.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", CATALOG_API_VERSION).
		SetBody(sourceRaw).
		SetResult(&sourceRaw).
		Post(source.UpdatePath())
	err = handleAPIResponse(ctx, response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update %s, got error: %s", source.String(), err))
		return
	}

	// FIXME Add wait attribute and if enabled then wait until last_import_completed_at > last_import_started_at

	// Save catalog source into Terraform state
	resp.Diagnostics.Append(source.FromAPI(ctx, sourceRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &source)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", source.String()))
}

func (self *CatalogSourceResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var source CatalogSourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &source)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(self.client.DeleteIt(ctx, &source)...)
	}
}

// Retrieve Workflow's integration attribute by calling the resources API endpoint
func (self *CatalogSourceModel) CompleteWorkflows(
	ctx context.Context,
	client *AriaClient,
) {
	// FIXME iterate over workflows and retrieve its details
  	/* client.Get() "/catalog/api/types/com.vmw.vro.workflow/data/workflows"
	  query = {
	    size   = ["20"]
	    page   = ["0"]
	    sort   = ["name,asc"]
	    filter = ["substringof('${aria_orchestrator_workflow.test.id}',id)"]
	  }
	} */
}

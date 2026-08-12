// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// OrchestratorTaskResource wraps the generic CRUD resource to special-case task creation.
type OrchestratorTaskResource struct {
	*GenericResource[OrchestratorTaskModel, *OrchestratorTaskModel, OrchestratorTaskAPIModel]
}

func NewOrchestratorTaskResource() resource.Resource {
	return &OrchestratorTaskResource{
		GenericResource: &GenericResource[
			OrchestratorTaskModel,
			*OrchestratorTaskModel,
			OrchestratorTaskAPIModel,
		]{
			config: GenericResourceConfig{
				TypeName:     "_orchestrator_task",
				SchemaFunc:   OrchestratorTaskSchema,
				CreateCodes:  []int{202},
				UpdateMethod: "POST",
			},
		},
	}
}

// Create special-cases the "suspended" state. The API refuses to create a task directly as
// "suspended" and silently falls back to "pending", which Terraform rejects as an inconsistent
// result after apply. When the plan asks for "suspended", create the task as "pending" first,
// then suspend it with a follow-up update.
func (self *OrchestratorTaskResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var model OrchestratorTaskModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	suspendAfterCreate := model.State.ValueString() == "suspended"

	toAPI, diags := model.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// A task must exist in a non-suspended state before it can be suspended. "pending" is the only
	// state the API accepts on creation: it rejects "scheduled" with 400 and silently downgrades
	// "suspended" to "pending".
	if suspendAfterCreate {
		toAPI.State = "pending"
	}

	var raw OrchestratorTaskAPIModel
	_, createDiags := self.client.CreateIt(&model, &raw, toAPI, self.config.CreateCodes...)
	resp.Diagnostics.Append(createDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(model.FromAPI(ctx, raw)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Persist the created task before attempting to suspend it. If the suspend update fails,
	// Terraform keeps tracking the (tainted) resource instead of forgetting it, and a later apply
	// reconciles it to the desired state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !suspendAfterCreate {
		tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", model.String()))
		return
	}

	// The create populated the id, which UpdatePath() needs. Switch to the desired state and
	// rebuild the request body from the freshly created task before updating.
	model.State = types.StringValue("suspended")

	suspendAPI, suspendDiags := model.ToAPI(ctx)
	resp.Diagnostics.Append(suspendDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var updated OrchestratorTaskAPIModel
	_, updateDiags := self.client.UpdateIt(
		&model, &updated, suspendAPI, self.config.getUpdateMethod(), self.config.UpdateCodes...)
	resp.Diagnostics.Append(updateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(model.FromAPI(ctx, updated)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	tflog.Debug(ctx, fmt.Sprintf("Created and suspended %s successfully", model.String()))
}

// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CatalogSourceConfigModel describes the resource data model.
type CatalogSourceConfigModel struct {
	SourceProjectId types.String `tfsdk:"source_project_id"`
	Workflows       types.List   `tfsdk:"workflows"`
}

// CatalogSourceConfigAPIModel describes the resource API model.
type CatalogSourceConfigAPIModel struct {
	SourceProjectId string                          `json:"sourceProjectId,omitempty"`
	Workflows       []CatalogSourceWorkflowAPIModel `json:"workflows,omitempty"`
}

func (self *CatalogSourceConfigModel) FromAPI(
	ctx context.Context,
	raw CatalogSourceConfigAPIModel,
) diag.Diagnostics {
	self.SourceProjectId = types.StringValue(raw.SourceProjectId)

	diags := diag.Diagnostics{}

	// Convert workflows from raw to list
	if raw.Workflows == nil {
		self.Workflows = types.ListNull(self.Workflows.ElementType(ctx))
	} else {
		workflows := []CatalogSourceWorkflowModel{}
		for _, workflowRaw := range raw.Workflows {
			workflow := CatalogSourceWorkflowModel{}
			diags.Append(workflow.FromAPI(ctx, workflowRaw)...)
			workflows = append(workflows, workflow)
		}

		var someDiags diag.Diagnostics
		self.Workflows, someDiags = types.ListValueFrom(
			ctx, self.Workflows.ElementType(ctx), workflows,
		)
		diags.Append(someDiags...)
	}

	return diags
}

func (self CatalogSourceConfigModel) ToAPI(
	ctx context.Context,
	name string,
) (CatalogSourceConfigAPIModel, diag.Diagnostics) {

	diags := diag.Diagnostics{}

	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/list
	var workflowsRaw []CatalogSourceWorkflowAPIModel
	if self.Workflows.IsUnknown() {
		diags.AddError(
			"Configuration error",
			fmt.Sprintf("Unable to manage %s, workflows is unknown", name))
	} else if self.Workflows.IsNull() {
		workflowsRaw = nil
	} else {
		// Extract workflows from list value and then convert to raw
		workflows := make([]CatalogSourceWorkflowModel, 0, len(self.Workflows.Elements()))
		diags.Append(self.Workflows.ElementsAs(ctx, &workflows, false)...)
		if !diags.HasError() {
			for _, workflow := range workflows {
				workflowRaw, someDiags := workflow.ToAPI(ctx)
				workflowsRaw = append(workflowsRaw, workflowRaw)
				diags.Append(someDiags...)
			}
		}
	}

	return CatalogSourceConfigAPIModel{
		SourceProjectId: self.SourceProjectId.ValueString(),
		Workflows:       workflowsRaw,
	}, diags
}

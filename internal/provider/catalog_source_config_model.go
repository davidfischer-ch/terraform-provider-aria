// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

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
	workflows, diags := CatalogSourceWorkflowModelListFromAPI(ctx, raw.Workflows)
	self.SourceProjectId = types.StringValue(raw.SourceProjectId)
	self.Workflows = workflows
	return diags
}

func (self CatalogSourceConfigModel) ToAPI(
	ctx context.Context,
	name string,
) (CatalogSourceConfigAPIModel, diag.Diagnostics) {
	workflowsRaw, diags := CatalogSourceWorkflowModelListToAPI(ctx, self.Workflows, name)
	return CatalogSourceConfigAPIModel{
		SourceProjectId: self.SourceProjectId.ValueString(),
		Workflows:       workflowsRaw,
	}, diags
}

// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CatalogSourceWorkflowModel describes the resource data model.
type CatalogSourceWorkflowModel struct {
	Id      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Version types.String `tfsdk:"version"`
}

// CatalogSourceWorkflowAPIModel describes the resource API model.
type CatalogSourceWorkflowAPIModel struct {
	Id      string `tfsdk:"id,omitempty"`
	Name    string `tfsdk:"name"`
	Version string `tfsdk:"version"`
}

// Used to convert structure to a types.Object.
func CatalogSourceWorkflowAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":      types.StringType,
		"name":    types.StringType,
		"version": types.StringType,
	}
}

func (self *CatalogSourceWorkflowModel) FromAPI(
	ctx context.Context,
	raw CatalogSourceWorkflowAPIModel,
) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Version = types.StringValue(raw.Version)
	return diag.Diagnostics{}
}

func (self CatalogSourceWorkflowModel) ToAPI(
	ctx context.Context,
) (CatalogSourceWorkflowAPIModel, diag.Diagnostics) {
	return CatalogSourceWorkflowAPIModel{
		Id:      self.Id.ValueString(),
		Name:    self.Name.ValueString(),
		Version: self.Version.ValueString(),
	}, diag.Diagnostics{}
}

// Utils -------------------------------------------------------------------------------------------

func CatalogSourceWorkflowModelListFromAPI(
	ctx context.Context,
	workflowsRaw []CatalogSourceWorkflowAPIModel,
) (types.List, diag.Diagnostics) {
	// Convert input workflows from raw
	diags := diag.Diagnostics{}
	workflows := []CatalogSourceWorkflowModel{}
	for _, workflowRaw := range workflowsRaw {
		workflow := CatalogSourceWorkflowModel{}
		diags.Append(workflow.FromAPI(ctx, workflowRaw)...)
		workflows = append(workflows, workflow)
	}

	// Store inputs workflows to list value
	workflowAttrs := types.ObjectType{AttrTypes: CatalogSourceWorkflowAttributeTypes()}
	workflowsList, workflowsDiags := types.ListValueFrom(ctx, workflowAttrs, workflows)
	diags.Append(workflowsDiags...)

	return workflowsList, diags
}

func CatalogSourceWorkflowModelListToAPI(
	ctx context.Context,
	workflowsList types.List,
	name string,
) ([]CatalogSourceWorkflowAPIModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	workflowsRaw := []CatalogSourceWorkflowAPIModel{}

	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/list
	if workflowsList.IsNull() || workflowsList.IsUnknown() {
		diags.AddError(
			"Configuration error",
			fmt.Sprintf("Unable to manage %s is either null or unknown", name))
		return workflowsRaw, diags
	}

	// Extract input workflows from list value
	workflows := make([]CatalogSourceWorkflowModel, 0, len(workflowsList.Elements()))
	diags.Append(workflowsList.ElementsAs(ctx, &workflows, false)...)

	// Convert input workflows to raw
	for _, workflow := range workflows {
		workflowRaw, workflowDiags := workflow.ToAPI(ctx)
		workflowsRaw = append(workflowsRaw, workflowRaw)
		diags.Append(workflowDiags...)
	}

	return workflowsRaw, diags
}

// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CatalogSourceWorkflowModel describes the resource data model.
type CatalogSourceWorkflowModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Version     types.String `tfsdk:"version"`
}

// CatalogSourceWorkflowAPIModel describes the resource API model.
type CatalogSourceWorkflowAPIModel struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	/*Integration map[string]string `json:"integration"`
	  "integration": {
	      "name": "embedded-VRO",
	      "endpointUri": "https://vralab.ceti.etat-ge.ch:443",
	      "endpointConfigurationLink": "/resources/endpoints/8a430db3-924c-4d58-a29a-da811f9c992e"
	  }*/
}

func (self *CatalogSourceWorkflowModel) FromAPI(
	ctx context.Context,
	raw CatalogSourceWorkflowAPIModel,
) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Description = types.StringValue(raw.Description)
	self.Version = types.StringValue(raw.Version)
	return diag.Diagnostics{}
}

func (self CatalogSourceWorkflowModel) ToAPI(
	ctx context.Context,
) (CatalogSourceWorkflowAPIModel, diag.Diagnostics) {
	return CatalogSourceWorkflowAPIModel{
		Id:          self.Id.ValueString(),
		Name:        self.Name.ValueString(),
		Description: self.Description.ValueString(),
		Version:     self.Version.ValueString(),
	}, diag.Diagnostics{}
}

// Utils -------------------------------------------------------------------------------------------

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

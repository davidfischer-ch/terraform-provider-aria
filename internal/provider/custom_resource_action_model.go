// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CustomResourceActionModel describes the resource data model.
type CustomResourceActionModel struct {
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Type            types.String `tfsdk:"type"`
	ProjectId       types.String `tfsdk:"project_id"`
	InputParameters types.List   `tfsdk:"input_parameters"`
}

// CustomResourceActionAPIModel describes the resource API model.
type CustomResourceActionAPIModel struct {
	Id              string   `json:"id"`
	Name            string   `json:"name"`
	Type            string   `json:"type"`
	ProjectId       string   `json:"projectId"`
	InputParameters []string `json:"inputParameters"`
}

func (self *CustomResourceActionModel) FromAPI(
	ctx context.Context,
	raw CustomResourceActionAPIModel,
) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Type = types.StringValue(raw.Type)
	self.ProjectId = types.StringValue(raw.ProjectId)
	// FIXME self.InputParameters =
	return diag.Diagnostics{}
}

func (self *CustomResourceActionModel) ToAPI(
	ctx context.Context,
) (CustomResourceActionAPIModel, diag.Diagnostics) {
	return CustomResourceActionAPIModel{
		Id:        self.Id.ValueString(),
		Name:      self.Name.ValueString(),
		Type:      self.Type.ValueString(),
		ProjectId: self.ProjectId.ValueString(),
		// FIXME InputParameters:
	}, diag.Diagnostics{}
}

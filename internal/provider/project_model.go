// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ProjectModel describes the resource data model.
type ProjectModel struct {
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	OperationTimeout types.Int32  `tfsdk:"operation_timeout"`
	SharedResources  types.Bool   `tfsdk:"shared_resources"`

	Constraints ProjectConstraintsModel `tfsdk:"constraints"`
	Properties  types.Map               `tfsdk:"properties"`
	/*Cost ProjectCostModel `tfsdk:"cost"`*/

	OrgId types.String `tfsdk:"org_id"`
}

// ProjectAPIModel describes the resource API model.
type ProjectAPIModel struct {
	Id               string `json:"id,omitempty"`
	Name             string `json:"name"`
	OperationTimeout int32  `json:"operationTimeout"`
	SharedResources  bool   `json:"sharedResources"`

	Constraints ProjectConstraintsAPIModel `json:"constraints"`
	Properties  map[string]string          `json:"properties"`
	/*Cost ProjectCostAPIModel `json:"cost"`*/

	OrgId string `json:"orgId,omitempty"`
}

func (self ProjectModel) String() string {
	return fmt.Sprintf(
		"Project %s (%s)",
		self.Id.ValueString(),
		self.Name.ValueString())
}

func (self *ProjectModel) FromAPI(
	ctx context.Context,
	raw ProjectAPIModel,
) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.OperationTimeout = types.Int32Value(raw.OperationTimeout)
	self.SharedResources = types.BoolValue(raw.SharedResources)
	self.OrgId = types.StringValue(raw.OrgId)

	diags := self.Constraints.FromAPI(ctx, raw.Constraints)

	properties, propertiesDiags := types.MapValueFrom(ctx, types.StringType, raw.Properties)
	self.Properties = properties
	diags.Append(propertiesDiags...)

	return diags
}

func (self ProjectModel) ToAPI(
	ctx context.Context,
) (ProjectAPIModel, diag.Diagnostics) {

	constraintsRaw, diags := self.Constraints.ToAPI(ctx)

	propertiesRaw := make(map[string]string, len(self.Properties.Elements()))
	diags.Append(self.Properties.ElementsAs(ctx, &propertiesRaw, false)...)

	return ProjectAPIModel{
		Id:               self.Id.ValueString(),
		Name:             self.Name.ValueString(),
		OperationTimeout: self.OperationTimeout.ValueInt32(),
		SharedResources:  self.SharedResources.ValueBool(),
		Constraints:      constraintsRaw,
		Properties:       propertiesRaw,
		OrgId:            self.OrgId.ValueString(),
	}, diags
}

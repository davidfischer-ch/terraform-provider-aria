// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// OrchestratorConfigurationArrayModel describes the resource data model.
type OrchestratorConfigurationArrayModel struct {
	// Of type OrchestratorConfigurationArrayElementModel
	Elements types.List `tfsdk:"elements"`
}

// OrchestratorConfigurationArrayAPIModel describes the resource API model.
type OrchestratorConfigurationArrayAPIModel struct {
	Elements []OrchestratorConfigurationArrayElementAPIModel `json:"elements"`
}

func (self *OrchestratorConfigurationArrayModel) FromAPI(
	ctx context.Context,
	raw OrchestratorConfigurationArrayAPIModel,
) diag.Diagnostics {
	diags := diag.Diagnostics{}

	// Convert elements from raw to list
	attrs := types.ObjectType{
		AttrTypes: OrchestratorConfigurationArrayElementModel{}.AttributeTypes(),
	}
	if raw.Elements == nil {
		self.Elements = types.ListNull(attrs)
	} else {
		elements := []OrchestratorConfigurationArrayElementModel{}
		for _, elementRaw := range raw.Elements {
			element := OrchestratorConfigurationArrayElementModel{}
			diags.Append(element.FromAPI(ctx, elementRaw)...)
			elements = append(elements, element)
		}

		var someDiags diag.Diagnostics
		self.Elements, someDiags = types.ListValueFrom(ctx, attrs, elements)
		diags.Append(someDiags...)
	}

	return diags
}

func (self OrchestratorConfigurationArrayModel) ToAPI(
	ctx context.Context,
) (OrchestratorConfigurationArrayAPIModel, diag.Diagnostics) {

	diags := diag.Diagnostics{}

	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/list
	var elementsRaw []OrchestratorConfigurationArrayElementAPIModel
	if self.Elements.IsUnknown() {
		diags.AddError("Configuration error", "Unable to manage array, elements is unknown")
	} else if self.Elements.IsNull() {
		elementsRaw = nil
	} else {
		// Extract elements from list value and then convert to raw
		elements := make(
			[]OrchestratorConfigurationArrayElementModel, 0, len(self.Elements.Elements()),
		)
		diags.Append(self.Elements.ElementsAs(ctx, &elements, false)...)
		if !diags.HasError() {
			for _, element := range elements {
				elementRaw, someDiags := element.ToAPI(ctx)
				elementsRaw = append(elementsRaw, elementRaw)
				diags.Append(someDiags...)
			}
		}
	}

	return OrchestratorConfigurationArrayAPIModel{
		Elements: elementsRaw,
	}, diags
}

// Utils -------------------------------------------------------------------------------------------

// Used to convert structure to a types.Object.
func (self OrchestratorConfigurationArrayModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"elements": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: OrchestratorConfigurationArrayElementModel{}.AttributeTypes(),
			},
		},
	}
}

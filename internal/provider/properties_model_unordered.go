// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// UnorderedPropertiesModel describes the resource data model.
type UnorderedPropertiesModel map[string]PropertyModel

// UnorderedPropertiesAPIModel describes the resource API model.
type UnorderedPropertiesAPIModel map[string]PropertyAPIModel

func (self *UnorderedPropertiesModel) FromAPI(
	ctx context.Context,
	raw UnorderedPropertiesAPIModel,
) diag.Diagnostics {
	diags := diag.Diagnostics{}
	*self = UnorderedPropertiesModel{}
	selfRef := *self
	for propertyName, propertyRaw := range raw {
		property := PropertyModel{}
		diags.Append(property.FromAPI(ctx, propertyName, propertyRaw)...)
		selfRef[propertyName] = property
	}
	return diags
}

func (self UnorderedPropertiesModel) ToAPI(
	ctx context.Context,
) (UnorderedPropertiesAPIModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	properties := UnorderedPropertiesAPIModel{}
	for propertyKey, property := range self {
		propertyName, propertyRaw, propertyDiags := property.ToAPI(ctx)
		properties[propertyName] = propertyRaw
		diags.Append(propertyDiags...)
		if propertyKey != propertyName {
			diags.AddError(
				"Configuration error",
				fmt.Sprintf(
					"%s must be declared in map on key %s and not %s",
					property.String(), propertyName, propertyKey))
		}
	}
	return properties, diags
}

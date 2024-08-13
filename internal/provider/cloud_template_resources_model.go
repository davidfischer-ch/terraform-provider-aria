// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// CloudTemplateResourcesModel describes the resource data model.
type CloudTemplateResourcesModel map[string]CloudTemplateResourceModel

// CloudTemplateResourcesAPIModel describes the resource API model.
type CloudTemplateResourcesAPIModel map[string]CloudTemplateResourceAPIModel

func (self *CloudTemplateResourcesModel) FromAPI(
	ctx context.Context,
	raw CloudTemplateResourcesAPIModel,
) diag.Diagnostics {
	diags := diag.Diagnostics{}
	*self = CloudTemplateResourcesModel{}
	selfRef := *self
	for resourceName, resourceRaw := range raw {
		resource := CloudTemplateResourceModel{}
		diags.Append(resource.FromAPI(ctx, resourceName, resourceRaw)...)
		selfRef[resourceName] = resource
	}
	return diags
}

func (self CloudTemplateResourcesModel) ToAPI(
	ctx context.Context,
) (CloudTemplateResourcesAPIModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	resources := CloudTemplateResourcesAPIModel{}
	for resourceKey, resource := range self {
		resourceName, resourceRaw, resourceDiags := resource.ToAPI(ctx)
		resources[resourceName] = resourceRaw
		diags.Append(resourceDiags...)
		if resourceKey != resourceName {
			diags.AddError(
				"Configuration error",
				fmt.Sprintf(
					"%s must be declared in map on key %s and not %s",
					resource.String(), resourceName, resourceKey))
		}
	}
	return resources, diags
}

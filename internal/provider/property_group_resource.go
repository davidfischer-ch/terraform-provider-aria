// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import "github.com/hashicorp/terraform-plugin-framework/resource"

func NewPropertyGroupResource() resource.Resource {
	return &GenericResource[PropertyGroupModel, *PropertyGroupModel, PropertyGroupAPIModel]{
		config: GenericResourceConfig{
			TypeName:   "_property_group",
			SchemaFunc: PropertyGroupSchema,
		},
	}
}

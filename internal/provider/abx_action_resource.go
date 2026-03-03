// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import "github.com/hashicorp/terraform-plugin-framework/resource"

func NewABXActionResource() resource.Resource {
	return &GenericResource[ABXActionModel, *ABXActionModel, ABXActionAPIModel]{
		config: GenericResourceConfig{
			TypeName:    "_abx_action",
			SchemaFunc:  ABXActionSchema,
			CreateCodes: []int{200},
		},
	}
}

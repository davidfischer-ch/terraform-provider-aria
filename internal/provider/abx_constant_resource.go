// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import "github.com/hashicorp/terraform-plugin-framework/resource"

func NewABXConstantResource() resource.Resource {
	return &SimpleGenericResource[ABXConstantModel, *ABXConstantModel, ABXConstantAPIModel]{
		config: GenericResourceConfig{
			TypeName:    "_abx_constant",
			SchemaFunc:  ABXConstantSchema,
			CreateCodes: []int{200},
		},
	}
}

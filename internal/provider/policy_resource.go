// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import "github.com/hashicorp/terraform-plugin-framework/resource"

func NewPolicyResource() resource.Resource {
	return &GenericResource[PolicyModel, *PolicyModel, PolicyAPIModel]{
		config: GenericResourceConfig{
			TypeName:     "_policy",
			SchemaFunc:   PolicySchema,
			UpdateMethod: "POST",
			UpdateCodes:  []int{201},
		},
	}
}

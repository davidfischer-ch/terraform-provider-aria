// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ComputedCustomNamingResourceTypeSchema() schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: "Resource type, one of `COMPUTE`, `COMPUTE_STORAGE`, " +
			"`NETWORK`, `LOAD_BALANCER`, `RESOURCE_GROUP`, `GATEWAY`, `NAT`, " +
			"`SECURITY_GROUP`, `GENERIC`",
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseNonNullStateForUnknown(),
		},
		Validators: []validator.String{
			stringvalidator.OneOf([]string{
				"COMPUTE",
				"COMPUTE_STORAGE",
				"NETWORK",
				"LOAD_BALANCER",
				"RESOURCE_GROUP",
				"GATEWAY",
				"NAT",
				"SECURITY_GROUP",
				"GENERIC",
			}...),
		},
	}
}

func RequiredCustomNamingResourceTypeSchema() schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: "Resource type, one of `COMPUTE`, `COMPUTE_STORAGE`, " +
			"`NETWORK`, `LOAD_BALANCER`, `RESOURCE_GROUP`, `GATEWAY`, `NAT`, " +
			"`SECURITY_GROUP`, `GENERIC`",
		Required: true,
		Validators: []validator.String{
			stringvalidator.OneOf([]string{
				"COMPUTE",
				"COMPUTE_STORAGE",
				"NETWORK",
				"LOAD_BALANCER",
				"RESOURCE_GROUP",
				"GATEWAY",
				"NAT",
				"SECURITY_GROUP",
				"GENERIC",
			}...),
		},
	}
}

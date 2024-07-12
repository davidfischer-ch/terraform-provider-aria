// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func ComputedIdentifierSchema(description string) schema.StringAttribute {
	if len(description) == 0 {
		description = "Identifier"
	}
	return schema.StringAttribute{
		MarkdownDescription: description,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

func ComputedMutableIdentifierSchema(description string) schema.StringAttribute {
	if len(description) == 0 {
		description = "Identifier"
	}
	return schema.StringAttribute{
		MarkdownDescription: description,
		Computed:            true,
	}
}

func ComputedOrganizationIdSchema(description string) schema.StringAttribute {
	if len(description) == 0 {
		description = "Organization identifier"
	}
	return schema.StringAttribute{
		MarkdownDescription: description,
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

func RequiredIdentifierSchema(description string) schema.StringAttribute {
	if len(description) == 0 {
		description = "Identifier"
	}
	return schema.StringAttribute{
		MarkdownDescription: description,
		Required:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

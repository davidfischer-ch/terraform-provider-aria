// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
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

func ComputedMutableIdentifierSchema() schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: "Identifier",
		Computed:            true,
	}
}

func ComputedOrganizationIdSchema() schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: "Organization identifier",
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

func OptionalImmutableProjectIdSchema() schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: "Project identifier",
		Computed:            true,
		Optional:            true,
		Default:             stringdefault.StaticString(""),
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
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

func RequiredProjectId() schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: "Project identifier",
		Required:            true,
	}
}

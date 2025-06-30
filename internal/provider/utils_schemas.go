// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// Identifier

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

func OptionalIdentifierSchema(description string) schema.StringAttribute {
	if len(description) == 0 {
		description = "Identifier"
	}
	return schema.StringAttribute{
		MarkdownDescription: description,
		Computed:            true,
		Optional:            true,
	}
}

func RequiredIdentifierSchema(description string) schema.StringAttribute {
	if len(description) == 0 {
		description = "Identifier"
	}
	return schema.StringAttribute{
		MarkdownDescription: description,
		Required:            true,
	}
}

func RequiredImmutableIdentifierSchema(description string) schema.StringAttribute {
	if len(description) == 0 {
		description = "Identifier"
	}
	return schema.StringAttribute{
		MarkdownDescription: description,
		Required:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
}

// Description

func ComputedDescriptionSchema() schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: "Describe the resource in few sentences",
		Computed:            true,
	}
}

func RequiredDescriptionSchema() schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: "Describe the resource in few sentences",
		Required:            true,
	}
}

// Organization ID

func ComputedOrganizationIdSchema() schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: "Organization identifier",
		Computed:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

// Project ID

func OptionalImmutableProjectIdSchema() schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: "Project identifier. " +
			"Empty or unset means available for all projects." +
			IMMUTABLE,
		Computed: true,
		Optional: true,
		Default:  stringdefault.StaticString(""),
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

func RequiredImmutableProjectIdSchema() schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: "Project identifier" + IMMUTABLE,
		Required:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
}

func RequiredProjectIdSchema() schema.StringAttribute {
	return schema.StringAttribute{
		MarkdownDescription: "Project identifier",
		Required:            true,
	}
}

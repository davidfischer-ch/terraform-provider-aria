// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func CustomResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Custom Resource resource",
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"display_name": schema.StringAttribute{
				MarkdownDescription: "A friendly name",
				Required:            true,
			},
			"description": RequiredDescriptionSchema(),
			"resource_type": schema.StringAttribute{
				MarkdownDescription: "Define the type (must be unique, e.g. Custom.DB.PostgreSQL)",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"schema_type": schema.StringAttribute{
				MarkdownDescription: "Type of resource, one of ABX_USER_DEFINED (and that's all, maybe)",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("ABX_USER_DEFINED"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"ABX_USER_DEFINED"}...),
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Resource status, one of DRAFT, ON, or RELEASED",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("RELEASED"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"DRAFT", "ON", "RELEASED"}...),
				},
			},
			"properties": UnorderedPropertiesSchema("Resource's properties"),
			/* "allocate:" TODO one of the optional main actions */
			"create":     ResourceActionRunnableSchema("Create action"),
			"read":       ResourceActionRunnableSchema("Read action"),
			"update":     ResourceActionRunnableSchema("Update action"),
			"delete":     ResourceActionRunnableSchema("Delete action"),
			"project_id": OptionalImmutableProjectIdSchema(),
			"org_id":     ComputedOrganizationIdSchema(),
		},
	}
}

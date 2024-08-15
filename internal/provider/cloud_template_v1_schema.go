// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func CloudTemplateV1Schema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Cloud Template (v1 format) resource (WORK IN PROGRESS, DO NOT USE)",
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "Name",
				Required:            true,
			},
			"description": RequiredDescriptionSchema(),
			"project_id":  RequiredImmutableProjectIdSchema(),
			"request_scope_org": schema.BoolAttribute{
				MarkdownDescription: "Requestable from any project in organization?",
				Required:            true,
			},
			"inputs":    UnorderedPropertiesSchema("Cloud Template's properties"),
			"resources": CloudTemplateResourcesSchema(),
			"status": schema.StringAttribute{
				MarkdownDescription: "Status",
				Computed:            true,
			},
			"valid": schema.BoolAttribute{
				MarkdownDescription: "Cloud Template validation result",
				Computed:            true,
			},
			"validation_messages": schema.ListNestedAttribute{
				MarkdownDescription: "Cloud Template validation (error) messages.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"resource_name": schema.StringAttribute{
							MarkdownDescription: "Resource name",
							Computed:            true,
						},
						"path": schema.StringAttribute{
							MarkdownDescription: "Path",
							Computed:            true,
						},
						"message": schema.StringAttribute{
							MarkdownDescription: "Message",
							Computed:            true,
						},
					},
				},
			},
			"org_id": ComputedOrganizationIdSchema(),
		},
	}
}

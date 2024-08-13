// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func CloudTemplateV1Schema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Cloud Template (v1 format) resource",
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
			"org_id":    ComputedOrganizationIdSchema(),
		},
	}
}

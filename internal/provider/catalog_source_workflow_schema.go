// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// A Workflow declared inside a CatalogSourceConfigSchema.
func CatalogSourceWorkflowSchema() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"id": RequiredIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "Workflow name",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Workflow description",
				Required:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "Workflow version",
				Required:            true,
			},
			"integration": NestedIntegrationSchema(),
		},
	}
}

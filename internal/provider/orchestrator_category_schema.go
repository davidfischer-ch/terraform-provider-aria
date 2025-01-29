// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func OrchestratorCategorySchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Orchestrator category resource",
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "Category's name",
				Required:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "Category's path",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Category's type, " +
					"ConfigurationElementCategory, " +
					"PolicyTemplateCategory, " +
					"ResourceElementCategory or " +
					"WorkflowCategory",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"ConfigurationElementCategory",
						"PolicyTemplateCategory",
						"ResourceElementCategory",
						"WorkflowCategory",
					}...),
				},
			},
			"parent_id": schema.StringAttribute{
				MarkdownDescription: "Category's parent (empty string to make a root category).",
				Required:            true,
			},
		},
	}
}

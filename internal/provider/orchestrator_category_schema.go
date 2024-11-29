// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
				MarkdownDescription: "Category's type (e.g. WorkflowCategory)",
				Required:            true,
			},
			"parent_id": schema.StringAttribute{
				MarkdownDescription: "Category's parent (empty string to make a root category).",
				Required:            true,
			},
		},
	}
}

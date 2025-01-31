// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	dataschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
)

func OrchestratorConfigurationSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Orchestrator configuration resource",
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "Configuration name",
				Required:            true,
			},
			"description": RequiredDescriptionSchema(),
			"category_id": schema.StringAttribute{
				MarkdownDescription: "Where to store the configuration (Category's identifier)",
				Required:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "Configuration version (e.g. 1.0.0)",
				Required:            true,
			},
			"version_id": schema.StringAttribute{
				MarkdownDescription: "Configuration's latest changeset identifier",
				Computed:            true,
			},
			"attributes": schema.ListNestedAttribute{
				MarkdownDescription: "Attributes to store",
				Required:            true,
				NestedObject:        OrchestratorConfigurationAttributeSchema(),
			},
			"force_delete": schema.BoolAttribute{
				MarkdownDescription: "Force destroying the configuration (bypass references check).",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
		},
	}
}

func OrchestratorConfigurationDataSourceSchema() dataschema.Schema {
	return dataschema.Schema{
		MarkdownDescription: "Orchestrator configuration resource",
		Attributes: map[string]dataschema.Attribute{
			"id": RequiredIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "Configuration name",
				Computed:            true,
			},
			"description": ComputedDescriptionSchema(),
			"category_id": schema.StringAttribute{
				MarkdownDescription: "Where to store the configuration (Category's identifier)",
				Computed:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "Configuration version (e.g. 1.0.0)",
				Computed:            true,
			},
			"version_id": schema.StringAttribute{
				MarkdownDescription: "Configuration's latest changeset identifier",
				Computed:            true,
			},
			"attributes": schema.ListNestedAttribute{
				MarkdownDescription: "Attributes to store",
				Computed:            true,
				NestedObject:        ComputedOrchestratorConfigurationAttributeSchema(),
			},
		},
	}
}

// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
)

func OrchestratorEnvironmentRepositorySchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Orchestrator Repository (for Environments) resource",
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "Repository name",
				Required:            true,
			},
			"runtime": schema.StringAttribute{
				MarkdownDescription: "Runtime",
				Required:            true,
			},
			"location": schema.StringAttribute{
				MarkdownDescription: "Location",
				Required:            true,
			},
			"basic_auth": schema.BoolAttribute{
				MarkdownDescription: "Is basic authentication enabled?",
				Computed:            true,
			},
			"system_user": schema.StringAttribute{
				MarkdownDescription: "Username for basic authentication",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(""),
			},
			"system_credentials": schema.StringAttribute{
				MarkdownDescription: "Credentials for basic authentication",
				Computed:            true,
				Optional:            true,
				Sensitive:           true,
				Default:             stringdefault.StaticString(""),
			},
		},
	}
}

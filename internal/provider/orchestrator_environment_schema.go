// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func OrchestratorEnvironmentSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Orchestrator Environment resource",
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "Environment name",
				Required:            true,
			},
			"description": RequiredDescriptionSchema(),
			"version": schema.StringAttribute{
				MarkdownDescription: "Environment version",
				Required:            true,
			},
			"version_id": schema.StringAttribute{
				MarkdownDescription: "Configuration's latest changeset identifier",
				Computed:            true,
			},
			"runtime": schema.StringAttribute{
				MarkdownDescription: "Runtime",
				Required:            true,
			},
			"runtime_memory_limit": schema.Int64Attribute{
				MarkdownDescription: "Runtime memory constraint in bytes (can be 0 for unlimited)",
				Required:            true,
			},
			"runtime_timeout": schema.Int32Attribute{
				MarkdownDescription: "How long an action can run (in seconds) " +
					"(can be 0 for unlimited)",
				Required: true,
			},
			"dependencies": schema.MapAttribute{
				MarkdownDescription: "Dependencies to install on this environment",
				ElementType:         types.StringType,
				Required:            true,
			},
			"repositories": schema.MapAttribute{
				MarkdownDescription: "Repositories to use for downloading dependencies",
				ElementType:         types.StringType,
				Required:            true,
			},
			"variables": schema.MapAttribute{
				MarkdownDescription: "Variables to expose on this environment",
				ElementType:         types.StringType,
				Required:            true,
			},
			"bundle_hash": schema.StringAttribute{
				MarkdownDescription: "Bundle hash",
				Computed:            true,
			},
			"dependencies_install_execution_id": schema.StringAttribute{
				MarkdownDescription: "Dependencies Install Execution Identifier",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Status, either `UP_TO_DATE` or `PENDING_DOWNLOAD`, " +
					"maybe more (reverse-engineered the values)",
				Computed: true,
			},
			"validation_message": schema.StringAttribute{
				MarkdownDescription: "Validation message (if any, e.g. `DEPRECATED_RUNTIME`)",
				Computed:            true,
			},
			"wait_up_to_date": schema.BoolAttribute{
				MarkdownDescription: "Wait for the environment to be up-to-date (up to 10 minutes)",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(true),
			},
		},
	}
}

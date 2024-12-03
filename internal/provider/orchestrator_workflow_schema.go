// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
)

func OrchestratorWorkflowSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Orchestrator workflow resource",
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "Workflow name (e.g. Send Mail)",
				Required:            true,
			},
			"category_id": schema.StringAttribute{
				MarkdownDescription: "Where to store the workflow (Category's identifier)",
				Required:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "Workflow version (e.g. 1.0.0)",
				Required:            true,
			},
			"restart_mode": schema.Int32Attribute{
				MarkdownDescription: strings.Join([]string{
					"Workflow restart mode:",
					"Skip (0) - do not resume run from failure.",
					"Resume (1) - Resume workflow run failure.",
				}, "\n"),
				Required: true,
			},
			"resume_from_failed_mode": schema.Int32Attribute{
				MarkdownDescription: strings.Join([]string{
					"Resume workflow from failed behavior:",
					"Default (0) - System default - Follows the default behavior.",
					"Enabled (1) - If a workflow run fails, a pop-up window displays an option to " +
						"resume the workflow run.",
					"Disabled (2) - If a workflow run fails, it cannot be resumed.",
				}, "\n"),
				Required: true,
			},
			"force_delete": schema.BoolAttribute{
				MarkdownDescription: "Force destroying the workflow (bypass references check).",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
		},
	}
}

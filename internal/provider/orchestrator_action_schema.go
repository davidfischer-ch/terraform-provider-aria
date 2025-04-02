// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func OrchestratorActionSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Orchestrator action resource",
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "Action name (e.g. getVRAHost)",
				Required:            true,
			},
			"module": schema.StringAttribute{
				MarkdownDescription: "Where to store the action (e.g. ch.ocsin.core)",
				Required:            true,
			},
			"fqn": schema.StringAttribute{
				MarkdownDescription: "Action fully qualified name " +
					"(aka FQN, e.g. ch.ocsin.core/getVRAHost)",
				Required: true,
			},
			"description": RequiredDescriptionSchema(),
			"version": schema.StringAttribute{
				MarkdownDescription: "Action version",
				Required:            true,
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Environment identifier (omit or empty string if using a " +
					"standard runtime)",
				Computed: true,
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("runtime"),
					}...),
				},
			},
			"runtime": schema.StringAttribute{
				MarkdownDescription: "Runtime (omit or empty string for javascript " +
					"or when using a custom execution environment)",
				Computed: true,
				Optional: true,
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
			"script": schema.StringAttribute{
				MarkdownDescription: "Action source code",
				Required:            true,
			},
			"input_parameters": schema.ListNestedAttribute{
				MarkdownDescription: "Action input parameters",
				Required:            true,
				NestedObject:        ParameterSchema(),
			},
			"output_type": schema.StringAttribute{
				MarkdownDescription: "Action return type",
				Required:            true,
			},
			"force_delete": schema.BoolAttribute{
				MarkdownDescription: "Force destroying the action (bypass references check).",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"validation_message": schema.StringAttribute{
				MarkdownDescription: "Validation message (in case of error). " +
					"You may use a postcondition or a check to ensure its empty.",
				Computed: true,
			},
		},
	}
}

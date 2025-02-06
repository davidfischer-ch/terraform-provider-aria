// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func OrchestratorTaskSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Orchestrator Task resource",
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "Name",
				Required:            true,
			},
			"description": RequiredDescriptionSchema(),
			"href": schema.StringAttribute{
				MarkdownDescription: "Task URL (HATEOAS)",
				Computed:            true,
			},
			"recurrence_cycle": schema.StringAttribute{
				MarkdownDescription: "Recurrence cycle, one of one-time, every-minutes, " +
					"every-hours, every-days, every-weeks or every-months",
				Required: true,
			},
			"recurrence_pattern": schema.StringAttribute{
				MarkdownDescription: strings.Join([]string{
					"Recurrence pattern, examples:",
					"* one-time -> \"(Europe/Zurich) \"",
					"* every-minutes -> \"(Europe/Zurich) 30,\"",
					"* every-hours -> \"(Europe/Zurich) 10:00,\"",
					"* every-days -> \"(Europe/Zurich) 01:00:00,19:30:00,\"",
					"* every-weeks -> \"(Europe/Zurich) Monday 02:00:00,Friday 22:00:00,\"",
					"* every-months -> \"(Europe/Zurich) 01 00:00:00,12 00:00:00,\"",
				}, "\n"),
				Required: true,
			},
			"recurrence_start_date": schema.StringAttribute{
				MarkdownDescription: "Recurrence start timestamp (RFC3339)",
				CustomType:          timetypes.RFC3339Type{},
				Required:            true,
			},
			"recurrence_end_date": schema.StringAttribute{
				MarkdownDescription: "Recurrence end timestamp (RFC3339)",
				CustomType:          timetypes.RFC3339Type{},
				Computed:            true,
				Optional:            true,
			},
			"running_instance_id": schema.StringAttribute{
				MarkdownDescription: "Running instance ID",
				Computed:            true,
			},
			"start_mode": schema.StringAttribute{
				MarkdownDescription: "Start mode, either normal or start-in-the-past",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"normal", "start-in-the-past"}...),
				},
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "State",
				Computed:            true,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"pending", "suspended"}...),
				},
			},
			"user": schema.StringAttribute{
				MarkdownDescription: "User",
				Computed:            true,
				Optional:            true,
			},
			/*"input_parameters": schema.ListNestedAttribute{
				MarkdownDescription: "Workflow input parameters",
				Computed:            true,
				// Make it optional and use an actual model
				NestedObject: OrchestratorTaskInputParametersSchema(),
			},*/
			"workflow": OrchestratorTaskWorkflowSchema(),
		},
	}
}

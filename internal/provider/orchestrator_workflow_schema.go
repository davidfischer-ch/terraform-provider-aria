// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
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
			"description": RequiredDescriptionSchema(),
			"category_id": schema.StringAttribute{
				MarkdownDescription: "Where to store the workflow (Category's identifier)",
				Required:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "Workflow version (e.g. 1.0.0)",
				Required:            true,
			},
			"version_id": schema.StringAttribute{
				MarkdownDescription: "Workflow's latest changeset identifier",
				Computed:            true,
			},
			"allowed_operations": schema.StringAttribute{
				MarkdownDescription: "TODO (default is \"vef\")",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("vef"),
			},
			"attrib": schema.StringAttribute{
				MarkdownDescription: "Workflow attributes",
				CustomType:          jsontypes.NormalizedType{},
				Required:            true,
			},
			"object_name": schema.StringAttribute{
				MarkdownDescription: "TODO",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("workflow:name=generic"),
			},
			"position": PositionSchema(),
			"presentation": schema.StringAttribute{
				MarkdownDescription: "Workflow presentation",
				CustomType:          jsontypes.NormalizedType{},
				Required:            true,
			},
			"restart_mode": schema.Int32Attribute{
				MarkdownDescription: strings.Join([]string{
					"Workflow restart mode:",
					"Skip (0) - do not resume run from failure.",
					"Resume (1) - Resume workflow run failure.",
				}, "\n"),
				Required: true,
				/*Validators: []validator.String{
					stringvalidator.OneOf([]string{"skip", "resume"}...),
				},*/
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
				/*Validators: []validator.String{
					stringvalidator.OneOf([]string{"default", "enabled", "disabled"}...),
				},*/
			},
			"root_name": schema.StringAttribute{
				MarkdownDescription: "TODO (default is \"item0\")",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("item0"),
			},
			"workflow_item": schema.StringAttribute{
				MarkdownDescription: "Workflow item",
				CustomType:          jsontypes.NormalizedType{},
				Required:            true,
			},
			"input_parameters": schema.ListNestedAttribute{
				MarkdownDescription: "Workflow input parameters",
				Required:            true,
				NestedObject:        ParameterSchema(),
			},
			"output_parameters": schema.ListNestedAttribute{
				MarkdownDescription: "Workflow output parameters",
				Required:            true,
				NestedObject:        ParameterSchema(),
			},
			"input_forms": schema.StringAttribute{
				MarkdownDescription: "Workflow input forms",
				CustomType:          jsontypes.NormalizedType{},
				Required:            true,
			},
			"api_version": schema.StringAttribute{
				MarkdownDescription: "Orchestrator API Version (default is \"6.0.0\").",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("6.0.0"),
			},
			"editor_version": schema.StringAttribute{
				MarkdownDescription: "Orchestrator Editor Version (default is \"2.0\").",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("2.0"),
			},
			"integration": ComputedIntegrationSchema(),
			"force_delete": schema.BoolAttribute{
				MarkdownDescription: "Force destroying the workflow " +
					"(bypass references check, default is false).",
				Computed: true,
				Optional: true,
				Default:  booldefault.StaticBool(false),
			},
			"wait_imported": schema.BoolAttribute{
				MarkdownDescription: strings.Join([]string{
					"Wait for the workflow to be imported in the service " +
						"broker (up to 15 minutes, checked every 30 seconds, default is true).",
					"This ensure the integration attribute is set",
					"This is useful when non-orchestrator resources such as " +
						"`aria_resource_action` refer to this instance, ensuring the workflow is " +
						"available.",
					"If using an `aria_catalog_source` then you can rely on its own " +
						"`wait_imported` feature. However the `aria_catalog_source` must be declared " +
						"in the `depends_on` clause of any non-orchestrator resources making use of " +
						"this workflow.",
				}, "\n"),
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
		},
	}
}

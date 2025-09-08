// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ABXActionSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "ABX action resource",
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "A name (must be unique)",
				Required:            true,
			},
			"description": RequiredDescriptionSchema(),
			"faas_provider": schema.StringAttribute{
				MarkdownDescription: "FaaS provider used for code execution, one of `auto` " +
					"(default), `on-prem`, `aws` or `azure` " +
					"(automatically set by the platform if unset)",
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString("auto"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"auto", "on-prem", "aws", "azure"}...),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of action, one of `SCRIPT` (default), `REST_CALL`, " +
					"`REST_POLL`, `FLOW`, `VAULT` or `CYBERARK`",
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString("SCRIPT"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"SCRIPT",
						"REST_CALL",
						"REST_POLL",
						"FLOW",
						"VAULT",
						"CYBERARK",
					}...),
				},
			},
			"runtime_name": schema.StringAttribute{
				MarkdownDescription: "Runtime name (`python`, `nodejs`, ...)",
				Required:            true,
			},
			"runtime_version": schema.StringAttribute{
				MarkdownDescription: "Runtime version (3.10, ...)",
				Computed:            true,
				Optional:            true,
			},
			"cpu_shares": schema.Int32Attribute{
				MarkdownDescription: "Runtime CPU shares",
				Computed:            true,
				Optional:            true,
				Default:             int32default.StaticInt32(1024),
			},
			"memory_in_mb": schema.Int32Attribute{
				MarkdownDescription: "Runtime memory constraint in MB",
				Required:            true,
			},
			"timeout_seconds": schema.Int32Attribute{
				MarkdownDescription: "How long an action can run (default to 600)",
				Computed:            true,
				Optional:            true,
				Default:             int32default.StaticInt32(600),
			},
			"deployment_timeout_seconds": schema.Int32Attribute{
				MarkdownDescription: "How long ??",
				Computed:            true,
				Optional:            true,
				Default:             int32default.StaticInt32(900),
			},
			"entrypoint": schema.StringAttribute{
				MarkdownDescription: "Main function's name",
				Required:            true,
			},
			"dependencies": schema.ListAttribute{
				MarkdownDescription: "Dependencies (python packages, ...)",
				ElementType:         types.StringType,
				Required:            true,
			},
			"constants": schema.SetAttribute{
				MarkdownDescription: "ABX Constants to expose to the action",
				ElementType:         types.StringType,
				Required:            true,
			},
			"inputs": schema.MapAttribute{
				MarkdownDescription: "Inputs to expose to the action" + JSON_INSTEAD_OF_DYNAMIC_DISCLAIMER,
				ElementType:         jsontypes.NormalizedType{},
				Required:            true,
			},
			"secrets": schema.SetAttribute{
				MarkdownDescription: "Secrets to expose to the action",
				ElementType:         types.StringType,
				Required:            true,
			},
			"source": schema.StringAttribute{
				MarkdownDescription: "Action source code",
				Required:            true,
			},
			"project_id": OptionalImmutableProjectIdSchema(),
			"shared": schema.BoolAttribute{
				MarkdownDescription: "Flag indicating if the action can be shared across projects",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"system": schema.BoolAttribute{
				MarkdownDescription: "Flag indicating if the action is a system action",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"async_deployed": schema.BoolAttribute{
				MarkdownDescription: "TODO",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"org_id": ComputedOrganizationIdSchema(),
			/* created_millis int64 */
			/* updated_millis int64 */
			/* metadata {} */
			/* configuration {} */
			/* "self_link": schema.StringAttribute{
			    MarkdownDescription: "URL to the action",
			    Computed:            true,
			    PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			}, */
		},
	}
}

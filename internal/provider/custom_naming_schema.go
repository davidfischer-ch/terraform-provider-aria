// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func CustomNamingSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: strings.Join([]string{
			"Custom Naming resource",
			"",
			"**CAUTION: Updating projects**",
			"",
			"See [#29 Resource aria_custom_naming manage updating projects properly (scoped mode)]" +
				"(https://github.com/davidfischer-ch/terraform-provider-aria/issues/29)",
			"",
			"If the custon naming is scoped to projects, then updating projects should trigger an " +
				"update in place and not a replacement. However current version of the provider " +
				"will plan to replace the resource! The counters will be reset and its probably " +
				"something not desirable.",
			"",
			"As a workaround, you can update the list of projects in the config and manually on " +
				"the platform. Then plan to ensure no changes are detected by Terraform.",
			"",
			"Switching from organization <-> projects mode will never be possible without " +
				"replacing the resource.",
			"",
			"**CAUTION: Updating templates**",
			"",
			"Updating templates attributes will be shown as updatable in place." +
				" The nested template resource will be recreated, including its internal counter.",
		}, "\n"),
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "A friendly name",
				Required:            true,
			},
			"description": RequiredDescriptionSchema(),
			"projects": schema.ListNestedAttribute{
				MarkdownDescription: "Restrict the naming template to given projects (by filters).",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": ComputedIdentifierSchema(""),
						"active": schema.BoolAttribute{
							MarkdownDescription: "TODO",
							Required:            true,
						},
						"org_default": schema.BoolAttribute{
							MarkdownDescription: "Default for the organization?",
							Required:            true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.RequiresReplace(),
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						"org_id": schema.StringAttribute{
							MarkdownDescription: "Organization identifier",
							Required:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"project_id": schema.StringAttribute{
							MarkdownDescription: "Projects identifier pattern (e.g. *).",
							Required:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"project_name": schema.StringAttribute{
							MarkdownDescription: "Projects name pattern (e.g. *).",
							Required:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
								stringplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
			"templates": schema.MapNestedAttribute{
				MarkdownDescription: strings.Join([]string{
					"Resource naming patterns.",
					"Map key must be set to " +
						"\"resource_type.resource_type_name > static_pattern\"" +
						" for the provider to correlate API with state data.",
					" See example in documentation for details.",
					"Inspired by https://discuss.hashicorp.com/t/terraform-framework-optional-" +
						"inside-a-setnestedattribute-produces-a-does-not-correlate-with-any-" +
						"element-in-actual/62974/2.",
				}, "\n"),
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": ComputedMutableIdentifierSchema(""),
						"name": schema.StringAttribute{
							MarkdownDescription: "Template name (valid for types that supports " +
								"named templates)",
							Required: true,
						},
						"resource_type": schema.StringAttribute{
							MarkdownDescription: "Resource type, one of COMPUTE, COMPUTE_STORAGE, " +
								"NETWORK, LOAD_BALANCER, RESOURCE_GROUP, GATEWAY, NAT, " +
								"SECURITY_GROUP, GENERIC",
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf([]string{
									"COMPUTE",
									"COMPUTE_STORAGE",
									"NETWORK",
									"LOAD_BALANCER",
									"RESOURCE_GROUP",
									"GATEWAY",
									"NAT",
									"SECURITY_GROUP",
									"GENERIC",
								}...),
							},
						},
						"resource_type_name": schema.StringAttribute{
							MarkdownDescription: "Resource type name (e.g. Machine)",
							Required:            true,
						},
						"resource_default": schema.BoolAttribute{
							MarkdownDescription: "True when static pattern is empty (automatically" +
								" inferred by the provider)",
							Computed: true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						"unique_name": schema.BoolAttribute{
							MarkdownDescription: "TODO",
							Required:            true,
						},
						"pattern": schema.StringAttribute{
							MarkdownDescription: "TODO",
							Required:            true,
						},
						"static_pattern": schema.StringAttribute{
							MarkdownDescription: "TODO",
							Required:            true,
						},
						"start_counter": schema.Int32Attribute{
							MarkdownDescription: "TODO",
							Computed:            true,
							Optional:            true,
							Default:             int32default.StaticInt32(1),
						},
						"increment_step": schema.Int32Attribute{
							MarkdownDescription: "TODO",
							Computed:            true,
							Optional:            true,
							Default:             int32default.StaticInt32(1),
						},
					},
				},
			},
		},
	}
}

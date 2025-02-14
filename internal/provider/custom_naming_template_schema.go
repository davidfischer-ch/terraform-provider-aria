// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// A Template declared inside a CustomNamingSchema.
func CustomNamingTemplateSchema() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"id": ComputedMutableIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "Template name (valid for types that supports " +
					"named templates)",
				Required: true,
			},
			"resource_type": RequiredCustomNamingResourceTypeSchema(),
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
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"increment_step": schema.Int32Attribute{
				MarkdownDescription: "TODO",
				Computed:            true,
				Optional:            true,
				Default:             int32default.StaticInt32(1),
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"counters": schema.ListNestedAttribute{
				MarkdownDescription: "Counters",
				Computed:            true,
				NestedObject:        CustomNamingTemplateCounterSchema(),
			},
		},
	}
}

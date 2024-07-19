// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Computed as of today, only used by the ResourceActionSchema.
func CustomFormSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "Form definition",
		Computed:            true,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.UseStateForUnknown(),
		},
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "Form name",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Form type, requestForm",
				Computed:            true,
				Default:             stringdefault.StaticString("requestForm"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"requestForm"}...),
				},
			},
			"form": schema.StringAttribute{
				MarkdownDescription: "Form content in JSON " +
					"(TODO nested attribute to define this instead of messing with JSON)",
				Computed: true,
			},
			"form_format": schema.StringAttribute{
				MarkdownDescription: "Form format either JSON or YAML, " +
					"will be forced to JSON by Aria so you have no choice...",
				Computed: true,
				Default:  stringdefault.StaticString("JSON"),
			},
			"styles": schema.StringAttribute{
				MarkdownDescription: "Form stylesheet",
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"source_id": schema.StringAttribute{
				MarkdownDescription: "Form source ientifier",
				Computed:            true,
			},
			"source_type": schema.StringAttribute{
				MarkdownDescription: "Form source type",
				Computed:            true,
				Default:             stringdefault.StaticString("resource.action"),
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Resource status, one of DRAFT, ON, or RELEASED",
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"DRAFT", "ON", "RELEASED"}...),
				},
			},
			"tenant": schema.StringAttribute{
				MarkdownDescription: "TODO",
				Computed:            true,
			},
		},
	}
}

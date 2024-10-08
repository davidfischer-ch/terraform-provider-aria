// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func CustomFormSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Form definition",
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "Form name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Form type, requestForm",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("requestForm"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"requestForm"}...),
				},
			},
			"form": schema.StringAttribute{
				MarkdownDescription: "Form content in JSON",
				CustomType:          jsontypes.NormalizedType{},
				Required:            true,
			},
			"form_format": schema.StringAttribute{
				MarkdownDescription: "Form format either JSON or YAML, " +
					"will be forced to JSON by Aria so you have no choice...",
				Computed: true,
				Default:  stringdefault.StaticString("JSON"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"styles": schema.StringAttribute{
				MarkdownDescription: "Form stylesheet",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(""),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source_id": schema.StringAttribute{
				MarkdownDescription: "Form source ientifier",
				Required:            true,
			},
			"source_type": schema.StringAttribute{
				MarkdownDescription: "Form source type",
				Required:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Resource status, one of DRAFT, ON, or RELEASED",
				Computed:            true,
				Default:             stringdefault.StaticString("ON"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"DRAFT", "ON", "RELEASED"}...),
				},
			},
			"tenant": schema.StringAttribute{
				MarkdownDescription: "TODO",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// The optional CustomForm embeded inside a ResourceActionSchema.
func NestedCustomFormSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "Form definition",
		Computed:            true,
		Optional:            true,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.UseStateForUnknown(),
		},
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "Form name",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				MarkdownDescription: "Form content in JSON",
				CustomType:          jsontypes.NormalizedType{},
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"form_format": schema.StringAttribute{
				MarkdownDescription: "Form format either JSON or YAML, " +
					"will be forced to JSON by Aria so you have no choice...",
				Computed: true,
				Default:  stringdefault.StaticString("JSON"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"styles": schema.StringAttribute{
				MarkdownDescription: "Form stylesheet",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(""),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source_id": schema.StringAttribute{
				MarkdownDescription: "Form source ientifier",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source_type": schema.StringAttribute{
				MarkdownDescription: "Form source type",
				Computed:            true,
				Default:             stringdefault.StaticString("resource.action"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Resource status, one of DRAFT, ON, or RELEASED",
				Computed:            true,
				Default:             stringdefault.StaticString("ON"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"DRAFT", "ON", "RELEASED"}...),
				},
			},
			"tenant": schema.StringAttribute{
				MarkdownDescription: "TODO",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ResourceActionRunnableSchema(description string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: description,
		Required:            true,
		Attributes: map[string]schema.Attribute{
			"id": RequiredIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "Runnable name",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Runnable type, either abx.action or vro.workflow",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"abx.action", "vro.workflow"}...),
				},
			},
			"project_id": RequiredProjectIdSchema(),
			"endpoint_link": schema.StringAttribute{
				MarkdownDescription: "Integration API endpoint (e.g. /resources/endpoints/8a430db3-924c-4d58-a29a-da811f9c992e)",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(""),
			},
			"input_bindings": schema.StringAttribute{
				MarkdownDescription: "Input bindings" + JSON_INSTEAD_OF_DYNAMIC_DISCLAIMER,
				CustomType:          jsontypes.NormalizedType{},
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					// The value must stay whatever the state holds, null included, when the
					// configuration leaves it unset. The framework turns a computed attribute with
					// a null configuration value into "known after apply" otherwise, and the plan
					// of an unchanged resource is then never empty.
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"input_parameters": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: "Type",
							Required:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name",
							Required:            true,
						},
						"description": RequiredDescriptionSchema(),
					},
				},
			},
			"output_parameters": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: "Type",
							Required:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name",
							Required:            true,
						},
						"description": RequiredDescriptionSchema(),
					},
				},
			},
		},
	}
}

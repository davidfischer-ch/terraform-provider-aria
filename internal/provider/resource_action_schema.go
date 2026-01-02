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

func ResourceActionSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Resource's action resource (aka Day 2)",
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "Action name" + IMMUTABLE,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Action display name",
				Required:            true,
			},
			"description": RequiredDescriptionSchema(),
			"provider_name": schema.StringAttribute{
				MarkdownDescription: "Provider name, one of `xaas` or `vro-workflow` (and that's all, maybe)",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("xaas"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseNonNullStateForUnknown(),
				},
			},
			"resource_id": schema.StringAttribute{
				MarkdownDescription: "Resource identifier " +
					"(required if its a custom resource)" +
					IMMUTABLE,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"resource_type": schema.StringAttribute{
				MarkdownDescription: "Resource type" + IMMUTABLE,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"runnable_item": ResourceActionRunnableSchema("Action's runnable"),
			"criteria": schema.StringAttribute{
				MarkdownDescription: "Filtering criteria" + JSON_INSTEAD_OF_DYNAMIC_DISCLAIMER,
				CustomType:          jsontypes.NormalizedType{},
				Optional:            true,
			},
			"form_definition": NestedCustomFormSchema(),
			"status": schema.StringAttribute{
				MarkdownDescription: "Action status, either `DRAFT` or `RELEASED`",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("RELEASED"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"DRAFT", "RELEASED"}...),
				},
			},
			"project_id": OptionalImmutableProjectIdSchema(),
			"org_id":     ComputedOrganizationIdSchema(),
		},
	}
}

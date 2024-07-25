// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ResourceActionSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Native resource's action resource",
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "Action name",
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
				MarkdownDescription: "Provider name, one of xaas (and that's all, maybe)",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("xaas"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_type": schema.StringAttribute{
				MarkdownDescription: "Native resource type",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"runnable_item": ResourceActionRunnableSchema("Action's runnable"),
			"criteria": schema.StringAttribute{
				MarkdownDescription: "Filtering criteria (JSON encoded)",
				Optional:            true,
			},
			"form_definition": CustomFormSchema(),
			"status": schema.StringAttribute{
				MarkdownDescription: "Action status, either DRAFT or RELEASED",
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

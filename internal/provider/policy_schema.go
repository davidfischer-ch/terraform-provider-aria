// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func PolicySchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Policy resource",
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "Policy name",
				Required:            true,
			},
			"description": RequiredDescriptionSchema(),
			"enforcement_type": schema.StringAttribute{
				MarkdownDescription: "Enforcement type, either SOFT or HARD" + IMMUTABLE,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"SOFT", "HARD"}...),
				},
			},
			"type_id": schema.StringAttribute{
				MarkdownDescription: "Policy type" + IMMUTABLE + ", one of " +
					"com.vmware.policy.approval, " +
					"com.vmware.policy.catalog.entitlement, " +
					"com.vmware.policy.deployment.action, " +
					"com.vmware.policy.deployment.lease, " +
					"com.vmware.policy.deployment.limit, " +
					"com.vmware.policy.resource.quota",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"com.vmware.policy.approval",
						"com.vmware.policy.catalog.entitlement",
						"com.vmware.policy.deployment.action",
						"com.vmware.policy.deployment.lease",
						"com.vmware.policy.deployment.limit",
						"com.vmware.policy.resource.quota",
					}...),
				},
			},
			// Only when its not com.vmware.policy.approval?
			"criteria": schema.StringAttribute{
				MarkdownDescription: "Filtering criteria" + JSON_INSTEAD_OF_DYNAMIC_DISCLAIMER,
				CustomType:          jsontypes.NormalizedType{},
				Optional:            true,
			},
			// Only when its com.vmware.policy.approval?
			"scope_criteria": schema.StringAttribute{
				MarkdownDescription: "Scoping criteria" +
					IMMUTABLE +
					JSON_INSTEAD_OF_DYNAMIC_DISCLAIMER,
				CustomType: jsontypes.NormalizedType{},
				Optional:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"definition": schema.StringAttribute{
				MarkdownDescription: "Definition" + JSON_INSTEAD_OF_DYNAMIC_DISCLAIMER,
				CustomType:          jsontypes.NormalizedType{},
				Required:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Creation timestamp (RFC3339)",
				CustomType:          timetypes.RFC3339Type{},
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_by": schema.StringAttribute{
				MarkdownDescription: "User who created the resource",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated_at": schema.StringAttribute{
				MarkdownDescription: "Last update timestamp (RFC3339)",
				CustomType:          timetypes.RFC3339Type{},
				Computed:            true,
			},
			"last_updated_by": schema.StringAttribute{
				MarkdownDescription: "Last user who updated the resource",
				Computed:            true,
			},
			"project_id": OptionalImmutableProjectIdSchema(),
			"org_id":     ComputedOrganizationIdSchema(),
		},
	}
}

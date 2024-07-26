// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ProjectSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Project resource",
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "Project name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"operation_timeout": schema.Int32Attribute{
				MarkdownDescription: "Timeout (in seconds) that should be used for " +
					"Cloud Template operations and Provisioning tasks",
				Required: true,
			},
			"shared_resources": schema.BoolAttribute{
				MarkdownDescription: "Specifies whetever the resources are shared between " +
					"project's members or not",
				Required: true,
			},
			/*"administrators": ProjectMembershipSchema("Administrators"),*/
			/*"members": ProjectMembershipSchema("Members"),*/
			/*"viewers": ProjectMembershipSchema("Viewers"),*/
			/*"supervisors": ProjectMembershipSchema("Supervisors"),*/
			"constraints": ProjectConstraintsSchema(),
			"properties": schema.MapAttribute{
				MarkdownDescription: "Custom properties to attach to project's resources",
				ElementType:         types.StringType,
				Required:            true,
			},
			/*
				https://github.com/davidfischer-ch/terraform-provider-aria/issues/52
				"cost": ProjectCostSchema(),
			*/
			"org_id": ComputedOrganizationIdSchema(),
		},
	}
}

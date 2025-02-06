// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// The Workflow embedded inside an OrchestratorTaskSchema.
func OrchestratorTaskWorkflowSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "Configuration",
		Required:            true,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplace(),
		},
		Attributes: map[string]schema.Attribute{
			"id": RequiredImmutableIdentifierSchema("Workflow identifier"),
			"name": schema.StringAttribute{
				MarkdownDescription: "Workflow name",
				Required:            true,
			},
		},
	}
}

// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// The Position embedded inside an Orchestrator Workflow.
func NestedPositionSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "Position",
		Required:            true,
		Attributes: map[string]schema.Attribute{
			"x": schema.Float64Attribute{
				MarkdownDescription: "X",
				Required:            true,
			},
			"y": schema.Float64Attribute{
				MarkdownDescription: "Y",
				Required:            true,
			},
		},
	}
}

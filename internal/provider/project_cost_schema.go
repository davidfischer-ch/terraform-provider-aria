// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

/* https://github.com/davidfischer-ch/terraform-provider-aria/issues/52

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ProjectCostSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "Project cost",
		Attributes: map[string]schema.Attribute{
			"cost": schema.NumberAttribute{
				MarkdownDescription: "Cost of project",
				Computed:            true,
			},
			"cost_unit": schema.StringAttribute{
				MarkdownDescription: "Cost currency, 3 letters currency code (e.g. USD)",
				Computed:            true,
			},
			"cost_sync_time": schema.StringAttribute{
				CustomType: timetypes.RFC3339Type,
				MarkdownDescription: "The date as of which project cost was calculated.",
				Computed: true,
			},
			"message": schema.StringAttribute{
				MarkdownDescription: "Message regarding the project cost",
				Computed: true,
			},
			"code": schema.StringAttribute{
				MarkdownDescription: "Unique code for the message",
				Computed: true,
			},
		},
	}
} */

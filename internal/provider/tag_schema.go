// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func TagSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Tag resource",
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"key": schema.StringAttribute{
				MarkdownDescription: "Key" + IMMUTABLE,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "Value" + IMMUTABLE,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(""),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"force_delete": schema.BoolAttribute{
				MarkdownDescription: "Force destroying the tag (bypass references check).",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"keep_on_destroy": schema.BoolAttribute{
				MarkdownDescription: strings.Join([]string{
					"Keep the tag on destroy?",
					"This can help preventing issues if this tag " +
						"should never be destroyed for good reasons.",
					"Default value is false.",
				}, "\n"),
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
		},
	}
}

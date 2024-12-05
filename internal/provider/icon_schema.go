// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"strings"

	dataschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func IconSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Icon resource",
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema("Identifier (Aria seem to compute it from content)"),
			"path": schema.StringAttribute{
				MarkdownDescription: "Path (force recreation on change)",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"hash": schema.StringAttribute{
				MarkdownDescription: "Content SHA-256 (force recreation on change)",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"keep_on_destroy": schema.BoolAttribute{
				MarkdownDescription: strings.Join([]string{
					"Keep the icon on destroy?",
					"This can help preventing issues if sharing the same icon for multiple catalog items.",
					"Default value is false.",
				}, "\n"),
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
		},
	}
}

func IconDataSourceSchema() dataschema.Schema {
	return dataschema.Schema{
		MarkdownDescription: "Icon data source",
		Attributes: map[string]dataschema.Attribute{
			"id": dataschema.StringAttribute{
				MarkdownDescription: "Icon identifier",
				Required:            true,
			},
			"content": dataschema.StringAttribute{
				MarkdownDescription: "Icon content",
				Computed:            true,
			},
		},
	}
}

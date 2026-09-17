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
				MarkdownDescription: "Path" + IMMUTABLE,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"hash": schema.StringAttribute{
				MarkdownDescription: strings.Join([]string{
					"Checksum of the file at `path`, of your choosing." + IMMUTABLE,
					"Set it to `filesha256(...)` to replace the icon whenever its content " +
						"changes, the provider itself never reads the value.",
					"Compare with `content_hash` for what the platform stores back.",
				}, "\n"),
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content_hash": schema.StringAttribute{
				MarkdownDescription: strings.Join([]string{
					"SHA-256 of the content stored by the platform.",
					"This is not always the checksum of the file at `path`: the platform " +
						"rewrites some formats, SVG among them.",
				}, "\n"),
				Computed: true,
				PlanModifiers: []planmodifier.String{
					// The content is immutable, only a replacement can change this checksum.
					// Changing keep_on_destroy then plans without a spurious "known after apply".
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"keep_on_destroy": schema.BoolAttribute{
				MarkdownDescription: strings.Join([]string{
					"Keep the icon on destroy?",
					"This can help preventing issues if sharing the same icon " +
						"for multiple catalog items.",
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

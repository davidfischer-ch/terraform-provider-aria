// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func PropertyOneOfSchema() schema.NestedAttributeObject{
	return  schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"const": schema.StringAttribute{
				MarkdownDescription: "Technical value",
				Required:            true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "Display value",
				Required:            true,
			},
			"encrypted": schema.BoolAttribute{
				MarkdownDescription: "Encrypted?",
				Required:            true,
			},
		},
	}
}

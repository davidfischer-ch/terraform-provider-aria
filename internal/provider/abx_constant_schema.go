// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
)

func ABXConstantSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "ABX constant resource",
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "Name",
				Required:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "Value",
				Required:            true,
			},
			"encrypted": schema.BoolAttribute{
				MarkdownDescription: "Should be always unencrypted!",
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"org_id": ComputedOrganizationIdSchema(),
		},
	}
}

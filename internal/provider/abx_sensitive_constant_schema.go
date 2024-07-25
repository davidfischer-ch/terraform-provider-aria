// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
)

func ABXSensitiveConstantSchema() schema.Schema {
    return schema.Schema{
        MarkdownDescription: "ABX sensitive constant resource",
        Attributes: map[string]schema.Attribute{
            "id": ComputedIdentifierSchema(""),
            "name": schema.StringAttribute{
                MarkdownDescription: "Name",
                Required:            true,
            },
            "value": schema.StringAttribute{
                MarkdownDescription: "Value",
                Required:            true,
                Sensitive:           true,
            },
            "encrypted": schema.BoolAttribute{
                MarkdownDescription: "Should be always encrypted!",
                Computed:            true,
                Default:             booldefault.StaticBool(true),
            },
            "org_id": ComputedOrganizationIdSchema(),
        },
    }
}

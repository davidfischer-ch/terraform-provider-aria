// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func OrchestratorActionInputParameterSchema() schema.NestedAttributeObject {
    return schema.NestedAttributeObject{
        Attributes: map[string]schema.Attribute{
            "name": schema.StringAttribute{
                MarkdownDescription: "Parameter name",
                Required:            true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Parameter description",
                Required:            true,
            },
            "type": schema.StringAttribute{
                MarkdownDescription: "Parameter type",
                Required:            true,
            },
        },
    }
}

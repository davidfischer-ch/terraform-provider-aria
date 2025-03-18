// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// An Attribute declared inside a OrchestratorConfigurationSchema.
func OrchestratorConfigurationAttributeSchema() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name",
				Required:            true,
			},
			"description": RequiredDescriptionSchema(),
			"type": schema.StringAttribute{
				MarkdownDescription: "Type",
				Required:            true,
			},
			"value": OrchestratorConfigurationAttributeValueSchema(),
		},
	}
}

// An Attribute declared inside a OrchestratorConfigurationDataSourceSchema.
func ComputedOrchestratorConfigurationAttributeSchema() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name",
				Computed:            true,
			},
			"description": ComputedDescriptionSchema(),
			"type": schema.StringAttribute{
				MarkdownDescription: "Type",
				Computed:            true,
			},
			"value": ComputedOrchestratorConfigurationAttributeValueSchema(),
		},
	}
}

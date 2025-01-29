// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// An Attribute declared inside a OrchestratorConfigurationSchema.
func OrchestratorConfigurationAttributeSchema() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type, one of boolean or string.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"boolean",
						"string",
					}...),
				},
			},
			"value": OrchestratorConfigurationValueSchema(),
		},
	}
}

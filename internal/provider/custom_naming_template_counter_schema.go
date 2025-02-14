// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// A Counter declared inside a CustomNamingSchema.
func CustomNamingTemplateCounterSchema() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"id":            ComputedMutableIdentifierSchema(""),
			"resource_type": ComputedCustomNamingResourceTypeSchema(),
			"current_counter": schema.Int32Attribute{
				MarkdownDescription: "Current counter value",
				Computed:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "TODO",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

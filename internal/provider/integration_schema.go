// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	/*"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"*/)

// The computed integration embedded inside a CatalogSourceWorkflowSchema.
func NestedIntegrationSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "Integration",
		Computed:            true,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.UseStateForUnknown(),
		},
		Attributes: map[string]schema.Attribute{
			"name": schema.ListNestedAttribute{
				MarkdownDescription: "Integration name",
				Computed:            true,
				/*PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},*/
			},
			"endpoint_configuration_link": schema.StringAttribute{
				MarkdownDescription: "Integration endpoint configuration link",
				Computed:            true,
				/*PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},*/
			},
			"endpoint_uri": schema.StringAttribute{
				MarkdownDescription: "Integration endpoint URI",
				Computed:            true,
				/*PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},*/
			},
		},
	}
}

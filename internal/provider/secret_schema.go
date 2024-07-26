// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func SecretDataSourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Secret data source",
		Attributes: map[string]schema.Attribute{
			"id": RequiredIdentifierSchema(),
			"name": schema.StringAttribute{
				MarkdownDescription: "Secret name",
				Computed:            true,
			},
			"description": ComputedDescriptionSchema(),
			"org_id":      ComputedOrganizationIdSchema(),
			"org_scoped": schema.BoolAttribute{
				MarkdownDescription: "Scoped to the organization?",
				Computed:            true,
			},
			"project_ids": schema.SetAttribute{
				MarkdownDescription: "Restrict to given projects (an empty list means all)",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Creation date-time",
				Computed:            true,
			},
			"created_by": schema.StringAttribute{
				MarkdownDescription: "Ask VMware",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "Changed date-time",
				Computed:            true,
			},
			"updated_by": schema.StringAttribute{
				MarkdownDescription: "Ask VMware",
				Computed:            true,
			},
		},
	}
}

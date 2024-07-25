// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func CloudTemplateResourceSchema() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				MarkdownDescription: "Resource type",
				Required:            true,
			},
			/*"metadata": schema.MapNestedAttribute{
			      MarkdownDescription: "Resource metadata",
			      Required: true,
			  },
			  "properties": schema.MapNestedAttribute{
			      MarkdownDescription: "Resource properties",
			      Required: true,
			  },*/
			"allocate_per_instance": schema.BoolAttribute{
				MarkdownDescription: "TODO",
				Optional:            true,
			},
		},
	}
}

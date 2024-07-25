// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func CloudTemplateResourcesSchema() schema.MapNestedAttribute {
	return schema.MapNestedAttribute{
		MarkdownDescription: "Cloud Template's resources",
		NestedObject:        CloudTemplateResourceSchema(),
		Required:            true,
	}
}

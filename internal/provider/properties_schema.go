// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func OrderedPropertiesSchema(description string) schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		MarkdownDescription: description,
		Required:            true,
		NestedObject:        PropertySchema(),
	}
}

func UnorderedPropertiesSchema(description string) schema.MapNestedAttribute {
	return schema.MapNestedAttribute{
		MarkdownDescription: description,
		Required:            true,
		NestedObject:        PropertySchema(),
	}
}

// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func OrderedPropertyGroupSchema(description string) schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		MarkdownDescription: description,
		Required:            true,
		NestedObject:        PropertySchema(),
	}
}

func UnorderedPropertyGroupSchema(description string) schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: description,
		Required:            true,
		NestedObject:        PropertySchema(),
	}
}

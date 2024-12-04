// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func CatalogItemIconSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: strings.Join([]string{
			"Catalog Item's Icon resource",
			"",
			"Manage the icon of a catalog item.",
			"",
			"The create operation is implemented like the update.",
			"The destroy operation is a no-op and the catalog item's icon will be left unchanged.",
			"",
		}, "\n"),
		Attributes: map[string]schema.Attribute{
			"item_id": RequiredIdentifierSchema("Item identifier"),
			"icon_id": RequiredIdentifierSchema("Icon identifier"),
		},
	}
}

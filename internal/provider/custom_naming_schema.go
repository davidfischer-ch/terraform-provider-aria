// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func CustomNamingSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: strings.Join([]string{
			"Custom Naming resource",
			"",
			"**CAUTION: Updating projects**",
			"",
			"See [#29 Resource aria_custom_naming manage updating projects properly (scoped mode)]" +
				"(https://github.com/davidfischer-ch/terraform-provider-aria/issues/29)",
			"",
			"If the custon naming is scoped to projects, then updating projects should trigger an " +
				"update in place and not a replacement. However current version of the provider " +
				"will plan to replace the resource! The counters will be reset and its probably " +
				"something not desirable.",
			"",
			"As a workaround, you can update the list of projects in the config and manually on " +
				"the platform. Then plan to ensure no changes are detected by Terraform.",
			"",
			"Switching from organization <-> projects mode will never be possible without " +
				"replacing the resource.",
			"",
			"**CAUTION: Updating templates**",
			"",
			"Updating templates attributes will be shown as updatable in place." +
				" The nested template resource will be recreated, including its internal counter.",
		}, "\n"),
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "A friendly name",
				Required:            true,
			},
			"description": RequiredDescriptionSchema(),
			"projects": schema.ListNestedAttribute{
				MarkdownDescription: "Restrict the naming template to given projects (by filters).",
				Required:            true,
				NestedObject:        CustomNamingProjectFilterSchema(),
			},
			"templates": schema.MapNestedAttribute{
				MarkdownDescription: strings.Join([]string{
					"Resource naming patterns.",
					"Map key must be set to " +
						"\"resource_type.resource_type_name > static_pattern\"" +
						" for the provider to correlate API with state data.",
					" See example in documentation for details.",
					"Inspired by https://discuss.hashicorp.com/t/terraform-framework-optional-" +
						"inside-a-setnestedattribute-produces-a-does-not-correlate-with-any-" +
						"element-in-actual/62974/2.",
				}, "\n"),
				Required:     true,
				NestedObject: CustomNamingTemplateSchema(),
			},
		},
	}
}

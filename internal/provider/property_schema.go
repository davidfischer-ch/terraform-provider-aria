// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func PropertySchema() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name",
				Required:            true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "Title",
				Required:            true,
			},
			"description": RequiredDescriptionSchema(),
			"type": schema.StringAttribute{
				MarkdownDescription: "Type, one of array, boolean, integer, object, number or " +
					"string.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"array",
						"boolean",
						"integer",
						"object",
						"number",
						"string",
					}...),
				},
			},
			"default": schema.StringAttribute{
				MarkdownDescription: strings.Join([]string{
					"Default value as string (will be seamlessly converted to appropriate type).",
					"The string should be a JSON for arrays and objects.",
					"",
					"We should have implemented this attribute as a dynamic type (and not string).",
					"Unfortunately Terraform SDK returns this issue:",
					"Dynamic types inside of collections are not currently supported in " +
						"terraform-plugin-framework.",
					"If underlying dynamic values are required, replace the 'properties' " +
						"attribute definition with DynamicAttribute instead.",
				}, "\n"),
				Optional: true,
			},
			"encrypted": schema.BoolAttribute{
				MarkdownDescription: "Encrypted?",
				Required:            true,
			},
			"read_only": schema.BoolAttribute{
				MarkdownDescription: "Make the field read-only (in the form)",
				Required:            true,
			},
			"recreate_on_update": schema.BoolAttribute{
				MarkdownDescription: "Mark this field as writable once (resource will be " +
					"recreated on change)",
				Required: true,
			},
			"minimum": schema.Int64Attribute{
				MarkdownDescription: "Minimum value (inclusive, valid for an integer)",
				Optional:            true,
			},
			"maximum": schema.Int64Attribute{
				MarkdownDescription: "Maximum value (inclusive, valid for an integer)",
				Optional:            true,
			},
			"min_length": schema.Int32Attribute{
				MarkdownDescription: "Minimum length (valid for a string)",
				Optional:            true,
			},
			"max_length": schema.Int32Attribute{
				MarkdownDescription: "Maximum length (valid for a string)",
				Optional:            true,
			},
			"pattern": schema.StringAttribute{
				MarkdownDescription: "Pattern (valid for a string)",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(""),
			},
			/*"items": schema.StringAttribute{
				MarkdownDescription: "Items in JSON"
				CustomType:
			}*/
			"one_of": PropertyOneOfSchema(),
		},
	}
}

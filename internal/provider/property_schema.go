// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
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
				MarkdownDescription: "Default value" + JSON_INSTEAD_OF_DYNAMIC_DISCLAIMER,
				CustomType:          jsontypes.NormalizedType{},
				Optional:            true,
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
			"one_of": schema.ListNestedAttribute{
				MarkdownDescription: "Enumerate possible values",
				Optional:            true,
				NestedObject: PropertyOneOfSchema(),
			},
		},
	}
}

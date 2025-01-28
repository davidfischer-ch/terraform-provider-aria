// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dataschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// The integration embedded inside a CatalogSourceWorkflowSchema.
func NestedIntegrationSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "Integration",
		Required:            true,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Integration name",
				Required:            true,
			},
			"endpoint_configuration_link": schema.StringAttribute{
				MarkdownDescription: "Integration endpoint configuration link",
				Required:            true,
			},
			"endpoint_uri": schema.StringAttribute{
				MarkdownDescription: "Integration endpoint URI",
				Required:            true,
			},
		},
	}
}

func IntegrationDataSourceSchema() dataschema.Schema {
	return dataschema.Schema{
		MarkdownDescription: "Integration data source",
		Attributes: map[string]dataschema.Attribute{
			"type_id": dataschema.StringAttribute{
				MarkdownDescription: "Source type (com.vmw.vro.workflow, and that's all for now)",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"com.vmw.vro.workflow"}...),
				},
			},
			"name": dataschema.StringAttribute{
				MarkdownDescription: "Integration name",
				Computed:            true,
			},
			"endpoint_configuration_link": dataschema.StringAttribute{
				MarkdownDescription: "Integration endpoint configuration link",
				Computed:            true,
			},
			"endpoint_uri": dataschema.StringAttribute{
				MarkdownDescription: "Integration endpoint URI",
				Computed:            true,
			},
		},
	}
}

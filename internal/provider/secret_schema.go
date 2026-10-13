// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func SecretDataSourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Secret data source",
		Attributes: map[string]schema.Attribute{
			"id": RequiredIdentifierSchema(""),
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
				MarkdownDescription: "User who performed this action",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "Changed date-time",
				Computed:            true,
			},
			"updated_by": schema.StringAttribute{
				MarkdownDescription: "User who performed this action",
				Computed:            true,
			},
		},
	}
}

func SecretSchema() resourceschema.Schema {
	return resourceschema.Schema{
		MarkdownDescription: strings.Join([]string{
			"Secret resource",
			"The platform never returns the value of a secret, and Terraform does not store it " +
				"either (write-only attribute, requires Terraform 1.11 or later).",
			"The provider detects a changed `value` through `value_checksum`, and a secret " +
				"modified outside of Terraform through `updated_at`: in both cases the next plan " +
				"proposes to overwrite the secret with your configuration.",
		}, "\n"),
		Attributes: map[string]resourceschema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"name": resourceschema.StringAttribute{
				MarkdownDescription: "Secret name",
				Required:            true,
			},
			"value": resourceschema.StringAttribute{
				MarkdownDescription: strings.Join([]string{
					"Secret value (write-only, never stored in the plan or the state).",
					"Changing it updates the secret, see `value_checksum`.",
				}, "\n"),
				Required:  true,
				Sensitive: true,
				WriteOnly: true,
			},
			"value_checksum": resourceschema.StringAttribute{
				MarkdownDescription: strings.Join([]string{
					"SHA-256 of the `value` last written by Terraform, managed by the provider.",
					"A different checksum, or none (resource imported, secret modified outside of " +
						"Terraform) plans an update of the secret.",
					"Beware that the checksum of a short or guessable value can be brute-forced.",
				}, "\n"),
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseNonNullStateForUnknown(),
				},
			},
			"description": resourceschema.StringAttribute{
				MarkdownDescription: "Describe the resource in few sentences",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"org_id": ComputedOrganizationIdSchema(),
			"org_scoped": resourceschema.BoolAttribute{
				MarkdownDescription: strings.Join([]string{
					"Scoped to the organization?",
					"Either this or `project_ids` must be set.",
					"Default value is false.",
				}, "\n"),
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"project_ids": resourceschema.SetAttribute{
				MarkdownDescription: "Restrict to given projects (an empty set means none)",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Default: setdefault.StaticValue(
					types.SetValueMust(types.StringType, []attr.Value{})),
			},
			"created_at": resourceschema.StringAttribute{
				MarkdownDescription: "Creation date-time",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseNonNullStateForUnknown(),
				},
			},
			"created_by": resourceschema.StringAttribute{
				MarkdownDescription: "User who performed this action",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseNonNullStateForUnknown(),
				},
			},
			// No UseStateForUnknown below: these change on every update.
			"updated_at": resourceschema.StringAttribute{
				MarkdownDescription: strings.Join([]string{
					"Changed date-time",
					"A value different from the one stored after the last apply means the secret " +
						"was modified outside of Terraform, the next plan proposes to overwrite it.",
				}, "\n"),
				Computed: true,
			},
			"updated_by": resourceschema.StringAttribute{
				MarkdownDescription: "User who performed this action",
				Computed:            true,
			},
		},
	}
}

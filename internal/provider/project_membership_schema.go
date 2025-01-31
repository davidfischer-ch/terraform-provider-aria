// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

/* import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Named ProjectPrincipalsAssignment in Project's API Swagger
func ProjectMembershipSchema() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		MarkdownDescription: "A project membership granted to a user or group",
		Attributes: map[string]schema.Attribute{
			"email": schema.StringAttribute{
				MarkdownDescription: strings.Join(",", []string{
					"example: administrator@vmware.com",
					"The username of the user or display name of the group. "+
					"When assigning a group, the email is expected to have the format "+
					"displayName@domain. In the case where the display name in Identity provider "+
					"is in the format:",
					"* name@domain - email should be written as name@domain@domain",
					"* name (and group has domain) - email should be written as name@domain",
					"* name (and group doesn't have domain) - email should be written as name@",
					"",
					"to ensure proper functioning." + IMMUTABLE,
				}),
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Principal type, either user or group",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"user", "group"}...),
				},
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "Access level, one of administrator, member, supervisor or "+
					" viewer",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"administrator",
						"member",
						"supervisor",
						"viewer",
					}...),
				},
			},
		},
	}
} */

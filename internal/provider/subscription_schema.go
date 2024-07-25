// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func SubscriptionSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Subscription resource ([event broker API]" +
			"(https://developer.broadcom.com/xapis/vrealize-automation-event-broker-service-api/" +
			"latest/subscription/))",
		Attributes: map[string]schema.Attribute{
			"id": ComputedIdentifierSchema(""),
			"name": schema.StringAttribute{
				MarkdownDescription: "Subscription name",
				Required:            true,
			},
			"description": RequiredDescriptionSchema(),
			"type": schema.StringAttribute{
				MarkdownDescription: "Subscription type, either RUNNABLE or SUBSCRIBABLE",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"RUNNABLE", "SUBSCRIBABLE"}...),
				},
			},
			"runnable_type": schema.StringAttribute{
				MarkdownDescription: "Runnable type, either extensibility.abx or extensibility.vro",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"extensibility.abx", "extensibility.vro"}...),
				},
			},
			"runnable_id": schema.StringAttribute{
				MarkdownDescription: "Runnable identifier",
				Required:            true,
			},
			"recover_runnable_type": schema.StringAttribute{
				MarkdownDescription: "Recovery runnable type, either extensibility.abx or " +
					"extensibility.vro",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"extensibility.abx", "extensibility.vro"}...),
				},
			},
			"recover_runnable_id": schema.StringAttribute{
				MarkdownDescription: "Recovery runnable identifier",
				Optional:            true,
			},
			"event_topic_id": schema.StringAttribute{
				MarkdownDescription: "Event topic identifier",
				Required:            true,
			},
			"project_ids": schema.SetAttribute{
				MarkdownDescription: "Restrict to given projects (an empty list means all)",
				ElementType:         types.StringType,
				Required:            true,
			},
			"blocking": schema.BoolAttribute{
				MarkdownDescription: "TODO",
				Required:            true,
			},
			"broadcast": schema.BoolAttribute{
				MarkdownDescription: "TODO",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"contextual": schema.BoolAttribute{
				MarkdownDescription: "TODO",
				Required:            true,
			},
			"criteria": schema.StringAttribute{
				MarkdownDescription: "TODO",
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(""),
			},
			"disabled": schema.BoolAttribute{
				MarkdownDescription: "TODO",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "TODO",
				Required:            true,
			},
			"system": schema.BoolAttribute{
				MarkdownDescription: "TODO",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"timeout": schema.Int64Attribute{
				MarkdownDescription: "TODO",
				Required:            true,
			},
			"org_id": ComputedOrganizationIdSchema(),
			"owner_id": schema.StringAttribute{
				MarkdownDescription: "Owner identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"subscriber_id": schema.StringAttribute{
				MarkdownDescription: "Subscriber identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

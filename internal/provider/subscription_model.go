// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SubscriptionModel describes the resource data model.
type SubscriptionModel struct {
	Id                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	Type                types.String `tfsdk:"type"`
	RunnableType        types.String `tfsdk:"runnable_type"`
	RunnableId          types.String `tfsdk:"runnable_id"`
	RecoverRunnableType types.String `tfsdk:"recover_runnable_type"`
	RecoverRunnableId   types.String `tfsdk:"recover_runnable_id"`
	EventTopicId        types.String `tfsdk:"event_topic_id"`

	ProjectIds types.Set `tfsdk:"project_ids"`

	Blocking   types.Bool   `tfsdk:"blocking"`
	Broadcast  types.Bool   `tfsdk:"broadcast"`
	Contextual types.Bool   `tfsdk:"contextual"`
	Criteria   types.String `tfsdk:"criteria"`
	Disabled   types.Bool   `tfsdk:"disabled"`
	Priority   types.Int64  `tfsdk:"priority"`
	System     types.Bool   `tfsdk:"system"`
	Timeout    types.Int64  `tfsdk:"timeout"`

	OrgId        types.String `tfsdk:"org_id"`
	OwnerId      types.String `tfsdk:"owner_id"`
	SubscriberId types.String `tfsdk:"subscriber_id"`
}

// SubscriptionAPIModel describes the resource API model.
type SubscriptionAPIModel struct {
	Id                  string `json:"id,omitempty"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	Type                string `json:"type"`
	RunnableType        string `json:"runnableType"`
	RunnableId          string `json:"runnableId"`
	RecoverRunnableType string `json:"recoverRunnableType"`
	RecoverRunnableId   string `json:"recoverRunnableId"`
	EventTopicId        string `json:"eventTopicId"`

	Constraints map[string][]string `json:"constraints"`

	Blocking   bool   `json:"blocking"`
	Broadcast  bool   `json:"broadcast"`
	Contextual bool   `json:"contextual"`
	Criteria   string `json:"criteria"`
	Disabled   bool   `json:"disabled"`
	Priority   int64  `json:"priority"`
	System     bool   `json:"system"`
	Timeout    int64  `json:"timeout"`

	OrgId        string `json:"orgId"`
	OwnerId      string `json:"ownerId"`
	SubscriberId string `json:"subscriberId"`
}

func (self SubscriptionModel) String() string {
	return fmt.Sprintf(
		"Subscription %s (%s)",
		self.Id.ValueString(),
		self.Name.ValueString())
}

func (self *SubscriptionModel) GenerateId() {
	if len(self.Id.ValueString()) == 0 {
		self.Id = types.StringValue(uuid.New().String())
	}
}

func (self *SubscriptionModel) FromAPI(
	ctx context.Context,
	raw SubscriptionAPIModel,
) diag.Diagnostics {
	projectIds, diags := types.SetValueFrom(ctx, types.StringType, raw.Constraints["projectId"])

	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Description = types.StringValue(raw.Description)
	self.Type = types.StringValue(raw.Type)
	self.RunnableType = types.StringValue(raw.RunnableType)
	self.RunnableId = types.StringValue(raw.RunnableId)
	self.RecoverRunnableType = StringOrNullValue(raw.RecoverRunnableType)
	self.RecoverRunnableId = StringOrNullValue(raw.RecoverRunnableId)
	self.EventTopicId = types.StringValue(raw.EventTopicId)
	self.ProjectIds = projectIds
	self.Blocking = types.BoolValue(raw.Blocking)
	self.Broadcast = types.BoolValue(raw.Broadcast)
	self.Contextual = types.BoolValue(raw.Contextual)
	self.Criteria = types.StringValue(raw.Criteria)
	self.Disabled = types.BoolValue(raw.Disabled)
	self.Priority = types.Int64Value(raw.Priority)
	self.System = types.BoolValue(raw.System)
	self.Timeout = types.Int64Value(raw.Timeout)
	self.OrgId = types.StringValue(raw.OrgId)
	self.OwnerId = types.StringValue(raw.OwnerId)
	self.SubscriberId = types.StringValue(raw.SubscriberId)

	return diags
}

func (self SubscriptionModel) ToAPI(
	ctx context.Context,
) (SubscriptionAPIModel, diag.Diagnostics) {

	var diags diag.Diagnostics

	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/set
	if self.ProjectIds.IsNull() || self.ProjectIds.IsUnknown() {
		diags.AddError(
			"Configuration error",
			fmt.Sprintf(
				"Unable to manage subscription %s, project_ids is either null or unknown",
				self.Id.ValueString()))
		return SubscriptionAPIModel{}, diags
	}

	projectIds := make([]string, 0, len(self.ProjectIds.Elements()))
	diags = self.ProjectIds.ElementsAs(ctx, &projectIds, false)
	if diags.HasError() {
		return SubscriptionAPIModel{}, diags
	}

	return SubscriptionAPIModel{
		Id:                  self.Id.ValueString(),
		Name:                self.Name.ValueString(),
		Description:         self.Description.ValueString(),
		Type:                self.Type.ValueString(),
		RunnableType:        self.RunnableType.ValueString(),
		RunnableId:          self.RunnableId.ValueString(),
		RecoverRunnableType: self.RecoverRunnableType.ValueString(),
		RecoverRunnableId:   self.RecoverRunnableId.ValueString(),
		EventTopicId:        self.EventTopicId.ValueString(),
		Constraints:         map[string][]string{"projectId": projectIds},
		Blocking:            self.Blocking.ValueBool(),
		Contextual:          self.Contextual.ValueBool(),
		Criteria:            self.Criteria.ValueString(),
		Disabled:            self.Disabled.ValueBool(),
		Priority:            self.Priority.ValueInt64(),
		Timeout:             self.Timeout.ValueInt64(),
	}, diags
}

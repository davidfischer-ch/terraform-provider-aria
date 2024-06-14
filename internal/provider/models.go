// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ABXConstantModel describes the resource data model.
type ABXConstantModel struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Value     types.String `tfsdk:"value"`
	Encrypted types.Bool   `tfsdk:"encrypted"`
	OrgId     types.String `tfsdk:"org_id"`
}

// ABXConstantAPIModel describes the resource API model.
type ABXConstantAPIModel struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Value         string `json:"value"`
	Encrypted     bool   `json:"encrypted"`
	OrgId         string `json:"orgId"`
	CreatedMillis uint64 `json:"createdMillis"`
}

func (self *ABXConstantModel) FromAPI(
	ctx context.Context,
	raw ABXConstantAPIModel,
) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Value = types.StringValue(raw.Value)
	self.Encrypted = types.BoolValue(raw.Encrypted)
	self.OrgId = types.StringValue(raw.OrgId)
	return diag.Diagnostics{}
}

func (self *ABXConstantModel) ToAPI() ABXConstantAPIModel {
	return ABXConstantAPIModel{
		Name:      self.Name.ValueString(),
		Value:     self.Value.ValueString(),
		Encrypted: self.Encrypted.ValueBool(),
	}
}

// ABXSecretModel describes the resource data model.
type ABXSecretModel struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Value     types.String `tfsdk:"value"`
	Encrypted types.Bool   `tfsdk:"encrypted"`
	OrgId     types.String `tfsdk:"org_id"`
}

// ABXSecretAPIModel describes the resource API model.
type ABXSecretAPIModel struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Value         string `json:"value"`
	Encrypted     bool   `json:"encrypted"`
	OrgId         string `json:"orgId"`
	CreatedMillis uint64 `json:"createdMillis"`
}

func (self *ABXSecretModel) FromAPI(
	ctx context.Context,
	raw ABXSecretAPIModel,
) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	// The value is returned with the following '*****' awesome :)
	// self.Value = types.StringValue(raw.Value)
	self.Encrypted = types.BoolValue(raw.Encrypted)
	self.OrgId = types.StringValue(raw.OrgId)
	return diag.Diagnostics{}
}

func (self *ABXSecretModel) ToAPI() ABXSecretAPIModel {
	return ABXSecretAPIModel{
		Name:      self.Name.ValueString(),
		Value:     self.Value.ValueString(),
		Encrypted: self.Encrypted.ValueBool(),
	}
}

// CatalogTypeModel describes the catalog type model.
type CatalogTypeModel struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	BaseURI   types.String `tfsdk:"base_uri"`
	CreatedAt types.String `tfsdk:"created_at"`
	CreatedBy types.String `tfsdk:"created_by"`
	IconId    types.String `tfsdk:"icon_id"`
}

// CatalogTypeAPIModel describes the catalog type API model.
type CatalogTypeAPIModel struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	BaseURI   string `json:"baseUri"`
	CreatedAt string `json:"createdAt"`
	CreatedBy string `json:"createdBy"`
	IconId    string `json:"iconId"`
}

func (self *CatalogTypeModel) FromAPI(
	ctx context.Context,
	raw CatalogTypeAPIModel,
) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.BaseURI = types.StringValue(raw.BaseURI)
	self.CreatedAt = types.StringValue(raw.CreatedAt)
	self.CreatedBy = types.StringValue(raw.CreatedBy)
	self.IconId = types.StringValue(raw.IconId)
	return diag.Diagnostics{}
}

// IconModel describes the resource data model.
type IconModel struct {
	Id      types.String `tfsdk:"id"`
	Content types.String `tfsdk:"content"`
}

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

	/*ProjectIds types.Set `tfsdk:"project_ids"`*/

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
	Id                  string `json:"id"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	Type                string `json:"type"`
	RunnableType        string `json:"runnableType"`
	RunnableId          string `json:"runnableId"`
	RecoverRunnableType string `json:"recoverRunnableType"`
	RecoverRunnableId   string `json:"recoverRunnableId"`
	EventTopicId        string `json:"eventTopicId"`

	/*Constraints map[string][]string `json:"constraints"`*/

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

func (self *SubscriptionModel) GenerateId() {
	if len(self.Id.ValueString()) == 0 {
		self.Id = types.StringValue(fmt.Sprintf("sub_%d", time.Now().UnixMilli()))
	}
}

func (self *SubscriptionModel) FromAPI(
	ctx context.Context,
	raw SubscriptionAPIModel,
) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Description = types.StringValue(raw.Description)
	self.Type = types.StringValue(raw.Type)
	self.RunnableType = types.StringValue(raw.RunnableType)
	self.RunnableId = types.StringValue(raw.RunnableId)
	self.RecoverRunnableType = StringOrNullValue(raw.RecoverRunnableType)
	self.RecoverRunnableId = StringOrNullValue(raw.RecoverRunnableId)
	self.EventTopicId = types.StringValue(raw.EventTopicId)
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

	/*self.projectIds, diags := types.SetValueFrom(
	ctx, types.StringType, raw.Constraints["projectId"])*/

	return diag.Diagnostics{}
}

func (self *SubscriptionModel) ToAPI() SubscriptionAPIModel {
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
		Blocking:            self.Blocking.ValueBool(),
		Contextual:          self.Contextual.ValueBool(),
		Criteria:            self.Criteria.ValueString(),
		Disabled:            self.Disabled.ValueBool(),
		Priority:            self.Priority.ValueInt64(),
		Timeout:             self.Timeout.ValueInt64(),
	}
}

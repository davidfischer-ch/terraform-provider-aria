// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CatalogItemModel describes the resource data model.
type CatalogItemModel struct {
	Id          types.String         `tfsdk:"id"`
	Name        types.String         `tfsdk:"name"`
	Description types.String         `tfsdk:"description"`
	Schema      jsontypes.Normalized `tfsdk:"schema"`

	ExternalId types.String `tfsdk:"external_id"`
	FormId     types.String `tfsdk:"form_id"`
	IconId     types.String `tfsdk:"icon_id"`
	TypeId     types.String `tfsdk:"type_id"`

	SourceId   types.String `tfsdk:"source_id"`
	SourceName types.String `tfsdk:"source_name"`

	CreatedAt     timetypes.RFC3339 `tfsdk:"created_at"`
	CreatedBy     types.String      `tfsdk:"created_by"`
	LastUpdatedAt timetypes.RFC3339 `tfsdk:"last_updated_at"`
	LastUpdatedBy types.String      `tfsdk:"last_updated_by"`
}

// CatalogItemAPIModel describes the resource API model.
type CatalogItemAPIModel struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Schema      any    `json:"schema"`

	ExternalId string `json:"externalId"`
	FormId     string `json:"formId"`
	IconId     string `json:"iconId"`

	Type CatalogItemTypeAPIModel `json:"type"`

	SourceId   string `json:"sourceId"`
	SourceName string `json:"sourceName"`

	CreatedAt     string `json:"createdAt,omitempty"`
	CreatedBy     string `json:"createdBy,omitempty"`
	LastUpdatedAt string `json:"lastUpdatedAt,omitempty"`
	LastUpdatedBy string `json:"lastUpdatedBy,omitempty"`
}

type CatalogItemTypeAPIModel struct {
	Id   string `json:"id"`
	Link string `json:"link"`
	Name string `json:"name"`
}

type CatalogItemLstAPIModel struct {
	Content          []CatalogItemAPIModel `json:"Content"`
	TotalElements    int                   `json:"totalElements"`
	NumberOfElements int                   `json:"numberOfElements"`
}

func (self CatalogItemModel) String() string {
	return fmt.Sprintf(
		"Catalog Item %s (%s)",
		self.Id.ValueString(),
		self.Name.ValueString())
}

func (self CatalogItemModel) ListPath() string {
	return "catalog/api/admin/items"
}

func (self CatalogItemModel) ReadPath() string {
	return "catalog/api/admin/items/" + self.Id.ValueString()
}

func (self *CatalogItemModel) FromAPI(raw CatalogItemAPIModel) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Description = types.StringValue(raw.Description)
	self.ExternalId = types.StringValue(raw.ExternalId)
	self.IconId = types.StringValue(raw.IconId)
	self.FormId = types.StringValue(raw.FormId)
	self.TypeId = types.StringValue(raw.Type.Id)
	self.SourceId = types.StringValue(raw.SourceId)
	self.SourceName = types.StringValue(raw.SourceName)
	self.CreatedBy = types.StringValue(raw.CreatedBy)
	self.LastUpdatedBy = types.StringValue(raw.LastUpdatedBy)

	diags := diag.Diagnostics{}

	var someDiags diag.Diagnostics

	self.Schema, someDiags = JSONNormalizedFromAny(self.String(), raw.Schema)
	diags.Append(someDiags...)

	self.CreatedAt, someDiags = timetypes.NewRFC3339Value(raw.CreatedAt)
	diags.Append(someDiags...)

	self.LastUpdatedAt, someDiags = timetypes.NewRFC3339Value(raw.LastUpdatedAt)
	diags.Append(someDiags...)

	return diags
}

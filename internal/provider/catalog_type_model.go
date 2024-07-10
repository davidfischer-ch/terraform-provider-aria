// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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
	Id        string `json:"id,omitempty"`
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

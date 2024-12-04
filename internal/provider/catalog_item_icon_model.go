// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CatalogItemIconModel describes the resource data model.
type CatalogItemIconModel struct {
	Id     types.String `tfsdk:"item_id"`
	IconId types.String `tfsdk:"icon_id"`
}

// CatalogItemIconAPIModel describes the resource API model.
type CatalogItemIconAPIModel struct {
	Id     string `json:"id"`
	IconId string `json:"iconId"`
}

func (self CatalogItemIconModel) String() string {
	return fmt.Sprintf(
		"Catalog Item %s Icon %s",
		self.Id.ValueString(),
		self.IconId.ValueString())
}

// Return an appropriate key that can be used for naming mutexes.
// Create: Identifier can be used to prevent concurrent creation of catalog items.
// Read Update Delete: Identifier can be used to prevent concurrent modifications on the instance.
func (self CatalogItemIconModel) LockKey() string {
	return "catalog-item-" + self.Id.ValueString()
}

func (self CatalogItemIconModel) CreatePath() string {
	return "catalog/api/admin/items/" + self.Id.ValueString()
}

func (self CatalogItemIconModel) ReadPath() string {
	return self.CreatePath()
}

func (self CatalogItemIconModel) UpdatePath() string {
	return self.ReadPath()
}

func (self CatalogItemIconModel) DeletePath() string {
	return self.ReadPath() // Even if not possible ...
}

func (self *CatalogItemIconModel) FromAPI(
	ctx context.Context,
	raw CatalogItemIconAPIModel,
) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	self.IconId = types.StringValue(raw.IconId)
	return diag.Diagnostics{}
}

func (self CatalogItemIconModel) ToAPI(
	ctx context.Context,
) (CatalogItemIconAPIModel, diag.Diagnostics) {
	return CatalogItemIconAPIModel{
		Id: self.Id.ValueString(),
		IconId: self.IconId.ValueString(),
	}, diag.Diagnostics{}
}

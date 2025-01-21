// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CatalogSourceModel describes the resource data model.
type CatalogSourceModel struct {
	Id     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	TypeId types.String `tfsdk:"type_id"`
	Global types.Bool   `tfsdk:"global"`

	Config CatalogSourceConfigModel `tfsdk:"config"`

	CreatedAt timetypes.RFC3339 `tfsdk:"created_at"`
	CreatedBy types.String      `tfsdk:"created_by"`
	/*LastUpdatedAt         timetypes.RFC3339 `tfsdk:"last_updated_at"`
	LastUpdatedBy         types.String      `tfsdk:"last_updated_by"`
	LastImportStartedAt   timetypes.RFC3339 `tfsdk:"last_import_started_at"`
	LastImportCompletedAt timetypes.RFC3339 `tfsdk:"last_import_completed_at"`
	LastImportErrors      types.List        `tfsdk:"last_import_errors"`*/

	ItemsImported types.Int32 `tfsdk:"items_imported"`
	ItemsFound    types.Int32 `tfsdk:"items_found"`
}

// CatalogSourceAPIModel describes the resource API model.
type CatalogSourceAPIModel struct {
	Id     string `tfsdk:"id,omitempty"`
	Name   string `tfsdk:"name"`
	TypeId string `tfsdk:"typeId"`
	Global *bool  `tfsdk:"global,omitempty"`

	Config CatalogSourceConfigAPIModel `json:"config"`

	CreatedAt string `tfsdk:"createdAt,omitempty"`
	CreatedBy string `tfsdk:"createdBy,omitempty"`
	/*LastUpdatedAt         string   `tfsdk:"lastUpdatedAt,omitempty"`
	LastUpdatedBy         string   `tfsdk:"lastUpdatedBy,omitempty"`
	LastImportStartedAt   string   `tfsdk:"lastImportStartedAt,omitempty"`
	LastImportCompletedAt string   `tfsdk:"lastImportCompletedAt,omitempty"`*/
	LastImportErrors []string `tfsdk:"lastImportErrors,omitempty"`

	ItemsImported int32 `tfsdk:"itemsImported,omitempty"`
	ItemsFound    int32 `tfsdk:"itemsFound,omitempty"`
}

func (self CatalogSourceModel) String() string {
	return fmt.Sprintf(
		"Catalog Source %s (%s)",
		self.Id.ValueString(),
		self.Name.ValueString())
}

// Return an appropriate key that can be used for naming mutexes.
// Create: Identifier can be used to prevent concurrent creation of projects.
// Read Update Delete: Identifier can be used to prevent concurrent modifications on the instance.
func (self CatalogSourceModel) LockKey() string {
	return "catalog-source-" + self.Id.ValueString()
}

func (self CatalogSourceModel) CreatePath() string {
	return "catalog/api/admin/sources"
}

func (self CatalogSourceModel) ReadPath() string {
	return "catalog/api/admin/sources/" + self.Id.ValueString()
}

func (self CatalogSourceModel) UpdatePath() string {
	return self.CreatePath()
}

func (self CatalogSourceModel) DeletePath() string {
	return self.ReadPath()
}

func (self *CatalogSourceModel) FromAPI(
	ctx context.Context,
	raw CatalogSourceAPIModel,
) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.TypeId = types.StringValue(raw.TypeId)
	self.Global = types.BoolValue(*raw.Global)
	self.CreatedBy = types.StringValue(raw.CreatedBy)
	/*self.LastUpdatedBy = types.StringValue(raw.LastUpdatedBy)*/
	self.ItemsImported = types.Int32Value(raw.ItemsImported)
	self.ItemsFound = types.Int32Value(raw.ItemsFound)

	diags := self.Config.FromAPI(ctx, raw.Config)

	var timestampDiags diag.Diagnostics
	self.CreatedAt, timestampDiags = timetypes.NewRFC3339Value(raw.CreatedAt)
	diags.Append(timestampDiags...)

	/*dateTime, timeDiags = timetypes.NewRFC3339PointerValue(raw.LastUpdatedAt)
	diags.Append(timeDiags...)
	self.LastUpdatedAt = dateTime

	dateTime, timeDiags = timetypes.NewRFC3339PointerValue(raw.LastImportStartedAt)
	diags.Append(timeDiags...)
	self.LastImportStartedAt = dateTime

	dateTime, timeDiags = timetypes.NewRFC3339PointerValue(raw.LastImportCompletedAt)
	diags.Append(timeDiags...)
	self.LastImportCompletedAt = dateTime

	errors, errorsDiags := types.ListValueFrom(ctx, types.StringType, raw.LastImportErrors)
	self.LastImportErrors = errors
	diags.Append(errorsDiags...)*/

	return diags
}

func (self CatalogSourceModel) ToAPI(
	ctx context.Context,
) (CatalogSourceAPIModel, diag.Diagnostics) {
	configRaw, diags := self.Config.ToAPI(ctx, self.String())
	return CatalogSourceAPIModel{
		Id:     self.Id.ValueString(),
		Name:   self.Name.ValueString(),
		TypeId: self.TypeId.ValueString(),
		Config: configRaw,
	}, diags
}

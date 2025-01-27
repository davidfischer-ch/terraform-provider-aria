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

	CreatedAt             timetypes.RFC3339 `tfsdk:"created_at"`
	CreatedBy             types.String      `tfsdk:"created_by"`
	LastUpdatedAt         timetypes.RFC3339 `tfsdk:"last_updated_at"`
	LastUpdatedBy         types.String      `tfsdk:"last_updated_by"`
	LastImportStartedAt   timetypes.RFC3339 `tfsdk:"last_import_started_at"`
	LastImportCompletedAt timetypes.RFC3339 `tfsdk:"last_import_completed_at"`
	/*LastImportErrors      types.List        `tfsdk:"last_import_errors"`*/

	ItemsImported types.Int32 `tfsdk:"items_imported"`
	ItemsFound    types.Int32 `tfsdk:"items_found"`
}

// CatalogSourceAPIModel describes the resource API model.
type CatalogSourceAPIModel struct {
	Id     string `json:"id,omitempty"`
	Name   string `json:"name"`
	TypeId string `json:"typeId"`
	Global bool   `json:"global,omitempty"`

	Config CatalogSourceConfigAPIModel `json:"config"`

	CreatedAt             string   `json:"createdAt,omitempty"`
	CreatedBy             string   `json:"createdBy,omitempty"`
	LastUpdatedAt         string   `json:"lastUpdatedAt,omitempty"`
	LastUpdatedBy         string   `json:"lastUpdatedBy,omitempty"`
	LastImportStartedAt   string   `json:"lastImportStartedAt,omitempty"`
	LastImportCompletedAt string   `json:"lastImportCompletedAt,omitempty"`
	LastImportErrors      []string `json:"lastImportErrors,omitempty"`

	ItemsImported int32 `json:"itemsImported,omitempty"`
	ItemsFound    int32 `json:"itemsFound,omitempty"`
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
	self.Global = types.BoolValue(raw.Global)
	self.CreatedBy = types.StringValue(raw.CreatedBy)
	self.LastUpdatedBy = types.StringValue(raw.LastUpdatedBy)
	self.ItemsImported = types.Int32Value(raw.ItemsImported)
	self.ItemsFound = types.Int32Value(raw.ItemsFound)

	diags := self.Config.FromAPI(ctx, raw.Config)

	var timeDiags diag.Diagnostics

	self.CreatedAt, timeDiags = timetypes.NewRFC3339Value(raw.CreatedAt)
	diags.Append(timeDiags...)

	self.LastUpdatedAt, timeDiags = timetypes.NewRFC3339Value(raw.LastUpdatedAt)
	diags.Append(timeDiags...)

	self.LastImportStartedAt, timeDiags = timetypes.NewRFC3339Value(raw.LastImportStartedAt)
	diags.Append(timeDiags...)

	if len(raw.LastImportCompletedAt) == 0 {
		self.LastImportCompletedAt = timetypes.NewRFC3339Null()
	} else {
		self.LastImportCompletedAt, timeDiags = timetypes.NewRFC3339Value(raw.LastImportCompletedAt)
		diags.Append(timeDiags...)
	}

	/*errors, errorsDiags := types.ListValueFrom(ctx, types.StringType, raw.LastImportErrors)
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

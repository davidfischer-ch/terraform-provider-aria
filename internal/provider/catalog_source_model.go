// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	LastImportErrors      types.List        `tfsdk:"last_import_errors"`

	ItemsImported types.Int32 `tfsdk:"items_imported"`
	ItemsFound    types.Int32 `tfsdk:"items_found"`

	WaitImported types.Bool `tfsdk:"wait_imported"`
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

	var someDiags diag.Diagnostics

	self.CreatedAt, someDiags = timetypes.NewRFC3339Value(raw.CreatedAt)
	diags.Append(someDiags...)

	self.LastUpdatedAt, someDiags = timetypes.NewRFC3339Value(raw.LastUpdatedAt)
	diags.Append(someDiags...)

	self.LastImportStartedAt, someDiags = timetypes.NewRFC3339Value(raw.LastImportStartedAt)
	diags.Append(someDiags...)

	if len(raw.LastImportCompletedAt) == 0 {
		self.LastImportCompletedAt = timetypes.NewRFC3339Null()
	} else {
		self.LastImportCompletedAt, someDiags = timetypes.NewRFC3339Value(raw.LastImportCompletedAt)
		diags.Append(someDiags...)
	}

	self.LastImportErrors, someDiags = types.ListValueFrom(
		ctx, types.StringType, raw.LastImportErrors,
	)
	diags.Append(someDiags...)

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

// Utils -------------------------------------------------------------------------------------------

func (self CatalogSourceModel) IsImporting(ctx context.Context) bool {
	startedAt, startedDiags := self.LastImportStartedAt.ValueRFC3339Time()
	completedAt, completedDiags := self.LastImportCompletedAt.ValueRFC3339Time()
	tflog.Debug(
		ctx,
		fmt.Sprintf(
			"Resource %s last_import_started_at=%s last_import_completed_at=%s",
			self.String(), startedAt.String(), completedAt.String()))

	// Is not importing since not started
	if startedDiags.HasError() {
		return false
	}

	// Is importing since not completed
	if completedDiags.HasError() {
		return true
	}

	return startedAt.After(completedAt)
}

// Return a tuple with waitAndSee, errors and diagnostics.
// If some errors may be fixed by the next integration's refresh process then waitAndSee is true.
func (self CatalogSourceModel) QualifyErrors(
	ctx context.Context,
) (bool, []string, diag.Diagnostics) {

	diags := diag.Diagnostics{}
	errors := make([]string, 0, len(self.LastImportErrors.Elements()))

	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/list
	if self.LastImportErrors.IsNull() || self.LastImportErrors.IsUnknown() {
		diags.AddError(
			"Configuration error",
			fmt.Sprintf(
				"Unable to qualify %s errors, last_import_errors is either null or unknown",
				self.String()))
		return false, errors, diags
	}

	diags.Append(self.LastImportErrors.ElementsAs(ctx, &errors, false)...)

	for _, error := range errors {
		// Next integration's refresh process may fix this issue
		if strings.Contains(error, "Error downloading catalog item") {
			return true, errors, diags
		}
	}
	return false, errors, diags
}

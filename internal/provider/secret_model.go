// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SecretModel describes the resource data model.
type SecretModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	// Write-only, always null in the plan and the state. Read it from the configuration.
	Value         types.String `tfsdk:"value"`
	ValueChecksum types.String `tfsdk:"value_checksum"`
	Description   types.String `tfsdk:"description"`
	OrgId         types.String `tfsdk:"org_id"`
	OrgScoped     types.Bool   `tfsdk:"org_scoped"`
	ProjectIds    types.Set    `tfsdk:"project_ids"`
	CreatedAt     types.String `tfsdk:"created_at"`
	CreatedBy     types.String `tfsdk:"created_by"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
	UpdatedBy     types.String `tfsdk:"updated_by"`
}

// SecretDataSourceModel describes the data source data model.
type SecretDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	OrgId       types.String `tfsdk:"org_id"`
	OrgScoped   types.Bool   `tfsdk:"org_scoped"`
	ProjectIds  types.Set    `tfsdk:"project_ids"`
	CreatedAt   types.String `tfsdk:"created_at"`
	CreatedBy   types.String `tfsdk:"created_by"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
	UpdatedBy   types.String `tfsdk:"updated_by"`
}

// SecretAPIModel describes the secret API model, as the platform returns it.
//
// The value is returned masked (`********`), it is deliberately not part of this model.
type SecretAPIModel struct {
	Id          string   `json:"id,omitempty"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	OrgId       string   `json:"orgId"`
	OrgScoped   bool     `json:"orgScoped"`
	ProjectIds  []string `json:"projectIds"`
	CreatedAt   string   `json:"createdAt"`
	CreatedBy   string   `json:"createdBy"`
	UpdatedAt   string   `json:"updatedAt"`
	UpdatedBy   string   `json:"updatedBy"`
}

// SecretWriteAPIModel describes what the platform accepts to create or update a secret.
//
// Every field is sent on update too: the platform drops `orgScoped` from a secret when a PATCH
// leaves it out.
type SecretWriteAPIModel struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Value       string   `json:"value"`
	OrgScoped   bool     `json:"orgScoped"`
	ProjectIds  []string `json:"projectIds"`
}

// SecretValueChecksum returns the checksum stored in place of the (write-only) value.
func SecretValueChecksum(value string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(value)))
}

// SecretValueChanged tells if the configured value differs from what was last written.
//
// A missing checksum (resource imported, or secret modified outside of Terraform) counts as a
// change, the platform never discloses the value so it cannot be verified.
func SecretValueChanged(configValue types.String, stateChecksum types.String) bool {
	if configValue.IsUnknown() {
		return true
	}
	if stateChecksum.IsNull() || stateChecksum.IsUnknown() {
		return true
	}
	return SecretValueChecksum(configValue.ValueString()) != stateChecksum.ValueString()
}

// SecretModifiedOutside tells if the platform's `updated_at` differs from the one we stored.
//
// The stored one comes from the read following our last write, so any difference means that
// somebody else modified the secret since. Nothing to compare when there is no previous value
// (import).
func SecretModifiedOutside(previousUpdatedAt types.String, currentUpdatedAt string) bool {
	if previousUpdatedAt.IsNull() || previousUpdatedAt.IsUnknown() {
		return false
	}
	return previousUpdatedAt.ValueString() != currentUpdatedAt
}

func projectIdsOrEmpty(projectIds []string) []string {
	if projectIds == nil {
		return []string{}
	}
	return projectIds
}

func (self SecretModel) String() string {
	return fmt.Sprintf(
		"Secret %s (%s)",
		self.Id.ValueString(),
		self.Name.ValueString())
}

func (self SecretModel) LockKey() string {
	return "secret-" + self.Id.ValueString()
}

func (self SecretModel) CreatePath() string {
	return "platform/api/secrets"
}

func (self SecretModel) ReadPath() string {
	return "platform/api/secrets/" + self.Id.ValueString()
}

func (self SecretModel) UpdatePath() string {
	return self.ReadPath()
}

func (self SecretModel) DeletePath() string {
	return self.ReadPath()
}

func (self SecretModel) ToAPI(
	ctx context.Context,
	value string,
) (SecretWriteAPIModel, diag.Diagnostics) {
	projectIds := []string{}
	diags := self.ProjectIds.ElementsAs(ctx, &projectIds, false)
	return SecretWriteAPIModel{
		Name:        self.Name.ValueString(),
		Description: self.Description.ValueString(),
		Value:       value,
		OrgScoped:   self.OrgScoped.ValueBool(),
		ProjectIds:  projectIds,
	}, diags
}

func (self *SecretModel) FromAPI(ctx context.Context, raw SecretAPIModel) diag.Diagnostics {
	projectIds, diags := types.SetValueFrom(ctx, types.StringType, projectIdsOrEmpty(raw.ProjectIds))

	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Description = types.StringValue(raw.Description)
	self.OrgId = types.StringValue(raw.OrgId)
	self.OrgScoped = types.BoolValue(raw.OrgScoped)
	self.ProjectIds = projectIds
	self.CreatedAt = types.StringValue(raw.CreatedAt)
	self.CreatedBy = types.StringValue(raw.CreatedBy)
	self.UpdatedAt = types.StringValue(raw.UpdatedAt)
	self.UpdatedBy = types.StringValue(raw.UpdatedBy)

	return diags
}

// -------------------------------------------------------------------------------------------------

func (self SecretDataSourceModel) String() string {
	return fmt.Sprintf(
		"Secret %s (%s)",
		self.Id.ValueString(),
		self.Name.ValueString())
}

func (self SecretDataSourceModel) ReadPath() string {
	return "platform/api/secrets/" + self.Id.ValueString()
}

func (self *SecretDataSourceModel) FromAPI(
	ctx context.Context,
	raw SecretAPIModel,
) diag.Diagnostics {
	projectIds, diags := types.SetValueFrom(ctx, types.StringType, projectIdsOrEmpty(raw.ProjectIds))

	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Description = types.StringValue(raw.Description)
	self.OrgId = types.StringValue(raw.OrgId)
	self.OrgScoped = types.BoolValue(raw.OrgScoped)
	self.ProjectIds = projectIds
	self.CreatedAt = types.StringValue(raw.CreatedAt)
	self.CreatedBy = types.StringValue(raw.CreatedBy)
	self.UpdatedAt = types.StringValue(raw.UpdatedAt)
	self.UpdatedBy = types.StringValue(raw.UpdatedBy)

	return diags
}

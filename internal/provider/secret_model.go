// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SecretModel describes the secret model.
type SecretModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	// Value
	Description types.String `tfsdk:"description"`
	OrgId       types.String `tfsdk:"org_id"`
	OrgScoped   types.Bool   `tfsdk:"org_scoped"`
	ProjectIds  types.Set    `tfsdk:"project_ids"`
	CreatedAt   types.String `tfsdk:"created_at"`
	CreatedBy   types.String `tfsdk:"created_by"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
	UpdatedBy   types.String `tfsdk:"updated_by"`
}

// SecretAPIModel describes the secret API model.
type SecretAPIModel struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name"`
	// Value
	Description string   `json:"description"`
	OrgId       string   `json:"orgId"`
	OrgScoped   bool     `json:"orgScoped"`
	ProjectIds  []string `json:"projectIds"`
	CreatedAt   string   `json:"createdAt"`
	CreatedBy   string   `json:"createdBy"`
	UpdatedAt   string   `json:"updatedAt"`
	UpdatedBy   string   `json:"updatedBy"`
}

func (self SecretModel) String() string {
	return fmt.Sprintf(
		"Secret %s (%s)",
		self.Id.ValueString(),
		self.Name.ValueString())
}

func (self SecretModel) ReadPath() string {
	return "platform/api/secrets/" + self.Id.ValueString()
}

func (self *SecretModel) FromAPI(
	ctx context.Context,
	raw SecretAPIModel,
) diag.Diagnostics {
	projectIds, diags := types.SetValueFrom(ctx, types.StringType, raw.ProjectIds)

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

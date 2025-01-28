// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// PolicyModel describes the resource data model.
type PolicyModel struct {
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	EnforcementType types.String `tfsdk:"enforcement_type"`
	TypeId          types.String `tfsdk:"type_id"`

	Criteria      jsontypes.Normalized `tfsdk:"criteria"`
	ScopeCriteria jsontypes.Normalized `tfsdk:"scope_criteria"`
	Definition    jsontypes.Normalized `tfsdk:"definition"`

	CreatedAt     timetypes.RFC3339 `tfsdk:"created_at"`
	CreatedBy     types.String      `tfsdk:"created_by"`
	LastUpdatedAt timetypes.RFC3339 `tfsdk:"last_updated_at"`
	LastUpdatedBy types.String      `tfsdk:"last_updated_by"`

	ProjectId types.String `tfsdk:"project_id"`
	OrgId     types.String `tfsdk:"org_id"`
}

// PolicyAPIModel describes the resource API model.
type PolicyAPIModel struct {
	Id              string `json:"id,omitempty"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	EnforcementType string `json:"enforcementType"`
	TypeId          string `json:"typeId"`

	Criteria      any `json:"criteria,omitempty"`
	ScopeCriteria any `json:"scopeCriteria,omitempty"`
	Definition    any `json:"definition"`

	CreatedAt     string `json:"createdAt,omitempty"`
	CreatedBy     string `json:"createdBy,omitempty"`
	LastUpdatedAt string `json:"lastUpdatedAt,omitempty"`
	LastUpdatedBy string `json:"lastUpdatedBy,omitempty"`

	ProjectId string `json:"projectId,omitempty"`
	OrgId     string `json:"orgId,omitempty"`
}

func (self PolicyModel) String() string {
	return fmt.Sprintf(
		"%s Policy %s (%s)",
		self.GetType(),
		self.Id.ValueString(),
		self.Name.ValueString())
}

// Return an appropriate key that can be used for naming mutexes.
// Create: Identifier can be used to prevent concurrent creation of policies.
// Read Update Delete: Identifier can be used to prevent concurrent modifications on the instance.
func (self PolicyModel) LockKey() string {
	return "policy-" + self.Id.ValueString()
}

func (self PolicyModel) CreatePath() string {
	return "policy/api/policies"
}

func (self PolicyModel) ReadPath() string {
	return "policy/api/policies/" + self.Id.ValueString()
}

func (self PolicyModel) UpdatePath() string {
	return self.CreatePath()
}

func (self PolicyModel) DeletePath() string {
	return self.ReadPath()
}

func (self *PolicyModel) FromAPI(ctx context.Context, raw PolicyAPIModel) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Description = types.StringValue(raw.Description)
	self.EnforcementType = types.StringValue(raw.EnforcementType)
	self.TypeId = types.StringValue(raw.TypeId)
	self.CreatedBy = types.StringValue(raw.CreatedBy)
	self.LastUpdatedBy = types.StringValue(raw.LastUpdatedBy)
	self.ProjectId = types.StringValue(raw.ProjectId)
	self.OrgId = types.StringValue(raw.OrgId)

	diags := diag.Diagnostics{}
	var someDiags diag.Diagnostics

	self.Criteria, someDiags = JSONNormalizedFromAny(self.String(), raw.Criteria)
	diags.Append(someDiags...)

	self.ScopeCriteria, someDiags = JSONNormalizedFromAny(self.String(), raw.ScopeCriteria)
	diags.Append(someDiags...)

	self.Definition, someDiags = JSONNormalizedFromAny(self.String(), raw.Definition)
	diags.Append(someDiags...)

	self.CreatedAt, someDiags = timetypes.NewRFC3339Value(raw.CreatedAt)
	diags.Append(someDiags...)

	self.LastUpdatedAt, someDiags = timetypes.NewRFC3339Value(raw.LastUpdatedAt)
	diags.Append(someDiags...)

	return diags
}

func (self PolicyModel) ToAPI(ctx context.Context) (PolicyAPIModel, diag.Diagnostics) {

	// Criterias & Definition JSON Encoded -> API data

	criteriaRaw, diags := JSONNormalizedToAny(self.Criteria)

	scopeCriteriaRaw, someDiags := JSONNormalizedToAny(self.ScopeCriteria)
	diags.Append(someDiags...)

	definitionRaw, someDiags := JSONNormalizedToAny(self.Definition)
	diags.Append(someDiags...)

	return PolicyAPIModel{
		Id:              self.Id.ValueString(),
		Name:            self.Name.ValueString(),
		Description:     self.Description.ValueString(),
		EnforcementType: self.EnforcementType.ValueString(),
		TypeId:          self.TypeId.ValueString(),
		Criteria:        criteriaRaw,
		ScopeCriteria:   scopeCriteriaRaw,
		Definition:      definitionRaw,
		ProjectId:       self.ProjectId.ValueString(),
		OrgId:           self.OrgId.ValueString(),
	}, diags
}

// Utils -------------------------------------------------------------------------------------------

// Return policy type (e.g. Approval, Resource Quota, ...)
func (self PolicyModel) GetType() string {
	type_id := strings.TrimPrefix(self.TypeId.ValueString(), "com.vmware.policy.")
	return cases.Title(language.English).String(strings.ReplaceAll(type_id, ".", " "))
}

// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// OrchestratorTaskModel describes the resource data model.
type OrchestratorTaskModel struct {
	Id                  types.String      `tfsdk:"id"`
	Name                types.String      `tfsdk:"name"`
	Description         types.String      `tfsdk:"description"`
	Href                types.String      `tfsdk:"href"`
	RecurrenceCycle     types.String      `tfsdk:"recurrence_cycle"`
	RecurrencePattern   types.String      `tfsdk:"recurrence_pattern"`
	RecurrenceStartDate timetypes.RFC3339 `tfsdk:"recurrence_start_date"`
	RecurrenceEndDate   timetypes.RFC3339 `tfsdk:"recurrence_end_date"`
	RunningInstanceId   types.String      `tfsdk:"running_instance_id"`
	StartMode           types.String      `tfsdk:"start_mode"`
	State               types.String      `tfsdk:"state"`
	User                types.String      `tfsdk:"user"`

	/* InputParameters types.List `tfsdk:"input_parameters"` */

	Workflow types.Object `tfsdk:"workflow"`

	/*
		// Of type RelationModel
		Relations types.List `tfsdk:"relations"`
	*/
}

// OrchestratorTaskAPIModel describes the resource API model.
type OrchestratorTaskAPIModel struct {
	Id                  string `json:"id,omitempty"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	Href                string `json:"href,omitempty"`
	RecurrenceCycle     string `json:"recurrence-cycle"`
	RecurrencePattern   string `json:"recurrence-pattern"`
	RecurrenceStartDate string `json:"recurrence-start-date"`
	RecurrenceEndDate   string `json:"recurrence-end-date,omitempty"`
	RunningInstanceId   string `json:"running-instance-id,omitempty"`
	StartMode           string `json:"start-mode"`
	State               string `json:"state,omitempty"`
	User                string `json:"user,omitempty"`

	InputParameters []any                            `json:"input-parameters"`
	Workflow        OrchestratorTaskWorkflowAPIModel `json:"workflow"`

	/*
		Relations []RelationAPIModel `json:"relations",
	*/
}

func (self OrchestratorTaskModel) String() string {
	return fmt.Sprintf(
		"Orchestrator Task %s (%s) of %s",
		self.Id.ValueString(),
		self.Name.ValueString(),
		self.User.ValueString())
}

// Return an appropriate key that can be used for naming mutexes.
// Create: Identifier can be used to prevent concurrent creation of tasks.
// Read Update Delete: Identifier can be used to prevent concurrent modifications on the instance.
func (self OrchestratorTaskModel) LockKey() string {
	return "orchestrator-task-" + self.Id.ValueString()
}

func (self OrchestratorTaskModel) CreatePath() string {
	return "vco/api/tasks"
}

func (self OrchestratorTaskModel) ReadPath() string {
	return "vco/api/tasks/" + self.Id.ValueString()
}

func (self OrchestratorTaskModel) UpdatePath() string {
	return self.ReadPath()
}

func (self OrchestratorTaskModel) DeletePath() string {
	return self.ReadPath()
}

func (self *OrchestratorTaskModel) FromAPI(
	ctx context.Context,
	raw OrchestratorTaskAPIModel,
) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Description = types.StringValue(raw.Description)
	self.Href = types.StringValue(raw.Href)
	self.RecurrenceCycle = types.StringValue(raw.RecurrenceCycle)
	self.RecurrencePattern = types.StringValue(raw.RecurrencePattern)
	self.RunningInstanceId = types.StringValue(raw.RunningInstanceId)
	self.StartMode = types.StringValue(raw.StartMode)
	self.State = types.StringValue(raw.State)
	self.User = types.StringValue(raw.User)

	diags := diag.Diagnostics{}
	var someDiags diag.Diagnostics

	self.RecurrenceStartDate, someDiags = timetypes.NewRFC3339Value(raw.RecurrenceStartDate)
	diags.Append(someDiags...)

	if len(raw.RecurrenceEndDate) == 0 {
		self.RecurrenceEndDate = timetypes.NewRFC3339Null()
	} else {
		self.RecurrenceEndDate, someDiags = timetypes.NewRFC3339Value(raw.RecurrenceEndDate)
		diags.Append(someDiags...)
	}

	/* InputParameters = types.List */

	// Convert workflow from raw and then to object
	workflow := OrchestratorTaskWorkflowModel{}
	workflow.FromAPI(raw.Workflow)
	self.Workflow, someDiags = types.ObjectValueFrom(ctx, workflow.AttributeTypes(), workflow)
	diags.Append(someDiags...)

	return diags
}

func (self OrchestratorTaskModel) ToAPI(
	ctx context.Context,
) (OrchestratorTaskAPIModel, diag.Diagnostics) {
	inputParametersRaw := []any{}

	diags := diag.Diagnostics{}
	workflowRaw := OrchestratorTaskWorkflowAPIModel{}

	// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/object
	if self.Workflow.IsNull() || self.Workflow.IsUnknown() {
		diags.AddError(
			"Configuration error",
			fmt.Sprintf("Unable to manage %s, workflow is either null or unknown", self.String()))
	} else {
		// Convert workflow from object to raw
		workflow := OrchestratorTaskWorkflowModel{}
		diags.Append(self.Workflow.As(ctx, &workflow, basetypes.ObjectAsOptions{})...)
		workflowRaw = workflow.ToAPI()
	}

	return OrchestratorTaskAPIModel{
		Id:                  self.Id.ValueString(),
		Name:                self.Name.ValueString(),
		Description:         self.Description.ValueString(),
		Href:                self.Href.ValueString(),
		RecurrenceCycle:     self.RecurrenceCycle.ValueString(),
		RecurrencePattern:   self.RecurrencePattern.ValueString(),
		RecurrenceStartDate: self.RecurrenceStartDate.ValueString(),
		RecurrenceEndDate:   self.RecurrenceEndDate.ValueString(),
		RunningInstanceId:   self.RunningInstanceId.ValueString(),
		StartMode:           self.StartMode.ValueString(),
		State:               self.State.ValueString(),
		User:                self.User.ValueString(),
		InputParameters:     inputParametersRaw,
		Workflow:            workflowRaw,
	}, diags
}

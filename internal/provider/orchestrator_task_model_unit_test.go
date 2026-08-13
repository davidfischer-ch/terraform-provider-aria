// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// taskWorkflowObject builds a valid workflow, for tests not exercising workflow edge cases.
func taskWorkflowObject(t *testing.T) types.Object {
	t.Helper()
	ctx := t.Context()
	workflow := OrchestratorTaskWorkflowModel{
		Id:   types.StringValue("wf-1"),
		Name: types.StringValue("wf"),
	}
	obj, diags := types.ObjectValueFrom(ctx, workflow.AttributeTypes(), workflow)
	if diags.HasError() {
		t.Fatalf("workflow object: %v", diags.Errors())
	}
	return obj
}

// LockKey and the CRUD paths are derived straight from the id; pin their exact shape since the
// mutex-naming and API-path conventions matter for the client-level dedup logic.
func TestOrchestratorTaskModelLockKeyAndPaths(t *testing.T) {
	model := OrchestratorTaskModel{
		Id:   types.StringValue("task-123"),
		Name: types.StringValue("test-task"),
		User: types.StringValue("tester"),
	}

	if got, want := model.String(), "Orchestrator Task task-123 (test-task) of tester"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
	if got, want := model.LockKey(), "orchestrator-task-task-123"; got != want {
		t.Errorf("LockKey() = %q, want %q", got, want)
	}
	if got, want := model.CreatePath(), "vco/api/tasks"; got != want {
		t.Errorf("CreatePath() = %q, want %q", got, want)
	}
	if got, want := model.ReadPath(), "vco/api/tasks/task-123"; got != want {
		t.Errorf("ReadPath() = %q, want %q", got, want)
	}
	if got, want := model.UpdatePath(), "vco/api/tasks/task-123"; got != want {
		t.Errorf("UpdatePath() = %q, want %q", got, want)
	}
	if got, want := model.DeletePath(), "vco/api/tasks/task-123"; got != want {
		t.Errorf("DeletePath() = %q, want %q", got, want)
	}
}

// A non-empty recurrence-end-date must be parsed rather than treated as absent: only the
// zero-length string (an unbounded recurrence) short-circuits to null.
func TestOrchestratorTaskModelFromAPIRecurrenceEndDateSet(t *testing.T) {
	ctx := t.Context()

	model := OrchestratorTaskModel{}
	diags := model.FromAPI(ctx, OrchestratorTaskAPIModel{
		RecurrenceStartDate: "2050-01-06T05:02:00Z",
		RecurrenceEndDate:   "2055-01-06T05:02:00Z",
		Workflow:            OrchestratorTaskWorkflowAPIModel{Id: "wf-1", Name: "wf"},
	})
	if diags.HasError() {
		t.Fatalf("FromAPI: %v", diags.Errors())
	}
	if model.RecurrenceEndDate.IsNull() {
		t.Fatal("recurrence_end_date is null, want the parsed value")
	}
	if got := model.RecurrenceEndDate.ValueString(); got != "2055-01-06T05:02:00Z" {
		t.Errorf("recurrence_end_date = %q, want %q", got, "2055-01-06T05:02:00Z")
	}
}

// A malformed recurrence-end-date must surface as a diagnostic, not panic or silently truncate.
func TestOrchestratorTaskModelFromAPIRecurrenceEndDateInvalid(t *testing.T) {
	ctx := t.Context()

	model := OrchestratorTaskModel{}
	diags := model.FromAPI(ctx, OrchestratorTaskAPIModel{
		RecurrenceStartDate: "2050-01-06T05:02:00Z",
		RecurrenceEndDate:   "not-a-date",
		Workflow:            OrchestratorTaskWorkflowAPIModel{Id: "wf-1", Name: "wf"},
	})
	if !diags.HasError() {
		t.Fatal("expected FromAPI to report an error for a malformed recurrence_end_date")
	}
}

// A nil InputParameters (task with none configured) must decode to an empty, non-null list to
// match the schema's [] default, not null.
func TestOrchestratorTaskModelFromAPIInputParametersNil(t *testing.T) {
	ctx := t.Context()

	model := OrchestratorTaskModel{}
	diags := model.FromAPI(ctx, OrchestratorTaskAPIModel{
		RecurrenceStartDate: "2050-01-06T05:02:00Z",
		Workflow:            OrchestratorTaskWorkflowAPIModel{Id: "wf-1", Name: "wf"},
		InputParameters:     nil,
	})
	if diags.HasError() {
		t.Fatalf("FromAPI: %v", diags.Errors())
	}
	if model.InputParameters.IsNull() {
		t.Fatal("input_parameters is null, want an empty list to match the schema's [] default")
	}
	if got := len(model.InputParameters.Elements()); got != 0 {
		t.Errorf("input_parameters has %d elements, want 0", got)
	}
}

// Regression: vRO's Task API echoes description back empty. FromAPI must keep the prior/planned
// description instead of the API's blank one.
func TestOrchestratorTaskModelFromAPIPreservesDescriptionDespiteAPI(t *testing.T) {
	ctx := t.Context()

	priorList, diags := ParameterModelListFromAPI(ctx, []ParameterAPIModel{
		{Name: "vraHost", Type: "VRA:Host", Description: "Target vRA host"},
	})
	if diags.HasError() {
		t.Fatalf("prior input parameters list: %v", diags.Errors())
	}

	model := OrchestratorTaskModel{InputParameters: priorList}
	fromAPIDiags := model.FromAPI(ctx, OrchestratorTaskAPIModel{
		RecurrenceStartDate: "2050-01-06T05:02:00Z",
		Workflow:            OrchestratorTaskWorkflowAPIModel{Id: "wf-1", Name: "wf"},
		// The real API's behavior being regression-tested: description comes back empty even
		// though "Target vRA host" was sent on create.
		InputParameters: []ParameterAPIModel{
			{Name: "vraHost", Type: "VRA:Host", Description: ""},
		},
	})
	if fromAPIDiags.HasError() {
		t.Fatalf("FromAPI: %v", fromAPIDiags.Errors())
	}

	var out []ParameterModel
	if diags := model.InputParameters.ElementsAs(ctx, &out, false); diags.HasError() {
		t.Fatalf("input parameters elements: %v", diags.Errors())
	}
	if len(out) != 1 {
		t.Fatalf("input_parameters has %d entries, want 1", len(out))
	}
	if got := out[0].Description.ValueString(); got != "Target vRA host" {
		t.Errorf("description = %q, want %q (preserved despite the API's blank response)",
			got, "Target vRA host")
	}
	if got := out[0].Type.ValueString(); got != "VRA:Host" {
		t.Errorf("type = %q, want %q (still sourced from the API)", got, "VRA:Host")
	}
}

// A parameter added out-of-band (unknown to priorList) must fall back to the API's description
// instead of panicking.
func TestOrchestratorTaskModelFromAPIInputParametersUnknownToPrior(t *testing.T) {
	ctx := t.Context()

	model := OrchestratorTaskModel{InputParameters: types.ListValueMust(
		types.ObjectType{AttrTypes: ParameterModel{}.AttributeTypes()}, []attr.Value{})}

	diags := model.FromAPI(ctx, OrchestratorTaskAPIModel{
		RecurrenceStartDate: "2050-01-06T05:02:00Z",
		Workflow:            OrchestratorTaskWorkflowAPIModel{Id: "wf-1", Name: "wf"},
		InputParameters: []ParameterAPIModel{
			{Name: "driftedParam", Type: "string", Description: ""},
		},
	})
	if diags.HasError() {
		t.Fatalf("FromAPI: %v", diags.Errors())
	}

	var out []ParameterModel
	if diags := model.InputParameters.ElementsAs(ctx, &out, false); diags.HasError() {
		t.Fatalf("input parameters elements: %v", diags.Errors())
	}
	if len(out) != 1 || out[0].Name.ValueString() != "driftedParam" {
		t.Fatalf("input_parameters = %#v, want a single driftedParam entry", out)
	}
}

// ToAPI must reject a null or unknown workflow instead of silently sending a zero-value one.
func TestOrchestratorTaskModelToAPIWorkflowNull(t *testing.T) {
	ctx := t.Context()

	model := OrchestratorTaskModel{
		InputParameters: types.ListValueMust(
			types.ObjectType{AttrTypes: ParameterModel{}.AttributeTypes()}, []attr.Value{}),
		Workflow: types.ObjectNull(OrchestratorTaskWorkflowModel{}.AttributeTypes()),
	}

	_, diags := model.ToAPI(ctx)
	if !diags.HasError() {
		t.Fatal("expected ToAPI to report an error for a null workflow")
	}
}

func TestOrchestratorTaskModelToAPIWorkflowUnknown(t *testing.T) {
	ctx := t.Context()

	model := OrchestratorTaskModel{
		InputParameters: types.ListValueMust(
			types.ObjectType{AttrTypes: ParameterModel{}.AttributeTypes()}, []attr.Value{}),
		Workflow: types.ObjectUnknown(OrchestratorTaskWorkflowModel{}.AttributeTypes()),
	}

	_, diags := model.ToAPI(ctx)
	if !diags.HasError() {
		t.Fatal("expected ToAPI to report an error for an unknown workflow")
	}
}

// ToAPI must reject a null or unknown input_parameters list instead of silently dropping it.
func TestOrchestratorTaskModelToAPIInputParametersNull(t *testing.T) {
	ctx := t.Context()

	model := OrchestratorTaskModel{
		InputParameters: types.ListNull(
			types.ObjectType{AttrTypes: ParameterModel{}.AttributeTypes()}),
		Workflow: taskWorkflowObject(t),
	}

	_, diags := model.ToAPI(ctx)
	if !diags.HasError() {
		t.Fatal("expected ToAPI to report an error for a null input_parameters list")
	}
}

func TestOrchestratorTaskModelToAPIInputParametersUnknown(t *testing.T) {
	ctx := t.Context()

	model := OrchestratorTaskModel{
		InputParameters: types.ListUnknown(
			types.ObjectType{AttrTypes: ParameterModel{}.AttributeTypes()}),
		Workflow: taskWorkflowObject(t),
	}

	_, diags := model.ToAPI(ctx)
	if !diags.HasError() {
		t.Fatal("expected ToAPI to report an error for an unknown input_parameters list")
	}
}

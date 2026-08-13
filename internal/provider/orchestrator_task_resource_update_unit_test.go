// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// taskUpdateResult captures what the fake API observed and what Update produced.
type taskUpdateResult struct {
	updateCalls int
	stateSent   bool   // whether the update body carried a "state" field
	stateValue  string // the "state" field value, when present
	finalState  string
	hasError    bool
	errDetail   string
}

// runTaskUpdate drives the real resource Update against a fake Aria API. priorState is the state
// Terraform already tracks (post-refresh), plannedState is the desired state. It reports whether
// the update request carried a "state" field, mirroring the API's transition semantics.
func runTaskUpdate(t *testing.T, priorState, plannedState string) taskUpdateResult {
	t.Helper()
	ctx := t.Context()

	var res taskUpdateResult

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var parsed map[string]any
		_ = json.Unmarshal(body, &parsed)

		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/vco/api/tasks/task-123":
			res.updateCalls++
			if v, ok := parsed["state"]; ok {
				res.stateSent = true
				res.stateValue, _ = v.(string)
			}
			writeJSON(w, taskAPIBody("task-123", plannedState))
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	client := &AriaClient{
		Host:               srv.URL,
		AccessToken:        "fake-token",
		OKAPICallsLogLevel: "DEBUG",
		KOAPICallsLogLevel: "WARN",
		Context:            ctx,
	}
	if diags := client.Init(); diags.HasError() {
		t.Fatalf("AriaClient.Init: %v", diags.Errors())
	}

	res0, ok := NewOrchestratorTaskResource().(*OrchestratorTaskResource)
	if !ok {
		t.Fatal("NewOrchestratorTaskResource() did not return *OrchestratorTaskResource")
	}
	res0.client = client

	schema := OrchestratorTaskSchema()

	workflow := OrchestratorTaskWorkflowModel{
		Id:   types.StringValue("wf-1"),
		Name: types.StringValue("wf"),
	}
	workflowObject, diags := types.ObjectValueFrom(ctx, workflow.AttributeTypes(), workflow)
	if diags.HasError() {
		t.Fatalf("workflow object: %v", diags.Errors())
	}

	startDate, diags := timetypes.NewRFC3339Value("2050-01-06T05:02:00Z")
	if diags.HasError() {
		t.Fatalf("start date: %v", diags.Errors())
	}

	model := func(state string) OrchestratorTaskModel {
		return OrchestratorTaskModel{
			Id:                  types.StringValue("task-123"),
			Name:                types.StringValue("test-task"),
			Description:         types.StringValue("desc"),
			Href:                types.StringValue("https://aria.example/vco/api/tasks/task-123"),
			RecurrenceCycle:     types.StringValue("every-months"),
			RecurrencePattern:   types.StringValue("(Europe/Zurich) 01 00:00:00,"),
			RecurrenceStartDate: startDate,
			RecurrenceEndDate:   timetypes.NewRFC3339Null(),
			RunningInstanceId:   types.StringNull(),
			StartMode:           types.StringValue("normal"),
			State:               types.StringValue(state),
			User:                types.StringValue("tester"),
			Workflow:            workflowObject,
		}
	}

	prior := tfsdk.State{Schema: schema}
	if diags := prior.Set(ctx, model(priorState)); diags.HasError() {
		t.Fatalf("prior.Set: %v", diags.Errors())
	}

	plan := tfsdk.Plan{Schema: schema}
	if diags := plan.Set(ctx, model(plannedState)); diags.HasError() {
		t.Fatalf("plan.Set: %v", diags.Errors())
	}

	req := resource.UpdateRequest{Plan: plan, State: prior}
	resp := &resource.UpdateResponse{State: tfsdk.State{Schema: schema}}
	res0.Update(ctx, req, resp)

	res.hasError = resp.Diagnostics.HasError()
	for _, e := range resp.Diagnostics.Errors() {
		res.errDetail += e.Summary() + ": " + e.Detail() + "; "
	}

	if !resp.State.Raw.IsNull() {
		var out OrchestratorTaskModel
		if diags := resp.State.Get(ctx, &out); diags.HasError() {
			t.Fatalf("state.Get: %v", diags.Errors())
		}
		res.finalState = out.State.ValueString()
	}
	return res
}

// When the state does not change, Update must not send it: the API treats a state in the body as a
// transition request and rejects the no-op with 409 "Cannot resume task".
func TestOrchestratorTaskResourceUpdateKeepsState(t *testing.T) {
	res := runTaskUpdate(t, "suspended", "suspended")

	if res.hasError {
		t.Fatalf("unexpected Update error: %s", res.errDetail)
	}
	if res.updateCalls != 1 {
		t.Errorf("update requests = %d, want 1", res.updateCalls)
	}
	if res.stateSent {
		t.Errorf("update body carried state = %q, want it omitted for an unchanged state",
			res.stateValue)
	}
}

// When the desired state differs from the tracked state (including a drift Read surfaced), Update
// must send it to reconcile the task's lifecycle.
func TestOrchestratorTaskResourceUpdateChangesState(t *testing.T) {
	res := runTaskUpdate(t, "pending", "suspended")

	if res.hasError {
		t.Fatalf("unexpected Update error: %s", res.errDetail)
	}
	if res.updateCalls != 1 {
		t.Errorf("update requests = %d, want 1", res.updateCalls)
	}
	if !res.stateSent {
		t.Fatal("update body omitted state, want it sent to reconcile the changed state")
	}
	if res.stateValue != "suspended" {
		t.Errorf("update body state = %q, want %q", res.stateValue, "suspended")
	}
}

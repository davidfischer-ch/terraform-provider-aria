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

// taskCreateResult captures what the fake API observed and what Create produced.
type taskCreateResult struct {
	createState string
	updateState string
	updateCalls int
	finalState  string
	finalID     string
	hasError    bool
	errDetail   string
}

// taskAPIBody builds an OrchestratorTaskAPIModel-compatible response body. A valid
// recurrence-start-date is required, otherwise FromAPI fails to parse it.
func taskAPIBody(id, state string) map[string]any {
	return map[string]any{
		"id":                    id,
		"name":                  "test-task",
		"description":           "desc",
		"href":                  "https://aria.example/vco/api/tasks/" + id,
		"recurrence-cycle":      "every-months",
		"recurrence-pattern":    "(Europe/Zurich) 01 00:00:00,",
		"recurrence-start-date": "2050-01-06T05:02:00Z",
		"start-mode":            "normal",
		"state":                 state,
		"user":                  "tester",
		"workflow":              map[string]any{"id": "wf-1", "name": "wf"},
	}
}

// runTaskCreate drives the real resource Create against a fake Aria API and reports what the API
// received (create/update request bodies) and the state Create persisted. When failUpdate is set,
// the suspend update responds with an error.
func runTaskCreate(t *testing.T, plannedState string, failUpdate bool) taskCreateResult {
	t.Helper()
	ctx := t.Context()

	var res taskCreateResult

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var parsed map[string]any
		_ = json.Unmarshal(body, &parsed)
		state, _ := parsed["state"].(string)

		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/vco/api/tasks":
			res.createState = state
			// Content-Type must be set before WriteHeader, otherwise it is ignored and resty skips
			// unmarshalling the response body.
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)
			_ = json.NewEncoder(w).Encode(taskAPIBody("task-123", state))
		case r.Method == http.MethodPost && r.URL.Path == "/vco/api/tasks/task-123":
			res.updateCalls++
			res.updateState = state
			if failUpdate {
				w.WriteHeader(http.StatusInternalServerError)
				writeJSON(w, map[string]any{"message": "boom"})
				return
			}
			writeJSON(w, taskAPIBody("task-123", "suspended"))
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

	model := OrchestratorTaskModel{
		Id:                  types.StringNull(),
		Name:                types.StringValue("test-task"),
		Description:         types.StringValue("desc"),
		Href:                types.StringNull(),
		RecurrenceCycle:     types.StringValue("every-months"),
		RecurrencePattern:   types.StringValue("(Europe/Zurich) 01 00:00:00,"),
		RecurrenceStartDate: startDate,
		RecurrenceEndDate:   timetypes.NewRFC3339Null(),
		RunningInstanceId:   types.StringNull(),
		StartMode:           types.StringValue("normal"),
		State:               types.StringValue(plannedState),
		User:                types.StringNull(),
		Workflow:            workflowObject,
	}

	plan := tfsdk.Plan{Schema: schema}
	if diags := plan.Set(ctx, &model); diags.HasError() {
		t.Fatalf("plan.Set: %v", diags.Errors())
	}

	req := resource.CreateRequest{Plan: plan}
	resp := &resource.CreateResponse{State: tfsdk.State{Schema: schema}}
	res0.Create(ctx, req, resp)

	res.hasError = resp.Diagnostics.HasError()
	for _, e := range resp.Diagnostics.Errors() {
		res.errDetail += e.Summary() + ": " + e.Detail() + "; "
	}

	// State may be unset when the create itself fails; otherwise read what Create persisted.
	if !resp.State.Raw.IsNull() {
		var out OrchestratorTaskModel
		if diags := resp.State.Get(ctx, &out); diags.HasError() {
			t.Fatalf("state.Get: %v", diags.Errors())
		}
		res.finalState = out.State.ValueString()
		res.finalID = out.Id.ValueString()
	}
	return res
}

// A task requested as "suspended" must be created as "pending" then suspended with an update,
// otherwise the API falls back to "pending" and the apply is inconsistent.
func TestOrchestratorTaskResourceCreateSuspended(t *testing.T) {
	res := runTaskCreate(t, "suspended", false)

	if res.hasError {
		t.Fatalf("unexpected Create error: %s", res.errDetail)
	}
	if res.createState != "pending" {
		t.Errorf("create request state = %q, want %q", res.createState, "pending")
	}
	if res.updateCalls != 1 {
		t.Errorf("update requests = %d, want 1", res.updateCalls)
	}
	if res.updateState != "suspended" {
		t.Errorf("update request state = %q, want %q", res.updateState, "suspended")
	}
	if res.finalState != "suspended" {
		t.Errorf("persisted state = %q, want %q", res.finalState, "suspended")
	}
	if res.finalID != "task-123" {
		t.Errorf("persisted id = %q, want %q", res.finalID, "task-123")
	}
}

// A task requested in any other state is created in a single request, with no suspend update.
func TestOrchestratorTaskResourceCreatePending(t *testing.T) {
	res := runTaskCreate(t, "pending", false)

	if res.hasError {
		t.Fatalf("unexpected Create error: %s", res.errDetail)
	}
	if res.createState != "pending" {
		t.Errorf("create request state = %q, want %q", res.createState, "pending")
	}
	if res.updateCalls != 0 {
		t.Errorf("update requests = %d, want 0 (no suspend step)", res.updateCalls)
	}
	if res.finalState != "pending" {
		t.Errorf("persisted state = %q, want %q", res.finalState, "pending")
	}
}

// When the suspend update fails, the already-created task must still be persisted (tainted) so
// Terraform tracks it and can reconcile it later, rather than leaking it.
func TestOrchestratorTaskResourceCreateSuspendedUpdateFails(t *testing.T) {
	res := runTaskCreate(t, "suspended", true)

	if !res.hasError {
		t.Fatal("expected Create to report an error when the suspend update fails")
	}
	if res.updateCalls != 1 {
		t.Errorf("update requests = %d, want 1", res.updateCalls)
	}
	if res.finalID != "task-123" {
		t.Errorf("persisted id = %q, want %q (created task must stay tracked)",
			res.finalID, "task-123")
	}
	if res.finalState != "pending" {
		t.Errorf("persisted state = %q, want %q (intermediate create state)",
			res.finalState, "pending")
	}
}

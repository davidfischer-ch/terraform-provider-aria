// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// taskModel returns a minimal OrchestratorTaskModel usable for its CRUD paths. An orchestrator task
// is a convenient concrete Model: CreatePath is "vco/api/tasks" and ReadPath is that plus the id.
func taskModel(id string) *OrchestratorTaskModel {
	return &OrchestratorTaskModel{Id: types.StringValue(id)}
}

func TestAriaClientCreateIt(t *testing.T) {
	server := newFakeAPI(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/vco/api/tasks" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Location", "/vco/api/tasks/task-123")
		writeJSONStatus(w, http.StatusAccepted,
			map[string]any{"id": "task-123", "state": "pending"})
	})
	client := newTestClient(t, server.URL)

	var raw OrchestratorTaskAPIModel
	response, diags := client.CreateIt(taskModel(""), &raw, map[string]any{"name": "x"}, 202)
	if diags.HasError() {
		t.Fatalf("CreateIt: %v", diags.Errors())
	}
	if raw.Id != "task-123" {
		t.Errorf("decoded id = %q, want %q", raw.Id, "task-123")
	}
	if raw.State != "pending" {
		t.Errorf("decoded state = %q, want %q", raw.State, "pending")
	}
	id, err := client.GetIdFromLocation(response)
	if err != nil {
		t.Fatalf("GetIdFromLocation: %v", err)
	}
	if id != "task-123" {
		t.Errorf("id from Location = %q, want %q", id, "task-123")
	}
}

func TestAriaClientCreateItUnexpectedStatus(t *testing.T) {
	server := newFakeAPI(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSONStatus(w, http.StatusInternalServerError, map[string]any{"message": "boom"})
	})
	client := newTestClient(t, server.URL)

	var raw OrchestratorTaskAPIModel
	_, diags := client.CreateIt(taskModel(""), &raw, map[string]any{"name": "x"}, 202)
	if !diags.HasError() {
		t.Fatal("expected an error diagnostic when the API returns an unexpected status")
	}
}

func TestAriaClientReadItFound(t *testing.T) {
	server := newFakeAPI(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/vco/api/tasks/task-123" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		writeJSONStatus(w, http.StatusOK, map[string]any{"id": "task-123", "state": "suspended"})
	})
	client := newTestClient(t, server.URL)

	var raw OrchestratorTaskAPIModel
	found, _, diags := client.ReadIt(taskModel("task-123"), &raw)
	if diags.HasError() {
		t.Fatalf("ReadIt: %v", diags.Errors())
	}
	if !found {
		t.Fatal("found = false, want true")
	}
	if raw.State != "suspended" {
		t.Errorf("decoded state = %q, want %q", raw.State, "suspended")
	}
}

func TestAriaClientReadItNotFound(t *testing.T) {
	server := newFakeAPI(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSONStatus(w, http.StatusNotFound, map[string]any{"message": "gone"})
	})
	client := newTestClient(t, server.URL)

	var raw OrchestratorTaskAPIModel
	found, _, diags := client.ReadIt(taskModel("missing"), &raw)
	if diags.HasError() {
		t.Fatalf("a 404 must not produce an error diagnostic: %v", diags.Errors())
	}
	if found {
		t.Error("found = true, want false for a 404")
	}
}

func TestAriaClientUpdateIt(t *testing.T) {
	server := newFakeAPI(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/vco/api/tasks/task-123" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		writeJSONStatus(w, http.StatusOK, map[string]any{"id": "task-123", "state": "suspended"})
	})
	client := newTestClient(t, server.URL)

	var raw OrchestratorTaskAPIModel
	_, diags := client.UpdateIt(
		taskModel("task-123"), &raw, map[string]any{"state": "suspended"}, "POST", 200)
	if diags.HasError() {
		t.Fatalf("UpdateIt: %v", diags.Errors())
	}
	if raw.State != "suspended" {
		t.Errorf("decoded state = %q, want %q", raw.State, "suspended")
	}
}

func TestAriaClientDeleteIt(t *testing.T) {
	server := newFakeAPI(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		case http.MethodGet:
			// Deletion poll: report the resource is gone so DeleteIt returns immediately.
			writeJSONStatus(w, http.StatusNotFound, map[string]any{"message": "gone"})
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	})
	client := newTestClient(t, server.URL)

	diags := client.DeleteIt(taskModel("task-123"))
	if diags.HasError() {
		t.Fatalf("DeleteIt: %v", diags.Errors())
	}
}

func TestGetVersionFromPath(t *testing.T) {
	client := AriaClient{}
	cases := map[string]string{
		"abx/api/resources": ABX_API_VERSION,
		"blueprint/x":       BLUEPRINT_API_VERSION,
		"catalog/x":         CATALOG_API_VERSION,
		"form-service/x":    FORM_API_VERSION,
		"iaas/api/login":    IAAS_API_VERSION,
		"policy/x":          POLICY_API_VERSION,
		"project-service/x": PROJECT_API_VERSION,
		"properties/x":      BLUEPRINT_API_VERSION,
		"vco/api/tasks":     ORCHESTRATOR_API_VERSION,
	}
	for path, want := range cases {
		if got := client.GetVersionFromPath(path); got != want {
			t.Errorf("GetVersionFromPath(%q) = %q, want %q", path, got, want)
		}
	}
}

func TestGetVersionFromPathPanicsOnUnknownPrefix(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected a panic for an unmapped path prefix")
		}
	}()
	client := AriaClient{}
	_ = client.GetVersionFromPath("unknown/whatever")
}

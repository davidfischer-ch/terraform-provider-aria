// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
)

// ---- test helpers -------------------------------------------------------------------------------

// captureLogger records all Logf messages for later inspection.
type captureLogger struct {
	mu      sync.Mutex
	entries []string
}

func (l *captureLogger) Logf(format string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = append(l.entries, fmt.Sprintf(format, args...))
}

func (l *captureLogger) hasEntry(sub string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, e := range l.entries {
		if strings.Contains(e, sub) {
			return true
		}
	}
	return false
}

// requestRecord stores one incoming HTTP request.
type requestRecord struct {
	method string
	path   string
	query  url.Values
}

// requestLog records all requests received by a test server.
type requestLog struct {
	mu      sync.Mutex
	records []requestRecord
}

func (rl *requestLog) record(r *http.Request) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.records = append(rl.records, requestRecord{
		method: r.Method,
		path:   r.URL.Path,
		query:  r.URL.Query(),
	})
}

func (rl *requestLog) hasDelete(path string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	for _, rec := range rl.records {
		if rec.method == "DELETE" && rec.path == path {
			return true
		}
	}
	return false
}

func (rl *requestLog) hasDeleteWithParam(path, key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	for _, rec := range rl.records {
		if rec.method == "DELETE" && rec.path == path && rec.query.Get(key) == "true" {
			return true
		}
	}
	return false
}

func (rl *requestLog) hasAnyDelete() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	for _, rec := range rl.records {
		if rec.method == "DELETE" {
			return true
		}
	}
	return false
}

func (rl *requestLog) hasGetWithParam(path, key, value string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	for _, rec := range rl.records {
		if rec.method == "GET" && rec.path == path && rec.query.Get(key) == value {
			return true
		}
	}
	return false
}

// newTestRunner creates a CleanupRunner backed by a test HTTP server.
func newTestRunner(t *testing.T, srv *httptest.Server, dryRun bool) (*CleanupRunner, *captureLogger) {
	t.Helper()
	client := &AriaClient{
		Host:               srv.URL,
		AccessToken:        "fake-token",
		OKAPICallsLogLevel: "DEBUG",
		KOAPICallsLogLevel: "WARN",
		Context:            t.Context(),
	}
	if diags := client.Init(); diags.HasError() {
		t.Fatalf("AriaClient.Init: %v", diags.Errors())
	}
	logger := &captureLogger{}
	return &CleanupRunner{Client: client, Log: logger, DryRun: dryRun}, logger
}

// writeJSON sends a JSON response.
func writeJSON(w http.ResponseWriter, body any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(body)
}

// vROListBody builds a vROLinksAPIModel-compatible response body.
func vROListBody(links ...vROLinkAPIModel) map[string]any {
	return map[string]any{"link": links, "total": len(links)}
}

// vROTestLink builds a single vROLinkAPIModel with arbitrary attributes.
func vROTestLink(href string, attrs ...vROLinkAttributeAPIModel) vROLinkAPIModel {
	return vROLinkAPIModel{Href: href, Rel: "child", Attributes: attrs}
}

// vROAttr shorthand for vROLinkAttributeAPIModel.
func vROAttr(name, value string) vROLinkAttributeAPIModel {
	return vROLinkAttributeAPIModel{Name: name, Value: value}
}

// contentListBody builds a contentListAPIModel-compatible response body.
func contentListBody(items ...map[string]any) map[string]any {
	return map[string]any{"content": items}
}

// ---- idFromHref ----------------------------------------------------------------------------------

func TestIdFromHref(t *testing.T) {
	cases := []struct {
		href     string
		expected string
	}{
		{"https://host/vco/api/workflows/abc-123", "abc-123"},
		{"https://host/vco/api/categories/def-456/", "def-456"},
		{"vco/api/actions/ghi-789", "ghi-789"},
		{"https://host:8281/vco/api/workflows/uuid-with-dashes-0000", "uuid-with-dashes-0000"},
	}
	for _, tc := range cases {
		result := idFromHref(tc.href)
		CheckEqual(t, result, tc.expected)
	}
}

// ---- applyCleanups -------------------------------------------------------------------------------

func TestApplyCleanupsDoesNotDeleteOnDryRun(t *testing.T) {
	var rl requestLog
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.record(r)
		if r.Method == "DELETE" {
			t.Errorf("unexpected DELETE %s during dry-run", r.URL.Path)
		}
		w.WriteHeader(204)
	}))
	defer srv.Close()

	runner, logger := newTestRunner(t, srv, true /* dryRun */)
	runner.applyCleanups([]cleanupEntry{
		{label: "Workflow \"ARIA_PROVIDER_TEST_W\"", deletePath: "vco/api/workflows/id-1"},
		{label: "Tag \"ARIA_PROVIDER_TEST_T\"", deletePath: "iaas/api/tags/id-2"},
	})

	if rl.hasAnyDelete() {
		t.Error("dry-run: DELETE should never be called")
	}
	if !logger.hasEntry("Would delete") {
		t.Error("dry-run: expected 'Would delete' log entries")
	}
}

func TestApplyCleanupsCallsDeleteForEachEntry(t *testing.T) {
	var rl requestLog
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.record(r)
		w.WriteHeader(204)
	}))
	defer srv.Close()

	runner, _ := newTestRunner(t, srv, false)
	runner.applyCleanups([]cleanupEntry{
		{label: "A", deletePath: "vco/api/workflows/id-1"},
		{label: "B", deletePath: "vco/api/workflows/id-2"},
	})

	if !rl.hasDelete("/vco/api/workflows/id-1") {
		t.Error("expected DELETE /vco/api/workflows/id-1")
	}
	if !rl.hasDelete("/vco/api/workflows/id-2") {
		t.Error("expected DELETE /vco/api/workflows/id-2")
	}
}

func TestApplyCleanupsTolerates404(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer srv.Close()

	runner, logger := newTestRunner(t, srv, false)
	// Must not panic and must not log a warning
	runner.applyCleanups([]cleanupEntry{
		{label: "already gone", deletePath: "vco/api/workflows/gone"},
	})

	if logger.hasEntry("Warning") {
		t.Error("404 on DELETE should not produce a warning")
	}
}

func TestApplyCleanupsLogsWarningOnAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer srv.Close()

	runner, logger := newTestRunner(t, srv, false)
	runner.applyCleanups([]cleanupEntry{
		{label: "broken resource", deletePath: "vco/api/workflows/broken"},
	})

	if !logger.hasEntry("Warning") {
		t.Error("expected warning log on 500 response")
	}
}

// ---- vROCleanupsByPrefix -------------------------------------------------------------------------

func TestVROCleanupsByPrefixFiltersCorrectly(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, vROListBody(
			vROTestLink("https://host/vco/api/workflows/match-1",
				vROAttr("name", "ARIA_PROVIDER_TEST_WORKFLOW_A")),
			vROTestLink("https://host/vco/api/workflows/match-2",
				vROAttr("name", "ARIA_PROVIDER_TEST_WORKFLOW_B")),
			vROTestLink("https://host/vco/api/workflows/no-match",
				vROAttr("name", "ProductionWorkflow")),
		))
	}))
	defer srv.Close()

	runner, _ := newTestRunner(t, srv, false)
	entries := runner.vROCleanupsByPrefix("vco/api/workflows", "Workflow", "name")

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	CheckEqual(t, entries[0].deletePath, "vco/api/workflows/match-1")
	CheckEqual(t, entries[1].deletePath, "vco/api/workflows/match-2")
}

func TestVROCleanupsByPrefixHandlesEmptyList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, vROListBody())
	}))
	defer srv.Close()

	runner, _ := newTestRunner(t, srv, false)
	entries := runner.vROCleanupsByPrefix("vco/api/workflows", "Workflow", "name")

	if entries != nil {
		t.Errorf("expected nil entries for empty list, got %v", entries)
	}
}

func TestVROCleanupsByPrefixHandlesListError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer srv.Close()

	runner, logger := newTestRunner(t, srv, false)
	entries := runner.vROCleanupsByPrefix("vco/api/workflows", "Workflow", "name")

	if entries != nil {
		t.Error("expected nil entries on list error")
	}
	if !logger.hasEntry("Warning") {
		t.Error("expected warning log on list error")
	}
}

func TestVROCleanupsByPrefixUsesCorrectAttribute(t *testing.T) {
	// A link has both "name" and "fqn"; only "fqn" should be checked when nameField="fqn".
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, vROListBody(
			vROTestLink("https://host/vco/api/actions/act-1",
				vROAttr("name", "getVRAHost"),
				vROAttr("fqn", "ARIA_PROVIDER_TEST_ACTIONS/getVRAHost")),
			vROTestLink("https://host/vco/api/actions/act-2",
				vROAttr("name", "ARIA_PROVIDER_TEST_shouldNotMatch"),
				vROAttr("fqn", "production/someAction")),
		))
	}))
	defer srv.Close()

	runner, _ := newTestRunner(t, srv, false)
	entries := runner.vROCleanupsByPrefix("vco/api/actions", "Action", "fqn")

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry (fqn match only), got %d", len(entries))
	}
	CheckEqual(t, entries[0].deletePath, "vco/api/actions/act-1")
}

// ---- contentCleanupsByPrefix --------------------------------------------------------------------

func TestContentCleanupsByPrefixFiltersCorrectly(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, contentListBody(
			map[string]any{"id": "t-1", "key": "ARIA_PROVIDER_TEST_TAG"},
			map[string]any{"id": "t-2", "key": "ARIA_PROVIDER_TEST_OTHER"},
			map[string]any{"id": "t-3", "key": "production:tag"},
		))
	}))
	defer srv.Close()

	runner, _ := newTestRunner(t, srv, false)
	entries := runner.contentCleanupsByPrefix("iaas/api/tags", "key")

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	CheckEqual(t, entries[0].deletePath, "iaas/api/tags/t-1")
	CheckEqual(t, entries[1].deletePath, "iaas/api/tags/t-2")
}

func TestContentCleanupsByPrefixSkipsItemsWithoutID(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, contentListBody(
			map[string]any{"name": "ARIA_PROVIDER_TEST_NOID"}, // no "id" field
			map[string]any{"id": "p-1", "name": "ARIA_PROVIDER_TEST_HAS_ID"},
		))
	}))
	defer srv.Close()

	runner, _ := newTestRunner(t, srv, false)
	entries := runner.contentCleanupsByPrefix("policy/api/policies", "name")

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry (item with id only), got %d", len(entries))
	}
	CheckEqual(t, entries[0].deletePath, "policy/api/policies/p-1")
}

func TestContentCleanupsByPrefixHandlesEmptyContent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, contentListBody())
	}))
	defer srv.Close()

	runner, _ := newTestRunner(t, srv, false)
	entries := runner.contentCleanupsByPrefix("policy/api/policies", "name")

	if entries != nil {
		t.Errorf("expected nil entries for empty content, got %v", entries)
	}
}

func TestContentCleanupsByPrefixHandlesListError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(503)
	}))
	defer srv.Close()

	runner, logger := newTestRunner(t, srv, false)
	entries := runner.contentCleanupsByPrefix("policy/api/policies", "name")

	if entries != nil {
		t.Error("expected nil entries on list error")
	}
	if !logger.hasEntry("Warning") {
		t.Error("expected warning log on list error")
	}
}

// ---- contentCleanupsByPrefixInProject -----------------------------------------------------------

func TestContentCleanupsByPrefixInProjectSendsProjectID(t *testing.T) {
	var rl requestLog
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.record(r)
		writeJSON(w, contentListBody())
	}))
	defer srv.Close()

	runner, _ := newTestRunner(t, srv, false)
	runner.contentCleanupsByPrefixInProject(
		"abx/api/resources/actions", TestPrefix, "proj-xyz")

	if !rl.hasGetWithParam("/abx/api/resources/actions", "projectId", "proj-xyz") {
		t.Error("expected GET with projectId=proj-xyz query param")
	}
}

// ---- sweep method spot-checks -------------------------------------------------------------------

func TestOrchestratorWorkflowsAppliesForceDelete(t *testing.T) {
	var rl requestLog
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.record(r)
		switch r.Method {
		case "GET":
			writeJSON(w, vROListBody(
				vROTestLink("https://host/vco/api/workflows/wf-1",
					vROAttr("name", "ARIA_PROVIDER_TEST_WORKFLOW")),
			))
		case "DELETE":
			w.WriteHeader(204)
		}
	}))
	defer srv.Close()

	runner, _ := newTestRunner(t, srv, false)
	runner.Force = true
	runner.OrchestratorWorkflows()

	if !rl.hasDeleteWithParam("/vco/api/workflows/wf-1", "force") {
		t.Error("OrchestratorWorkflows: expected DELETE with force=true when Force=true")
	}
}

func TestOrchestratorWorkflowsNoForceWithoutFlag(t *testing.T) {
	var rl requestLog
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.record(r)
		switch r.Method {
		case "GET":
			writeJSON(w, vROListBody(
				vROTestLink("https://host/vco/api/workflows/wf-1",
					vROAttr("name", "ARIA_PROVIDER_TEST_WORKFLOW")),
			))
		case "DELETE":
			w.WriteHeader(204)
		}
	}))
	defer srv.Close()

	runner, _ := newTestRunner(t, srv, false)
	runner.OrchestratorWorkflows()

	if rl.hasDeleteWithParam("/vco/api/workflows/wf-1", "force") {
		t.Error("OrchestratorWorkflows: force=true must not be sent when Force=false")
	}
	if !rl.hasDelete("/vco/api/workflows/wf-1") {
		t.Error("OrchestratorWorkflows: DELETE should still be issued without Force")
	}
}

func TestOrchestratorWorkflowsDoesNotDeleteNonMatching(t *testing.T) {
	var rl requestLog
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.record(r)
		switch r.Method {
		case "GET":
			writeJSON(w, vROListBody(
				vROTestLink("https://host/vco/api/workflows/prod-1",
					vROAttr("name", "ProductionWorkflow")),
				vROTestLink("https://host/vco/api/workflows/prod-2",
					vROAttr("name", "ImportantWorkflow")),
			))
		case "DELETE":
			t.Errorf("unexpected DELETE %s: non-test resources must never be deleted", r.URL.Path)
			w.WriteHeader(204)
		}
	}))
	defer srv.Close()

	runner, _ := newTestRunner(t, srv, false)
	runner.OrchestratorWorkflows()

	if rl.hasAnyDelete() {
		t.Error("non-matching workflows must not be deleted")
	}
}

func TestOrchestratorActionsFiltersOnFQN(t *testing.T) {
	var rl requestLog
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.record(r)
		switch r.Method {
		case "GET":
			writeJSON(w, vROListBody(
				// fqn starts with ARIA_PROVIDER_TEST → should be deleted
				vROTestLink("https://host/vco/api/actions/act-match",
					vROAttr("name", "myAction"),
					vROAttr("fqn", "ARIA_PROVIDER_TEST_ACTIONS/myAction")),
				// fqn does NOT start with ARIA_PROVIDER_TEST → must be kept
				vROTestLink("https://host/vco/api/actions/act-keep",
					vROAttr("name", "ARIA_PROVIDER_TEST_shouldNotMatchOnName"),
					vROAttr("fqn", "com.vmware.production/importantAction")),
			))
		case "DELETE":
			w.WriteHeader(204)
		}
	}))
	defer srv.Close()

	runner, _ := newTestRunner(t, srv, false)
	runner.OrchestratorActions()

	if !rl.hasDelete("/vco/api/actions/act-match") {
		t.Error("expected DELETE for action with matching fqn")
	}
	if rl.hasDelete("/vco/api/actions/act-keep") {
		t.Error("action with non-matching fqn must NOT be deleted")
	}
}

func TestOrchestratorActionsAppliesForceDelete(t *testing.T) {
	var rl requestLog
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.record(r)
		switch r.Method {
		case "GET":
			writeJSON(w, vROListBody(
				vROTestLink("https://host/vco/api/actions/act-1",
					vROAttr("fqn", "ARIA_PROVIDER_TEST_ACTIONS/act")),
			))
		case "DELETE":
			w.WriteHeader(204)
		}
	}))
	defer srv.Close()

	runner, _ := newTestRunner(t, srv, false)
	runner.Force = true
	runner.OrchestratorActions()

	if !rl.hasDeleteWithParam("/vco/api/actions/act-1", "force") {
		t.Error("OrchestratorActions: expected DELETE with force=true when Force=true")
	}
}

func TestOrchestratorActionsNoForceWithoutFlag(t *testing.T) {
	var rl requestLog
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.record(r)
		switch r.Method {
		case "GET":
			writeJSON(w, vROListBody(
				vROTestLink("https://host/vco/api/actions/act-1",
					vROAttr("fqn", "ARIA_PROVIDER_TEST_ACTIONS/act")),
			))
		case "DELETE":
			w.WriteHeader(204)
		}
	}))
	defer srv.Close()

	runner, _ := newTestRunner(t, srv, false)
	runner.OrchestratorActions()

	if rl.hasDeleteWithParam("/vco/api/actions/act-1", "force") {
		t.Error("OrchestratorActions: force=true must not be sent when Force=false")
	}
	if !rl.hasDelete("/vco/api/actions/act-1") {
		t.Error("OrchestratorActions: DELETE should still be issued without Force")
	}
}

func TestOrchestratorCategoriesUsesCategoryPrefix(t *testing.T) {
	var rl requestLog
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.record(r)
		switch r.Method {
		case "GET":
			writeJSON(w, vROListBody(
				// Matches TestPrefix = "ARIA_PROVIDER_TEST"
				vROTestLink("https://host/vco/api/categories/cat-1",
					vROAttr("name", "ARIA_PROVIDER_TEST")),
				vROTestLink("https://host/vco/api/categories/cat-2",
					vROAttr("name", "ARIA_PROVIDER_TEST_ACTIONS")),
				// Does NOT match — unrelated production category
				vROTestLink("https://host/vco/api/categories/cat-keep",
					vROAttr("name", "com.vmware.orchestrator")),
			))
		case "DELETE":
			w.WriteHeader(204)
		}
	}))
	defer srv.Close()

	runner, _ := newTestRunner(t, srv, false)
	runner.OrchestratorCategories()

	if !rl.hasDelete("/vco/api/categories/cat-1") {
		t.Error("expected DELETE for category matching TestPrefix")
	}
	if !rl.hasDelete("/vco/api/categories/cat-2") {
		t.Error("expected DELETE for category matching TestPrefix")
	}
	if rl.hasDelete("/vco/api/categories/cat-keep") {
		t.Error("unrelated production category must NOT be deleted by OrchestratorCategories")
	}
}

func TestTagsUsesKeyField(t *testing.T) {
	var rl requestLog
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.record(r)
		switch r.Method {
		case "GET":
			writeJSON(w, contentListBody(
				map[string]any{"id": "tag-1", "key": "ARIA_PROVIDER_TEST_KEY", "value": "v"},
				// "name" field is not what Tags() filters on
				map[string]any{"id": "tag-2", "name": "ARIA_PROVIDER_TEST_NAME", "key": "prod:key"},
			))
		case "DELETE":
			w.WriteHeader(204)
		}
	}))
	defer srv.Close()

	runner, _ := newTestRunner(t, srv, false)
	runner.Tags()

	if !rl.hasDelete("/iaas/api/tags/tag-1") {
		t.Error("expected DELETE for tag with matching key")
	}
	if rl.hasDelete("/iaas/api/tags/tag-2") {
		t.Error("tag with non-matching key must NOT be deleted (name field is irrelevant)")
	}
}

func TestTagsAppliesIgnoreUsage(t *testing.T) {
	var rl requestLog
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.record(r)
		switch r.Method {
		case "GET":
			writeJSON(w, contentListBody(
				map[string]any{"id": "tag-1", "key": "ARIA_PROVIDER_TEST_TAG"},
			))
		case "DELETE":
			w.WriteHeader(204)
		}
	}))
	defer srv.Close()

	runner, _ := newTestRunner(t, srv, false)
	runner.Force = true
	runner.Tags()

	if !rl.hasDeleteWithParam("/iaas/api/tags/tag-1", "ignoreUsage") {
		t.Error("Tags: expected DELETE with ignoreUsage=true when Force=true")
	}
}

func TestTagsNoIgnoreUsageWithoutFlag(t *testing.T) {
	var rl requestLog
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.record(r)
		switch r.Method {
		case "GET":
			writeJSON(w, contentListBody(
				map[string]any{"id": "tag-1", "key": "ARIA_PROVIDER_TEST_TAG"},
			))
		case "DELETE":
			w.WriteHeader(204)
		}
	}))
	defer srv.Close()

	runner, _ := newTestRunner(t, srv, false)
	runner.Tags()

	if rl.hasDeleteWithParam("/iaas/api/tags/tag-1", "ignoreUsage") {
		t.Error("Tags: ignoreUsage=true must not be sent when Force=false")
	}
	if !rl.hasDelete("/iaas/api/tags/tag-1") {
		t.Error("Tags: DELETE should still be issued without Force")
	}
}

func TestCustomFormsDeletesMatchingForm(t *testing.T) {
	var rl requestLog
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.record(r)
		switch r.Method {
		case "GET":
			writeJSON(w, map[string]any{
				"id":   "form-1",
				"name": "ARIA_PROVIDER_TEST_FORM",
			})
		case "DELETE":
			w.WriteHeader(204)
		}
	}))
	defer srv.Close()

	runner, _ := newTestRunner(t, srv, false)
	runner.CustomForms("src-id", "com.vmware.type")

	if !rl.hasDelete("/form-service/api/forms/form-1") {
		t.Error("expected DELETE for custom form with matching name")
	}
}

func TestCustomFormsSkipsWhenNotFound(t *testing.T) {
	var rl requestLog
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.record(r)
		w.WriteHeader(404)
	}))
	defer srv.Close()

	runner, logger := newTestRunner(t, srv, false)
	runner.CustomForms("src-id", "com.vmware.type")

	if rl.hasAnyDelete() {
		t.Error("404 on fetchBySourceAndType: no DELETE should be made")
	}
	if logger.hasEntry("Warning") {
		t.Error("404 on fetchBySourceAndType should not produce a warning")
	}
}

func TestCustomFormsSkipsNonMatchingName(t *testing.T) {
	var rl requestLog
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.record(r)
		switch r.Method {
		case "GET":
			writeJSON(w, map[string]any{
				"id":   "form-prod",
				"name": "ProductionForm",
			})
		case "DELETE":
			t.Errorf("must not DELETE form whose name does not start with %s", TestPrefix)
			w.WriteHeader(204)
		}
	}))
	defer srv.Close()

	runner, _ := newTestRunner(t, srv, false)
	runner.CustomForms("src-id", "com.vmware.type")

	if rl.hasAnyDelete() {
		t.Error("custom form with non-matching name must NOT be deleted")
	}
}

func TestOrchestratorCategoriesSkipsSubCategoriesWithNonMatchingName(t *testing.T) {
	rl := &requestLog{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.record(r)
		switch r.Method + " " + r.URL.Path {
		case "GET /vco/api/categories":
			writeJSON(w, vROListBody(
				vROTestLink("https://host/vco/api/categories/id-root", vROAttr("name", "ARIA_PROVIDER_TEST")),
				vROTestLink("https://host/vco/api/categories/id-a", vROAttr("name", "A")),
				vROTestLink("https://host/vco/api/categories/id-b", vROAttr("name", "B")),
			))
		default:
			w.WriteHeader(204)
		}
	}))
	defer srv.Close()

	runner, _ := newTestRunner(t, srv, false)
	runner.OrchestratorCategories()

	if !rl.hasDelete("/vco/api/categories/id-root") {
		t.Error("expected root category ARIA_PROVIDER_TEST to be deleted")
	}
	if rl.hasDelete("/vco/api/categories/id-a") {
		t.Error("sub-category A must not be individually deleted")
	}
	if rl.hasDelete("/vco/api/categories/id-b") {
		t.Error("sub-category B must not be individually deleted")
	}
}

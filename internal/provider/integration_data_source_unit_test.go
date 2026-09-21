// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// The integrations the fake API lists. The endpoint mixes every kind, an organization holding
// several Orchestrators returns several entries of type vro.
func fakeIntegrationsResponse(entries ...IntegrationEntryAPIModel) map[string]any {
	return map[string]any{
		"content":       entries,
		"totalElements": len(entries),
	}
}

// Return one page of entries starting at the requested offset, of at most pageSize entries. An API
// serving fewer entries than asked for is what the walk has to cope with.
func fakeIntegrationsPage(
	skip int,
	pageSize int,
	entries []IntegrationEntryAPIModel,
) map[string]any {
	page := []IntegrationEntryAPIModel{}
	for index := skip; index < len(entries) && len(page) < pageSize; index++ {
		page = append(page, entries[index])
	}
	return map[string]any{
		"content":       page,
		"totalElements": len(entries),
	}
}

type integrationReadResult struct {
	hasError  bool
	errDetail string
	name      string
	endpoint  string
	link      string
}

// readIntegration runs the data source against a fake API returning the given integrations,
// filtering on name when it is not empty.
func readIntegration(
	t *testing.T,
	name string,
	entries ...IntegrationEntryAPIModel,
) integrationReadResult {
	t.Helper()
	return readIntegrationResponse(t, name, fakeIntegrationsResponse(entries...))
}

// readIntegrationResponse runs the data source against a fake API returning the given body.
func readIntegrationResponse(
	t *testing.T,
	name string,
	body map[string]any,
) integrationReadResult {
	t.Helper()

	server := newFakeAPI(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/iaas/api/integrations") {
			writeJSONStatus(w, http.StatusNotFound, map[string]any{})
			return
		}
		writeJSONStatus(w, http.StatusOK, body)
	})

	return readIntegrationAt(t, server.URL, name)
}

// readIntegrationAt runs the data source against the API served at host.
func readIntegrationAt(t *testing.T, host string, name string) integrationReadResult {
	t.Helper()

	ctx := t.Context()

	source, ok := NewIntegrationDataSource().(*IntegrationDataSource)
	if !ok {
		t.Fatal("NewIntegrationDataSource() did not return *IntegrationDataSource")
	}
	source.client = newTestClient(t, host)

	schema := IntegrationDataSourceSchema()

	model := IntegrationDataSourceModel{TypeId: types.StringValue("com.vmw.vro.workflow")}
	model.Name = types.StringNull()
	if len(name) > 0 {
		model.Name = types.StringValue(name)
	}
	model.EndpointConfigurationLink = types.StringNull()
	model.EndpointURI = types.StringNull()

	// Config holds no setter, the value is built through a State sharing the same schema.
	configState := tfsdk.State{Schema: schema}
	if diags := configState.Set(ctx, &model); diags.HasError() {
		t.Fatalf("config.Set: %v", diags.Errors())
	}
	config := tfsdk.Config{Schema: schema, Raw: configState.Raw}

	req := datasource.ReadRequest{Config: config}
	resp := &datasource.ReadResponse{State: tfsdk.State{Schema: schema}}
	source.Read(ctx, req, resp)

	res := integrationReadResult{hasError: resp.Diagnostics.HasError()}
	for _, e := range resp.Diagnostics.Errors() {
		res.errDetail += e.Summary() + ": " + e.Detail() + "; "
	}

	if !res.hasError {
		var out IntegrationDataSourceModel
		if diags := resp.State.Get(ctx, &out); diags.HasError() {
			t.Fatalf("state.Get: %v", diags.Errors())
		}
		res.name = out.Name.ValueString()
		res.endpoint = out.EndpointURI.ValueString()
		res.link = out.EndpointConfigurationLink.ValueString()
	}

	return res
}

var (
	embeddedIntegration = IntegrationEntryAPIModel{
		Id:               "8a430db3-924c-4d58-a29a-da811f9c992e",
		Name:             "embedded-VRO",
		IntegrationType:  "vro",
		CustomProperties: map[string]string{"hostName": "https://embedded.example.net"},
	}
	externalIntegration = IntegrationEntryAPIModel{
		Id:               "1ce2fa30-8f37-4a3e-8c0c-8b1f39aaf6ab",
		Name:             "External Orchestrator",
		IntegrationType:  "vro",
		CustomProperties: map[string]string{"hostName": "https://external.example.net"},
	}
	abxIntegration = IntegrationEntryAPIModel{
		Id:               "e2cb828c-cf8c-4fd3-82e5-85ff2a93deec",
		Name:             "embedded-ABX-onprem",
		IntegrationType:  "abx.endpoint",
		CustomProperties: map[string]string{"apiEndpoint": "http://gateway.example.net:8080"},
	}
)

func TestIntegrationDataSourceReadSingle(t *testing.T) {
	// The integrations of another kind are not candidates.
	res := readIntegration(t, "", embeddedIntegration, abxIntegration)
	if res.hasError {
		t.Fatalf("unexpected error: %s", res.errDetail)
	}
	if res.name != "embedded-VRO" {
		t.Errorf("name = %q, want %q", res.name, "embedded-VRO")
	}
	if res.endpoint != "https://embedded.example.net" {
		t.Errorf("endpoint_uri = %q, want %q", res.endpoint, "https://embedded.example.net")
	}
	want := "/resources/endpoints/8a430db3-924c-4d58-a29a-da811f9c992e"
	if res.link != want {
		t.Errorf("endpoint_configuration_link = %q, want %q", res.link, want)
	}
}

func TestIntegrationDataSourceReadAmbiguous(t *testing.T) {
	res := readIntegration(t, "", embeddedIntegration, externalIntegration)
	if !res.hasError {
		t.Fatalf("read succeeded, want an error naming both candidates")
	}
	for _, name := range []string{"embedded-VRO", "External Orchestrator"} {
		if !strings.Contains(res.errDetail, name) {
			t.Errorf("error %q does not name candidate %q", res.errDetail, name)
		}
	}
}

func TestIntegrationDataSourceReadByName(t *testing.T) {
	res := readIntegration(t, "External Orchestrator", embeddedIntegration, externalIntegration)
	if res.hasError {
		t.Fatalf("unexpected error: %s", res.errDetail)
	}
	if res.name != "External Orchestrator" {
		t.Errorf("name = %q, want %q", res.name, "External Orchestrator")
	}
	if res.endpoint != "https://external.example.net" {
		t.Errorf("endpoint_uri = %q, want %q", res.endpoint, "https://external.example.net")
	}
}

func TestIntegrationDataSourceReadUnknownName(t *testing.T) {
	res := readIntegration(t, "Nowhere Orchestrator", embeddedIntegration)
	if !res.hasError {
		t.Fatalf("read succeeded, want an error for an unknown name")
	}
	if !strings.Contains(res.errDetail, "embedded-VRO") {
		t.Errorf("error %q does not name the available candidate", res.errDetail)
	}
}

func TestIntegrationDataSourceReadNoCandidate(t *testing.T) {
	// An organization holding no Orchestrator at all, the ABX integration is not a candidate.
	res := readIntegration(t, "", abxIntegration)
	if !res.hasError {
		t.Fatalf("read succeeded, want an error when no integration of the type exists")
	}
	if !strings.Contains(res.errDetail, "0 found") {
		t.Errorf("error %q does not report how many candidates were found", res.errDetail)
	}
}

// Names are not unique, filtering on one can leave several candidates.
func TestIntegrationDataSourceReadDuplicateName(t *testing.T) {
	other := externalIntegration
	other.Id = "6f1b0e64-1f38-4e5f-8f27-2f9b5c4a1d70"
	other.CustomProperties = map[string]string{"hostName": "https://other.example.net"}

	res := readIntegration(t, "External Orchestrator", externalIntegration, other)
	if !res.hasError {
		t.Fatal("read succeeded, want an error when two integrations share the name")
	}
	if !strings.Contains(res.errDetail, "2 integrations") {
		t.Errorf("error %q does not report how many match", res.errDetail)
	}
	for _, uri := range []string{"https://external.example.net", "https://other.example.net"} {
		if !strings.Contains(res.errDetail, uri) {
			t.Errorf("error %q does not tell the candidates apart by %s", res.errDetail, uri)
		}
	}
}

// readIntegrationPaged runs the data source against a fake API serving pageSize entries at a time,
// ignoring the requested offset when honourSkip is false.
func readIntegrationPaged(
	t *testing.T,
	name string,
	pageSize int,
	honourSkip bool,
	entries []IntegrationEntryAPIModel,
) integrationReadResult {
	t.Helper()

	server := newFakeAPI(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/iaas/api/integrations") {
			writeJSONStatus(w, http.StatusNotFound, map[string]any{})
			return
		}
		if order := r.URL.Query().Get("$orderby"); order != "name asc" {
			t.Errorf("request orders by %q, want the offset to walk a stable order", order)
		}
		skip := 0
		if honourSkip {
			skip, _ = strconv.Atoi(r.URL.Query().Get("$skip"))
		}
		writeJSONStatus(w, http.StatusOK, fakeIntegrationsPage(skip, pageSize, entries))
	})

	return readIntegrationAt(t, server.URL, name)
}

// Three Orchestrators served one per page, every one of them has to be seen.
func TestIntegrationDataSourceReadPaginated(t *testing.T) {
	third := externalIntegration
	third.Id = "6f1b0e64-1f38-4e5f-8f27-2f9b5c4a1d70"
	third.Name = "Third Orchestrator"
	third.CustomProperties = map[string]string{"hostName": "https://third.example.net"}

	entries := []IntegrationEntryAPIModel{embeddedIntegration, externalIntegration, third}

	res := readIntegrationPaged(t, "Third Orchestrator", 1, true, entries)
	if res.hasError {
		t.Fatalf("unexpected error: %s", res.errDetail)
	}
	if res.endpoint != "https://third.example.net" {
		t.Errorf("endpoint_uri = %q, want the entry of the last page", res.endpoint)
	}

	// Without a name the three candidates are ambiguous, which proves all of them were walked.
	res = readIntegrationPaged(t, "", 1, true, entries)
	if !res.hasError {
		t.Fatal("read succeeded, want an error naming the three candidates")
	}
	if !strings.Contains(res.errDetail, "3 integrations") {
		t.Errorf("error %q does not report the three candidates", res.errDetail)
	}
}

// An API ignoring the offset serves the same page forever, the walk must not.
func TestIntegrationDataSourceReadPaginationIgnored(t *testing.T) {
	entries := []IntegrationEntryAPIModel{embeddedIntegration, externalIntegration}

	res := readIntegrationPaged(t, "", 1, false, entries)
	if !res.hasError {
		t.Fatal("read succeeded, want an error for an API ignoring the offset")
	}
	if !strings.Contains(res.errDetail, "stopped returning new ones") {
		t.Errorf("error %q does not report the lack of progress", res.errDetail)
	}
}

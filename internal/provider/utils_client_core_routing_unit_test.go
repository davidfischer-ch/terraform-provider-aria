// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// newTestClientWithVRO returns a client addressing Aria at host and a standalone Orchestrator at
// vroHost.
func newTestClientWithVRO(t *testing.T, host string, vroHost string) *AriaClient {
	t.Helper()
	client := &AriaClient{
		Host:               host,
		VROHost:            vroHost,
		AccessToken:        "fake-token",
		OKAPICallsLogLevel: "DEBUG",
		KOAPICallsLogLevel: "WARN",
		Context:            t.Context(),
	}
	if diags := client.Init(); diags.HasError() {
		t.Fatalf("AriaClient.Init: %v", diags.Errors())
	}
	return client
}

func TestAriaClientRoutesOrchestratorCallsToVROHost(t *testing.T) {
	ariaPaths := []string{}
	aria := newFakeAPI(t, func(w http.ResponseWriter, r *http.Request) {
		ariaPaths = append(ariaPaths, r.URL.Path)
		writeJSONStatus(w, http.StatusOK, map[string]any{"id": "id-1"})
	})

	vroPaths := []string{}
	vro := newFakeAPI(t, func(w http.ResponseWriter, r *http.Request) {
		vroPaths = append(vroPaths, r.URL.Path)
		if got := r.Header.Get("Authorization"); got != "Bearer fake-token" {
			t.Errorf("Authorization = %q, want %q", got, "Bearer fake-token")
		}
		writeJSONStatus(w, http.StatusOK, map[string]any{"id": "task-123", "state": "pending"})
	})

	client := newTestClientWithVRO(t, aria.URL, vro.URL)

	// Orchestrator's own API is served by the standalone Orchestrator.
	var task OrchestratorTaskAPIModel
	if _, _, diags := client.ReadIt(taskModel("task-123"), &task); diags.HasError() {
		t.Fatalf("ReadIt(vco): %v", diags.Errors())
	}

	// The Orchestrator gateway is a service of Aria Automation and stays on the main host.
	workflow := &OrchestratorWorkflowModel{Id: types.StringValue("id-1")}
	var gateway OrchestratorWorkflowGatewayAPIModel
	_, _, diags := client.ReadIt(workflow, &gateway, workflow.ReadGatewayPath())
	if diags.HasError() {
		t.Fatalf("ReadIt(vro): %v", diags.Errors())
	}

	if want := []string{"/vco/api/tasks/task-123"}; !slices.Equal(vroPaths, want) {
		t.Errorf("Orchestrator received %v, want %v", vroPaths, want)
	}
	if want := []string{"/vro/workflows/id-1"}; !slices.Equal(ariaPaths, want) {
		t.Errorf("Aria received %v, want %v", ariaPaths, want)
	}
}

func TestAriaClientRoutesEverythingToHostWithoutVROHost(t *testing.T) {
	paths := []string{}
	aria := newFakeAPI(t, func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		writeJSONStatus(w, http.StatusOK, map[string]any{"id": "task-123", "state": "pending"})
	})

	client := newTestClient(t, aria.URL)
	if client.VROClient != nil {
		t.Fatal("VROClient is set, want nil when no Orchestrator host is configured")
	}

	var task OrchestratorTaskAPIModel
	if _, _, diags := client.ReadIt(taskModel("task-123"), &task); diags.HasError() {
		t.Fatalf("ReadIt(vco): %v", diags.Errors())
	}

	if want := []string{"/vco/api/tasks/task-123"}; !slices.Equal(paths, want) {
		t.Errorf("Aria received %v, want %v", paths, want)
	}
}

// newTestClientWithVROIntegration returns a client resolving the Orchestrator to address from the
// name of its integration.
func newTestClientWithVROIntegration(
	t *testing.T,
	host string,
	name string,
) (*AriaClient, diag.Diagnostics) {
	t.Helper()
	client := &AriaClient{
		Host:               host,
		VROIntegrationName: name,
		AccessToken:        "fake-token",
		OKAPICallsLogLevel: "DEBUG",
		KOAPICallsLogLevel: "WARN",
		Context:            t.Context(),
	}
	return client, client.Init()
}

// newFakeAriaWithIntegrations serves the integrations listing, everything else answering 404.
func newFakeAriaWithIntegrations(
	t *testing.T,
	entries ...IntegrationEntryAPIModel,
) *httptest.Server {
	t.Helper()
	return newFakeAPI(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/"+INTEGRATIONS_PATH {
			writeJSONStatus(w, http.StatusNotFound, map[string]any{})
			return
		}
		writeJSONStatus(w, http.StatusOK, fakeIntegrationsResponse(entries...))
	})
}

func TestAriaClientResolvesVROHostFromIntegrationName(t *testing.T) {
	aria := newFakeAriaWithIntegrations(t, embeddedIntegration, externalIntegration, abxIntegration)

	client, diags := newTestClientWithVROIntegration(t, aria.URL, "External Orchestrator")
	if diags.HasError() {
		t.Fatalf("AriaClient.Init: %v", diags.Errors())
	}
	if client.VROHost != externalIntegration.CustomProperties["hostName"] {
		t.Errorf(
			"VROHost = %q, want %q",
			client.VROHost, externalIntegration.CustomProperties["hostName"])
	}
	if client.VROClient == nil {
		t.Fatal("VROClient is nil, want the client addressing the resolved host")
	}
	if got := client.ClientForPath("vco/api/tasks"); got != client.VROClient {
		t.Error("vco calls are not addressed to the resolved Orchestrator")
	}
}

func TestAriaClientSkipsVROResolutionWithoutIntegrationName(t *testing.T) {
	// No name and no host is the embedded Orchestrator: nothing is looked up.
	aria := newFakeAriaWithIntegrations(t, embeddedIntegration, externalIntegration)

	client, diags := newTestClientWithVROIntegration(t, aria.URL, "")
	if diags.HasError() {
		t.Fatalf("AriaClient.Init: %v", diags.Errors())
	}
	if client.VROClient != nil {
		t.Error("VROClient is set, want nil when no Orchestrator is configured")
	}
}

func TestAriaClientResolvesVROHostUnknownIntegrationName(t *testing.T) {
	aria := newFakeAriaWithIntegrations(t, embeddedIntegration, externalIntegration)

	_, diags := newTestClientWithVROIntegration(t, aria.URL, "Nowhere Orchestrator")
	if !diags.HasError() {
		t.Fatal("Init succeeded, want an error naming the integrations to choose from")
	}
	detail := diags.Errors()[0].Detail()
	for _, name := range []string{"embedded-VRO", "External Orchestrator"} {
		if !strings.Contains(detail, name) {
			t.Errorf("error %q does not name candidate %q", detail, name)
		}
	}
}

func TestAriaClientCheckConfigRejectsConflictingVROConfiguration(t *testing.T) {
	client := &AriaClient{
		Host:               "https://aria.example.net",
		VROHost:            "https://vro.example.net",
		VROIntegrationName: "External Orchestrator",
		AccessToken:        "fake-token",
		Context:            t.Context(),
	}
	if diags := client.CheckConfig(); !diags.HasError() {
		t.Fatal("CheckConfig accepted both an Orchestrator host and an integration name")
	}
}

func TestAriaClientCheckConfigRejectsSchemelessVROHost(t *testing.T) {
	client := &AriaClient{
		Host:        "https://aria.example.net",
		VROHost:     "vro.example.net",
		AccessToken: "fake-token",
		Context:     t.Context(),
	}
	diags := client.CheckConfig()
	if !diags.HasError() {
		t.Fatal("CheckConfig accepted an Orchestrator host without a scheme")
	}
}

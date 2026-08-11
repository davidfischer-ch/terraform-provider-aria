// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Shared helpers for unit tests that exercise the client and resources against a fake Aria API,
// without requiring a real platform.

// newFakeAPI starts an httptest server running handler and closes it when the test ends.
func newFakeAPI(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return server
}

// newTestClient builds an AriaClient pointed at host, pre-authenticated with a fake access token
// so Init does not attempt a token exchange.
func newTestClient(t *testing.T, host string) *AriaClient {
	t.Helper()
	client := &AriaClient{
		Host:               host,
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

// writeJSONStatus writes a JSON response with the given status. The Content-Type must be set before
// WriteHeader, otherwise it is ignored and resty skips unmarshalling the body.
func writeJSONStatus(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

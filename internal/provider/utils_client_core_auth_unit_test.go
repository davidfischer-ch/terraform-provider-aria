// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetAccessTokenLegacySuccess(t *testing.T) {
	var gotMethod, gotPath, gotContentType, gotRefreshToken string
	server := newFakeAPI(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath, gotContentType = r.Method, r.URL.Path, r.Header.Get("Content-Type")
		var body map[string]string
		_ = json.NewDecoder(r.Body).Decode(&body)
		gotRefreshToken = body["refreshToken"]
		writeJSONStatus(w, http.StatusOK, map[string]string{
			"tokenType": "Bearer",
			"token":     "legacy-access-token",
		})
	})

	client := &AriaClient{
		Host:               server.URL,
		RefreshToken:       "legacy-refresh-token",
		OKAPICallsLogLevel: "DEBUG",
		KOAPICallsLogLevel: "WARN",
		Context:            t.Context(),
	}
	diags := client.Init()
	if diags.HasError() {
		t.Fatalf("Init: %v", diags.Errors())
	}

	if gotMethod != http.MethodPost || gotPath != "/iaas/api/login" {
		t.Errorf("unexpected request: %s %s", gotMethod, gotPath)
	}
	if gotContentType != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", gotContentType)
	}
	if gotRefreshToken != "legacy-refresh-token" {
		t.Errorf("sent refreshToken = %q, want %q", gotRefreshToken, "legacy-refresh-token")
	}
	if client.AccessToken != "legacy-access-token" {
		t.Errorf("AccessToken = %q, want %q", client.AccessToken, "legacy-access-token")
	}
}

func TestGetAccessTokenLegacyError(t *testing.T) {
	server := newFakeAPI(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSONStatus(w, http.StatusUnauthorized, map[string]string{"message": "invalid token"})
	})

	client := &AriaClient{
		Host:               server.URL,
		RefreshToken:       "bad-refresh-token",
		OKAPICallsLogLevel: "DEBUG",
		KOAPICallsLogLevel: "WARN",
		Context:            t.Context(),
	}
	diags := client.Init()
	if !diags.HasError() {
		t.Fatal("expected an error diagnostic for a rejected refresh token")
	}
	if client.AccessToken != "" {
		t.Errorf("AccessToken = %q, want empty on failure", client.AccessToken)
	}
}

func TestGetAccessTokenVCFSuccess(t *testing.T) {
	var gotMethod, gotPath, gotContentType, gotGrantType, gotRefreshToken string
	server := newFakeAPI(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath, gotContentType = r.Method, r.URL.Path, r.Header.Get("Content-Type")
		if err := r.ParseForm(); err != nil {
			t.Fatalf("ParseForm: %v", err)
		}
		gotGrantType = r.PostFormValue("grant_type")
		gotRefreshToken = r.PostFormValue("refresh_token")
		writeJSONStatus(w, http.StatusOK, map[string]string{"access_token": "vcf-access-token"})
	})

	client := &AriaClient{
		Host:               server.URL,
		Tenant:             "classic",
		RefreshToken:       "vcf-api-token",
		OKAPICallsLogLevel: "DEBUG",
		KOAPICallsLogLevel: "WARN",
		Context:            t.Context(),
	}
	diags := client.Init()
	if diags.HasError() {
		t.Fatalf("Init: %v", diags.Errors())
	}

	if gotMethod != http.MethodPost || gotPath != "/tm/oauth/tenant/classic/token" {
		t.Errorf("unexpected request: %s %s", gotMethod, gotPath)
	}
	if gotContentType != "application/x-www-form-urlencoded" {
		t.Errorf("Content-Type = %q, want application/x-www-form-urlencoded", gotContentType)
	}
	if gotGrantType != "refresh_token" {
		t.Errorf("grant_type = %q, want refresh_token", gotGrantType)
	}
	if gotRefreshToken != "vcf-api-token" {
		t.Errorf("refresh_token = %q, want %q", gotRefreshToken, "vcf-api-token")
	}
	if client.AccessToken != "vcf-access-token" {
		t.Errorf("AccessToken = %q, want %q", client.AccessToken, "vcf-access-token")
	}
}

func TestGetAccessTokenVCFError(t *testing.T) {
	server := newFakeAPI(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSONStatus(w, http.StatusUnauthorized, map[string]string{"message": "invalid token"})
	})

	client := &AriaClient{
		Host:               server.URL,
		Tenant:             "classic",
		RefreshToken:       "bad-api-token",
		OKAPICallsLogLevel: "DEBUG",
		KOAPICallsLogLevel: "WARN",
		Context:            t.Context(),
	}
	diags := client.Init()
	if !diags.HasError() {
		t.Fatal("expected an error diagnostic for a rejected API token")
	}
	if client.AccessToken != "" {
		t.Errorf("AccessToken = %q, want empty on failure", client.AccessToken)
	}
}

// TestGetAccessTokenPrefersVCFWhenTenantSet asserts that a non-empty Tenant always routes through
// the VCF 9 endpoint, never the legacy one.
func TestGetAccessTokenPrefersVCFWhenTenantSet(t *testing.T) {
	server := newFakeAPI(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/iaas/api/login" {
			t.Errorf("legacy endpoint %s must not be called when Tenant is set", r.URL.Path)
		}
		writeJSONStatus(w, http.StatusOK, map[string]string{"access_token": "vcf-access-token"})
	})

	client := &AriaClient{
		Host:               server.URL,
		Tenant:             "classic",
		RefreshToken:       "vcf-api-token",
		OKAPICallsLogLevel: "DEBUG",
		KOAPICallsLogLevel: "WARN",
		Context:            t.Context(),
	}
	diags := client.Init()
	if diags.HasError() {
		t.Fatalf("Init: %v", diags.Errors())
	}
	if client.AccessToken != "vcf-access-token" {
		t.Errorf("AccessToken = %q, want %q", client.AccessToken, "vcf-access-token")
	}
}

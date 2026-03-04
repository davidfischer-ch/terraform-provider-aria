// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"encoding/json"
	"testing"
)

func unmarshalOrFail[T any](t *testing.T, data []byte) T {
	t.Helper()
	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	return result
}

func assertString(t *testing.T, val any, expected string) {
	t.Helper()
	s, ok := val.(string)
	if !ok {
		t.Fatalf("expected string, got %T", val)
	}
	CheckEqual(t, s, expected)
}

func assertMap(t *testing.T, val any) map[string]any {
	t.Helper()
	m, ok := val.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", val)
	}
	return m
}

func TestRedactSensitiveKeys_SimpleObject(t *testing.T) {
	input := []byte(`{"refreshToken":"my-secret","host":"example.com"}`)
	result := redactJSON(input)
	data := unmarshalOrFail[map[string]any](t, result)
	assertString(t, data["refreshToken"], "<REDACTED>")
	assertString(t, data["host"], "example.com")
}

func TestRedactSensitiveKeys_Token(t *testing.T) {
	input := []byte(`{"tokenType":"Bearer","token":"abc123"}`)
	result := redactJSON(input)
	data := unmarshalOrFail[map[string]any](t, result)
	assertString(t, data["token"], "<REDACTED>")
	assertString(t, data["tokenType"], "Bearer")
}

func TestRedactSensitiveKeys_SystemCredentials(t *testing.T) {
	input := []byte(`{"name":"repo","systemCredentials":"p@ssw0rd","location":"/path"}`)
	result := redactJSON(input)
	data := unmarshalOrFail[map[string]any](t, result)
	assertString(t, data["systemCredentials"], "<REDACTED>")
	assertString(t, data["name"], "repo")
	assertString(t, data["location"], "/path")
}

func TestRedactSensitiveKeys_NestedObject(t *testing.T) {
	input := []byte(`{"outer":{"token":"secret","name":"test"}}`)
	result := redactJSON(input)
	data := unmarshalOrFail[map[string]any](t, result)
	inner := assertMap(t, data["outer"])
	assertString(t, inner["token"], "<REDACTED>")
	assertString(t, inner["name"], "test")
}

func TestRedactSensitiveKeys_Array(t *testing.T) {
	input := []byte(`[{"token":"s1","id":"1"},{"token":"s2","id":"2"}]`)
	result := redactJSON(input)
	data := unmarshalOrFail[[]any](t, result)
	assertString(t, assertMap(t, data[0])["token"], "<REDACTED>")
	assertString(t, assertMap(t, data[0])["id"], "1")
	assertString(t, assertMap(t, data[1])["token"], "<REDACTED>")
}

func TestRedactSensitiveKeys_NoSensitiveFields(t *testing.T) {
	input := []byte(`{"name":"test","value":"safe"}`)
	result := redactJSON(input)
	data := unmarshalOrFail[map[string]any](t, result)
	assertString(t, data["name"], "test")
	assertString(t, data["value"], "safe")
}

func TestRedactJSON_InvalidJSON(t *testing.T) {
	input := []byte(`not json`)
	result := redactJSON(input)
	CheckEqual(t, string(result), "not json")
}

func TestRedactJSON_NullLiteral(t *testing.T) {
	input := []byte(`null`)
	result := redactJSON(input)
	CheckEqual(t, string(result), "null")
}

func TestRedactJSON_EmptyBody(t *testing.T) {
	input := []byte(`<body>`)
	result := redactJSON(input)
	CheckEqual(t, string(result), "<body>")
}

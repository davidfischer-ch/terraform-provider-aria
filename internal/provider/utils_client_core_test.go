// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"encoding/json"
	"testing"
)

func TestRedactSensitiveKeys_SimpleObject(t *testing.T) {
	input := []byte(`{"refreshToken":"my-secret","host":"example.com"}`)
	result := redactJSON(input)
	var data map[string]any
	json.Unmarshal(result, &data)
	CheckEqual(t, data["refreshToken"].(string), "<REDACTED>")
	CheckEqual(t, data["host"].(string), "example.com")
}

func TestRedactSensitiveKeys_Token(t *testing.T) {
	input := []byte(`{"tokenType":"Bearer","token":"abc123"}`)
	result := redactJSON(input)
	var data map[string]any
	json.Unmarshal(result, &data)
	CheckEqual(t, data["token"].(string), "<REDACTED>")
	CheckEqual(t, data["tokenType"].(string), "Bearer")
}

func TestRedactSensitiveKeys_SystemCredentials(t *testing.T) {
	input := []byte(`{"name":"repo","systemCredentials":"p@ssw0rd","location":"/path"}`)
	result := redactJSON(input)
	var data map[string]any
	json.Unmarshal(result, &data)
	CheckEqual(t, data["systemCredentials"].(string), "<REDACTED>")
	CheckEqual(t, data["name"].(string), "repo")
	CheckEqual(t, data["location"].(string), "/path")
}

func TestRedactSensitiveKeys_NestedObject(t *testing.T) {
	input := []byte(`{"outer":{"token":"secret","name":"test"}}`)
	result := redactJSON(input)
	var data map[string]any
	json.Unmarshal(result, &data)
	inner := data["outer"].(map[string]any)
	CheckEqual(t, inner["token"].(string), "<REDACTED>")
	CheckEqual(t, inner["name"].(string), "test")
}

func TestRedactSensitiveKeys_Array(t *testing.T) {
	input := []byte(`[{"token":"s1","id":"1"},{"token":"s2","id":"2"}]`)
	result := redactJSON(input)
	var data []any
	json.Unmarshal(result, &data)
	CheckEqual(t, data[0].(map[string]any)["token"].(string), "<REDACTED>")
	CheckEqual(t, data[0].(map[string]any)["id"].(string), "1")
	CheckEqual(t, data[1].(map[string]any)["token"].(string), "<REDACTED>")
}

func TestRedactSensitiveKeys_NoSensitiveFields(t *testing.T) {
	input := []byte(`{"name":"test","value":"safe"}`)
	result := redactJSON(input)
	var data map[string]any
	json.Unmarshal(result, &data)
	CheckEqual(t, data["name"].(string), "test")
	CheckEqual(t, data["value"].(string), "safe")
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

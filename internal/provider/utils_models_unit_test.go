// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import "testing"

func TestStringOrNullValue(t *testing.T) {
	if got := StringOrNullValue(""); !got.IsNull() {
		t.Errorf("empty string must map to null, got %#v", got)
	}
	if got := StringOrNullValue("x"); got.ValueString() != "x" {
		t.Errorf("non-empty string = %q, want %q", got.ValueString(), "x")
	}
}

func TestJSONNormalizedFromAnyNil(t *testing.T) {
	got, diags := JSONNormalizedFromAny("field", nil)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags.Errors())
	}
	if !got.IsNull() {
		t.Error("nil value must map to a null normalized JSON")
	}
}

func TestJSONNormalizedRoundTrip(t *testing.T) {
	input := map[string]any{"key": "value", "count": float64(3)}

	normalized, diags := JSONNormalizedFromAny("field", input)
	if diags.HasError() {
		t.Fatalf("JSONNormalizedFromAny: %v", diags.Errors())
	}

	back, diags := JSONNormalizedToAny(normalized)
	if diags.HasError() {
		t.Fatalf("JSONNormalizedToAny: %v", diags.Errors())
	}

	got, ok := back.(map[string]any)
	if !ok {
		t.Fatalf("round trip type = %T, want map[string]any", back)
	}
	if got["key"] != "value" || got["count"] != float64(3) {
		t.Errorf("round trip mismatch: got %#v", got)
	}
}

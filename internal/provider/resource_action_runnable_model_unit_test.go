// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Happy path: a fully populated runnable, including a valid JSON input_bindings value and
// non-empty parameter lists, round-trips through ToAPI with no diagnostics.
func TestResourceActionRunnableModelToAPI(t *testing.T) {
	runnable := ResourceActionRunnableModel{
		Id:            types.StringValue("c974e486-9039-4b84-9152-0e5aa2074d26"),
		Name:          types.StringValue("someAction"),
		Type:          types.StringValue("abx.action"),
		ProjectId:     types.StringValue("175bed78-dd9e-4999-8669-cc62388e9abb"),
		EndpointLink:  types.StringValue("/resources/endpoints/8a430db3-924c-4d58-a29a-da811f9c992e"),
		InputBindings: jsontypes.NewNormalizedValue(`{"foo":"bar"}`),
		InputParameters: []ParameterModel{
			{Name: types.StringValue("resourceId"), Type: types.StringValue("string")},
		},
		OutputParameters: []ParameterModel{},
	}

	raw, diags := runnable.ToAPI(t.Context())
	CheckDiagnostics(t, diags, "", "")

	CheckEqual(t, raw.Id, "c974e486-9039-4b84-9152-0e5aa2074d26")
	CheckEqual(t, raw.Name, "someAction")
	CheckEqual(t, raw.Type, "abx.action")
	CheckEqual(t, raw.ProjectId, "175bed78-dd9e-4999-8669-cc62388e9abb")
	CheckEqual(t, raw.EndpointLink, "/resources/endpoints/8a430db3-924c-4d58-a29a-da811f9c992e")
	CheckDeepEqual(t, raw.InputBindings, map[string]any{"foo": "bar"})
	CheckDeepEqual(t, raw.InputParameters, []ParameterAPIModel{
		{Name: "resourceId", Type: "string"},
	})
	CheckDeepEqual(t, raw.OutputParameters, []ParameterAPIModel{})
}

// A null input_bindings must be dropped from the request entirely (thanks to omitempty), not sent
// as JSON null.
func TestResourceActionRunnableModelToAPIInputBindingsNull(t *testing.T) {
	runnable := ResourceActionRunnableModel{
		InputBindings:    jsontypes.NewNormalizedNull(),
		InputParameters:  []ParameterModel{},
		OutputParameters: []ParameterModel{},
	}

	raw, diags := runnable.ToAPI(t.Context())
	CheckDiagnostics(t, diags, "", "")

	if raw.InputBindings != nil {
		t.Errorf("InputBindings = %#v, want nil", raw.InputBindings)
	}
}

// Regression test for aec55bc: an unknown input_bindings (the realistic state for a Computed and
// Optional attribute the practitioner never set) must not be passed to JSONNormalizedToAny.
// Unmarshal on an unknown jsontypes.Normalized errors.
func TestResourceActionRunnableModelToAPIInputBindingsUnknown(t *testing.T) {
	runnable := ResourceActionRunnableModel{
		InputBindings:    jsontypes.NewNormalizedUnknown(),
		InputParameters:  []ParameterModel{},
		OutputParameters: []ParameterModel{},
	}

	raw, diags := runnable.ToAPI(t.Context())
	if diags.HasError() {
		t.Fatalf("ToAPI: %v", diags.Errors())
	}

	if raw.InputBindings != nil {
		t.Errorf("InputBindings = %#v, want nil", raw.InputBindings)
	}
}

// Only Null and Unknown are special-cased; a non-null, non-unknown but malformed JSON value must
// still surface as an error, proving ToAPI didn't accidentally widen the guard to swallow all
// failures.
func TestResourceActionRunnableModelToAPIInputBindingsInvalidJSON(t *testing.T) {
	runnable := ResourceActionRunnableModel{
		InputBindings:    jsontypes.NewNormalizedValue("not valid json"),
		InputParameters:  []ParameterModel{},
		OutputParameters: []ParameterModel{},
	}

	_, diags := runnable.ToAPI(t.Context())
	if !diags.HasError() {
		t.Fatal("expected ToAPI to report an error for malformed input_bindings JSON")
	}
}

// Real API behavior: the response simply omits inputBindings when there are none
// (raw.InputBindings is the Go zero value nil for the `any` field), not an artificial edge case.
// FromAPI must decode this to Null, not error or produce an empty-string JSON value.
func TestResourceActionRunnableModelFromAPIInputBindingsAbsent(t *testing.T) {
	runnable := ResourceActionRunnableModel{}
	diags := runnable.FromAPI(t.Context(), ResourceActionRunnableAPIModel{
		Id:   "c974e486-9039-4b84-9152-0e5aa2074d26",
		Name: "someAction",
		Type: "abx.action",
	})
	CheckDiagnostics(t, diags, "", "")

	if !runnable.InputBindings.IsNull() {
		t.Errorf("InputBindings = %#v, want null", runnable.InputBindings)
	}
}

// A populated inputBindings value round-trips through FromAPI into an equivalent JSON-encoded
// string, compared by decoding both sides back to Go values, not raw string equality, because
// json.Marshal's key ordering is not guaranteed to match the input literal.
func TestResourceActionRunnableModelFromAPIInputBindingsPresent(t *testing.T) {
	runnable := ResourceActionRunnableModel{}
	diags := runnable.FromAPI(t.Context(), ResourceActionRunnableAPIModel{
		Id:            "c974e486-9039-4b84-9152-0e5aa2074d26",
		Name:          "someAction",
		Type:          "abx.action",
		InputBindings: map[string]any{"foo": "bar", "count": float64(2)},
	})
	CheckDiagnostics(t, diags, "", "")

	var decoded map[string]any
	if diags := runnable.InputBindings.Unmarshal(&decoded); diags.HasError() {
		t.Fatalf("InputBindings.Unmarshal: %v", diags.Errors())
	}
	CheckDeepEqual(t, decoded, map[string]any{"foo": "bar", "count": float64(2)})
}

// String() must include id, name, and project_id for log/error-message identification.
func TestResourceActionRunnableModelString(t *testing.T) {
	runnable := ResourceActionRunnableModel{
		Id:        types.StringValue("c974e486-9039-4b84-9152-0e5aa2074d26"),
		Name:      types.StringValue("someAction"),
		ProjectId: types.StringValue("175bed78-dd9e-4999-8669-cc62388e9abb"),
	}

	got := runnable.String()
	want := "Resource Action Runnable c974e486-9039-4b84-9152-0e5aa2074d26 (someAction) project " +
		"175bed78-dd9e-4999-8669-cc62388e9abb"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

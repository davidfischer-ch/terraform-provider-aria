// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// resourceActionCreateResult captures what the fake API observed and what Create produced.
type resourceActionCreateResult struct {
	requestHasInputBindings bool
	requestInputBindings    any
	hasError                bool
	errDetail               string
	stateSet                bool
	inputBindingsNull       bool
	inputBindings           map[string]any
}

// runResourceActionCreate drives the real resource Create (native, non-custom path) against a fake
// Aria API, planning a runnable with the given input_bindings. When respondWithInputBindings is
// true, the fake API echoes back requestInputBindings (or, if nil, a fixed value); otherwise the
// response omits inputBindings entirely, mirroring the real API's behavior when there are none.
func runResourceActionCreate(
	t *testing.T,
	inputBindings jsontypes.Normalized,
	respondWithInputBindings bool,
) resourceActionCreateResult {
	t.Helper()
	ctx := t.Context()

	var res resourceActionCreateResult

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/form-service/api/custom/resource-actions" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		body, _ := io.ReadAll(r.Body)
		var parsed map[string]any
		_ = json.Unmarshal(body, &parsed)

		runnableItem, _ := parsed["runnableItem"].(map[string]any)
		requestInputBindings, present := runnableItem["inputBindings"]
		res.requestHasInputBindings = present
		res.requestInputBindings = requestInputBindings

		responseRunnable := ResourceActionRunnableAPIModel{
			Id:               "c974e486-9039-4b84-9152-0e5aa2074d26",
			Name:             "someAction",
			Type:             "abx.action",
			ProjectId:        "175bed78-dd9e-4999-8669-cc62388e9abb",
			InputParameters:  []ParameterAPIModel{},
			OutputParameters: []ParameterAPIModel{},
		}
		if respondWithInputBindings {
			if requestInputBindings != nil {
				responseRunnable.InputBindings = requestInputBindings
			} else {
				responseRunnable.InputBindings = map[string]any{"foo": "bar"}
			}
		}

		responseBody := ResourceActionAPIModel{
			Id:           "resource-action-1",
			Name:         "someResourceAction",
			DisplayName:  "Some Resource Action",
			Description:  "desc",
			ProviderName: "xaas",
			ResourceType: "Custom.MyCustom",
			RunnableItem: responseRunnable,
			Status:       "RELEASED",
			ProjectId:    "175bed78-dd9e-4999-8669-cc62388e9abb",
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(responseBody)
	}))
	defer srv.Close()

	client := &AriaClient{
		Host:               srv.URL,
		AccessToken:        "fake-token",
		OKAPICallsLogLevel: "DEBUG",
		KOAPICallsLogLevel: "WARN",
		Context:            ctx,
	}
	if diags := client.Init(); diags.HasError() {
		t.Fatalf("AriaClient.Init: %v", diags.Errors())
	}

	res0, ok := NewResourceActionResource().(*ResourceActionResource)
	if !ok {
		t.Fatal("NewResourceActionResource() did not return *ResourceActionResource")
	}
	res0.client = client

	schema := ResourceActionSchema()

	model := ResourceActionModel{
		Id:           types.StringNull(),
		Name:         types.StringValue("someResourceAction"),
		DisplayName:  types.StringValue("Some Resource Action"),
		Description:  types.StringValue("desc"),
		ProviderName: types.StringValue("xaas"),
		ResourceId:   types.StringNull(),
		ResourceType: types.StringValue("Custom.MyCustom"),
		RunnableItem: ResourceActionRunnableModel{
			Id:               types.StringValue("c974e486-9039-4b84-9152-0e5aa2074d26"),
			Name:             types.StringValue("someAction"),
			Type:             types.StringValue("abx.action"),
			ProjectId:        types.StringValue("175bed78-dd9e-4999-8669-cc62388e9abb"),
			EndpointLink:     types.StringValue(""),
			InputBindings:    inputBindings,
			InputParameters:  []ParameterModel{},
			OutputParameters: []ParameterModel{},
		},
		Criteria:       jsontypes.NewNormalizedNull(),
		Status:         types.StringValue("RELEASED"),
		FormDefinition: types.ObjectNull(CustomFormModel{}.AttributeTypes()),
		ProjectId:      types.StringValue("175bed78-dd9e-4999-8669-cc62388e9abb"),
		OrgId:          types.StringNull(),
	}

	plan := tfsdk.Plan{Schema: schema}
	if diags := plan.Set(ctx, &model); diags.HasError() {
		t.Fatalf("plan.Set: %v", diags.Errors())
	}

	req := resource.CreateRequest{Plan: plan}
	resp := &resource.CreateResponse{State: tfsdk.State{Schema: schema}}
	res0.Create(ctx, req, resp)

	res.hasError = resp.Diagnostics.HasError()
	for _, e := range resp.Diagnostics.Errors() {
		res.errDetail += e.Summary() + ": " + e.Detail() + "; "
	}

	res.stateSet = !resp.State.Raw.IsNull()
	if res.stateSet {
		var out ResourceActionModel
		if diags := resp.State.Get(ctx, &out); diags.HasError() {
			t.Fatalf("state.Get: %v", diags.Errors())
		}
		res.inputBindingsNull = out.RunnableItem.InputBindings.IsNull()
		if !res.inputBindingsNull {
			var decoded map[string]any
			if diags := out.RunnableItem.InputBindings.Unmarshal(&decoded); diags.HasError() {
				t.Fatalf("persisted InputBindings.Unmarshal: %v", diags.Errors())
			}
			res.inputBindings = decoded
		}
	}

	return res
}

// End-to-end regression test for aec55bc: at apply time a Computed+Optional input_bindings the
// practitioner never configured arrives in the plan as Unknown, not Null. ToAPI must still omit
// it from the create request, and when the fake API's response likewise omits inputBindings, the
// persisted state must end up Null (not an error, not a stray empty-string).
func TestResourceActionResourceCreateInputBindingsUnknown(t *testing.T) {
	res := runResourceActionCreate(t, jsontypes.NewNormalizedUnknown(), false)

	if res.hasError {
		t.Fatalf("unexpected Create error: %s", res.errDetail)
	}
	if res.requestHasInputBindings {
		t.Errorf("create request runnableItem has inputBindings = %#v, want key absent",
			res.requestInputBindings)
	}
	if !res.stateSet {
		t.Fatal("state was not set")
	}
	if !res.inputBindingsNull {
		t.Errorf("persisted input_bindings = %#v, want null", res.inputBindings)
	}
}

// A practitioner-supplied input_bindings value must survive the create round trip verbatim: sent
// as decoded JSON in the request body's runnableItem.inputBindings, and read back into state from
// the API's response.
func TestResourceActionResourceCreateInputBindingsSet(t *testing.T) {
	res := runResourceActionCreate(t, jsontypes.NewNormalizedValue(`{"foo":"bar"}`), true)

	if res.hasError {
		t.Fatalf("unexpected Create error: %s", res.errDetail)
	}
	if !res.requestHasInputBindings {
		t.Fatal("create request runnableItem has no inputBindings key, want it present")
	}
	CheckDeepEqual(t, res.requestInputBindings, map[string]any{"foo": "bar"})

	if !res.stateSet {
		t.Fatal("state was not set")
	}
	if res.inputBindingsNull {
		t.Fatal("persisted input_bindings is null, want the round-tripped value")
	}
	CheckDeepEqual(t, res.inputBindings, map[string]any{"foo": "bar"})
}

// Explicit Null behaves identically to Unknown at the wire level: both are the two guarded cases
// in ToAPI's fix. This locks in that Null alone (not just Unknown) also produces a keyless request
// body.
func TestResourceActionResourceCreateInputBindingsNull(t *testing.T) {
	res := runResourceActionCreate(t, jsontypes.NewNormalizedNull(), false)

	if res.hasError {
		t.Fatalf("unexpected Create error: %s", res.errDetail)
	}
	if res.requestHasInputBindings {
		t.Errorf("create request runnableItem has inputBindings = %#v, want key absent",
			res.requestInputBindings)
	}
	if !res.stateSet {
		t.Fatal("state was not set")
	}
	if !res.inputBindingsNull {
		t.Errorf("persisted input_bindings = %#v, want null", res.inputBindings)
	}
}

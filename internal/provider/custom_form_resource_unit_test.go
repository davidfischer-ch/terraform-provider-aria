// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type customFormCreateResult struct {
	hasError  bool
	errDetail string
	stateSet  bool
	id        string
	createdId string
}

// createCustomForm runs Create against a fake API minting createdId for the form. The form is
// readable back at that identifier only, the way an API ignoring the generated one behaves.
func createCustomForm(t *testing.T, createdId string, readable bool) customFormCreateResult {
	t.Helper()

	ctx := t.Context()
	res := customFormCreateResult{}

	server := newFakeAPI(t, func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		// No form exists yet for that source
		case strings.HasSuffix(path, "/form-service/api/forms/fetchBySourceAndType"):
			writeJSONStatus(w, http.StatusNotFound, map[string]any{})

		case r.Method == http.MethodPost && strings.HasSuffix(path, "/form-service/api/forms"):
			var body CustomFormAPIModel
			_ = json.NewDecoder(r.Body).Decode(&body)
			res.createdId = body.Id
			writeJSONStatus(w, http.StatusCreated, CustomFormAPIModel{
				Id:         createdId,
				Name:       body.Name,
				Type:       body.Type,
				Form:       body.Form,
				FormFormat: body.FormFormat,
				SourceId:   body.SourceId,
				SourceType: body.SourceType,
				Status:     body.Status,
			})

		case r.Method == http.MethodGet &&
			strings.HasSuffix(path, "/form-service/api/forms/"+createdId):
			if !readable {
				writeJSONStatus(w, http.StatusNotFound, map[string]any{})
				return
			}
			writeJSONStatus(w, http.StatusOK, CustomFormAPIModel{
				Id:         createdId,
				Name:       "someForm",
				Type:       "requestForm",
				Form:       "{}",
				FormFormat: "JSON",
				SourceId:   "8a7480d3-8e53-5332-018e-857e0d4f3437",
				SourceType: "com.vmw.blueprint",
				Status:     "ON",
			})

		default:
			writeJSONStatus(w, http.StatusNotFound, map[string]any{})
		}
	})

	res0, ok := NewCustomFormResource().(*CustomFormResource)
	if !ok {
		t.Fatal("NewCustomFormResource() did not return *CustomFormResource")
	}
	res0.client = newTestClient(t, server.URL)

	schema := CustomFormSchema()

	model := CustomFormModel{
		Id:         types.StringNull(),
		Name:       types.StringValue("someForm"),
		Type:       types.StringValue("requestForm"),
		Form:       jsontypes.NewNormalizedValue("{}"),
		FormFormat: types.StringValue("JSON"),
		Styles:     types.StringValue(""),
		SourceId:   types.StringValue("8a7480d3-8e53-5332-018e-857e0d4f3437"),
		SourceType: types.StringValue("com.vmw.blueprint"),
		Tenant:     types.StringNull(),
		Status:     types.StringValue("ON"),
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
		var out CustomFormModel
		if diags := resp.State.Get(ctx, &out); diags.HasError() {
			t.Fatalf("state.Get: %v", diags.Errors())
		}
		res.id = out.Id.ValueString()
	}

	return res
}

// The API is free to mint its own identifier, reading the form back requires the one it kept.
func TestCustomFormResourceCreateAdoptsAPIIdentifier(t *testing.T) {
	res := createCustomForm(t, "b4d3f6d0-6c1f-4f0e-9a7b-8bb0c8f4e5a1", true)
	if res.hasError {
		t.Fatalf("unexpected error: %s", res.errDetail)
	}
	if !res.stateSet {
		t.Fatal("no state saved after create")
	}
	if res.id != "b4d3f6d0-6c1f-4f0e-9a7b-8bb0c8f4e5a1" {
		t.Errorf("id = %q, want the identifier returned by the API", res.id)
	}
	if len(res.createdId) == 0 {
		t.Error("create request carries no identifier, one is generated when the form is new")
	}
}

// Terraform rejects an apply that saves no state, an error must be reported instead.
func TestCustomFormResourceCreateMissingAfterCreate(t *testing.T) {
	res := createCustomForm(t, "b4d3f6d0-6c1f-4f0e-9a7b-8bb0c8f4e5a1", false)
	if !res.hasError {
		t.Fatal("create succeeded, want an error when the form cannot be read back")
	}
	if !strings.Contains(res.errDetail, "read") {
		t.Errorf("error %q does not report the failing read back", res.errDetail)
	}
}

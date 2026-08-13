// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Regression: vRO's Task API echoes description back empty. FromAPI must keep the prior/planned
// description instead of the API's blank one.
func TestOrchestratorTaskModelFromAPIPreservesDescriptionDespiteAPI(t *testing.T) {
	ctx := t.Context()

	priorList, diags := ParameterModelListFromAPI(ctx, []ParameterAPIModel{
		{Name: "vraHost", Type: "VRA:Host", Description: "Target vRA host"},
	})
	if diags.HasError() {
		t.Fatalf("prior input parameters list: %v", diags.Errors())
	}

	model := OrchestratorTaskModel{InputParameters: priorList}
	fromAPIDiags := model.FromAPI(ctx, OrchestratorTaskAPIModel{
		RecurrenceStartDate: "2050-01-06T05:02:00Z",
		Workflow:            OrchestratorTaskWorkflowAPIModel{Id: "wf-1", Name: "wf"},
		// The real API's behavior being regression-tested: description comes back empty even
		// though "Target vRA host" was sent on create.
		InputParameters: []ParameterAPIModel{
			{Name: "vraHost", Type: "VRA:Host", Description: ""},
		},
	})
	if fromAPIDiags.HasError() {
		t.Fatalf("FromAPI: %v", fromAPIDiags.Errors())
	}

	var out []ParameterModel
	if diags := model.InputParameters.ElementsAs(ctx, &out, false); diags.HasError() {
		t.Fatalf("input parameters elements: %v", diags.Errors())
	}
	if len(out) != 1 {
		t.Fatalf("input_parameters has %d entries, want 1", len(out))
	}
	if got := out[0].Description.ValueString(); got != "Target vRA host" {
		t.Errorf("description = %q, want %q (preserved despite the API's blank response)",
			got, "Target vRA host")
	}
	if got := out[0].Type.ValueString(); got != "VRA:Host" {
		t.Errorf("type = %q, want %q (still sourced from the API)", got, "VRA:Host")
	}
}

// A parameter added out-of-band (unknown to priorList) must fall back to the API's description
// instead of panicking.
func TestOrchestratorTaskModelFromAPIInputParametersUnknownToPrior(t *testing.T) {
	ctx := t.Context()

	model := OrchestratorTaskModel{InputParameters: types.ListValueMust(
		types.ObjectType{AttrTypes: ParameterModel{}.AttributeTypes()}, []attr.Value{})}

	diags := model.FromAPI(ctx, OrchestratorTaskAPIModel{
		RecurrenceStartDate: "2050-01-06T05:02:00Z",
		Workflow:            OrchestratorTaskWorkflowAPIModel{Id: "wf-1", Name: "wf"},
		InputParameters: []ParameterAPIModel{
			{Name: "driftedParam", Type: "string", Description: ""},
		},
	})
	if diags.HasError() {
		t.Fatalf("FromAPI: %v", diags.Errors())
	}

	var out []ParameterModel
	if diags := model.InputParameters.ElementsAs(ctx, &out, false); diags.HasError() {
		t.Fatalf("input parameters elements: %v", diags.Errors())
	}
	if len(out) != 1 || out[0].Name.ValueString() != "driftedParam" {
		t.Fatalf("input_parameters = %#v, want a single driftedParam entry", out)
	}
}

// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func CheckDiagnostics(
	t *testing.T,
	diags diag.Diagnostics,
	warningMessage string,
	errorMessage string,
) {

	warnings := diags.Warnings()
	if warningMessage != "" {
		if len(warnings) == 0 {
			t.Errorf("Message \"%s\" not found, there are no warnings.", warningMessage)
		} else {
			detail := warnings[len(warnings)-1].Detail()
			if strings.Contains(detail, warningMessage) {
				warnings = diag.Diagnostics{} // Warnings are processed
			} else {
				t.Errorf(
					"Message \"%s\" not found in latest warning \"%s\".",
					warningMessage, detail)
			}
		}
	}

	errors := diags.Errors()
	if errorMessage != "" {
		if len(errors) == 0 {
			t.Errorf("Message \"%s\" not found, there are no errors.", errorMessage)
		} else {
			detail := errors[len(errors)-1].Detail()
			if strings.Contains(detail, errorMessage) {
				errors = diag.Diagnostics{} // Errors are processed
			} else {
				t.Errorf("Message \"%s\" not found in latest error \"%s\".", errorMessage, detail)
			}
		}
	}

	if len(warnings) > 0 {
		t.Errorf("Diagnostics contains unexpected warnings.")
		for counter, warning := range warnings {
			t.Errorf("Diagnostic Warning %d - %s", counter+1, warning.Detail())
		}
	}

	if len(errors) > 0 {
		t.Errorf("Diagnostics contains unexpected errors.")
		for counter, error := range errors {
			t.Errorf("Diagnostic Error %d - %s", counter+1, error.Detail())
		}
	}
}

func CheckEqual(t *testing.T, actual interface{}, expected interface{}) {
	if actual != expected {
		t.Errorf("Result was incorrect,\nactual  : %s\nexpected: %s", actual, expected)
	}
}

func CheckDeepEqual(t *testing.T, actual interface{}, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Result was incorrect,\nactual  : %s\nexpected: %s", actual, expected)
	}
}

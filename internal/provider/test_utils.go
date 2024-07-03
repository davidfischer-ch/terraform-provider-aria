// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func CheckDiagnostics(t *testing.T, diags diag.Diagnostics, errorMessage string) {
	errors := diags.Errors()
	if errorMessage != "" {
		detail := errors[len(errors)-1].Detail()
		if strings.Contains(detail, errorMessage) {
			return
		}
		t.Errorf("Message \"%s\" not found in latest error detail \"%s\".", errorMessage, detail)
	}
	if len(errors) > 0 {
		t.Errorf("Diagnostics contains unexpected errors.")
		for counter, error := range diags.Errors() {
			t.Errorf("Diagnostic Error %d - %s", counter, error.Detail())
		}
	}
}

func CheckEqual(t *testing.T, actual interface{}, expected interface{}) {
	if actual != expected {
		t.Errorf("Result was incorrect, got: %s, expected %s", actual, expected)
	}
}

func CheckDeepEqual(t *testing.T, actual interface{}, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Result was incorrect, got: %s, expected %s", actual, expected)
	}
}

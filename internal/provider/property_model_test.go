// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestPropertyModelDefaultToAPI(t *testing.T) {
	cases := []struct {
		name             string
		propertyType     string
		propertyInternal string
		propertyRaw      any
		warningMessage   string
		errorMessage     string
	}{
		{
			name:             "boolean value",
			propertyType:     "boolean",
			propertyInternal: "false",
			propertyRaw:      false,
		},
		{
			name:             "boolean value",
			propertyType:     "boolean",
			propertyInternal: "true",
			propertyRaw:      true,
		},
		{
			name:             "wrong boolean value",
			propertyType:     "boolean",
			propertyInternal: "not really a boolean",
			propertyRaw:      nil,
			errorMessage:     "invalid syntax",
		},
		{
			name:             "integer value",
			propertyType:     "integer",
			propertyInternal: "42",
			propertyRaw:      int64(42),
		},
		{
			name:             "wrong integer value",
			propertyType:     "integer",
			propertyInternal: "1.2",
			propertyRaw:      nil,
			errorMessage:     "invalid syntax",
		},
		{
			name:             "number value (integer)",
			propertyType:     "number",
			propertyInternal: "-100",
			propertyRaw:      int64(-100),
		},
		{
			name:             "number value (float)",
			propertyType:     "number",
			propertyInternal: "3.141592",
			propertyRaw:      3.141592,
		},
		{
			name:             "string value",
			propertyType:     "string",
			propertyInternal: "some text",
			propertyRaw:      "some text",
		},
		{
			name:             "array value",
			propertyType:     "array",
			propertyInternal: "[1, 2, 3]",
			errorMessage:     "type array is not yet implemented",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			property := PropertyModel{
				Title:   types.StringValue("P"),
				Type:    types.StringValue(tc.propertyType),
				Default: types.StringValue(tc.propertyInternal),
			}
			raw, diags := property.ToAPI(context.Background())
			CheckDiagnostics(t, diags, tc.warningMessage, tc.errorMessage)
			CheckEqual(t, raw.Default, tc.propertyRaw)
		})
	}
}

func TestPropertyModelDefaultFromAPI(t *testing.T) {
	cases := []struct {
		name             string
		propertyType     string
		propertyRaw      any
		propertyInternal string
		warningMessage   string
		errorMessage     string
	}{
		{
			name:             "boolean value",
			propertyType:     "boolean",
			propertyRaw:      false,
			propertyInternal: "false",
		},
		{
			name:             "boolean value",
			propertyType:     "boolean",
			propertyRaw:      true,
			propertyInternal: "true",
		},
		{
			name:             "wrong boolean value",
			propertyType:     "boolean",
			propertyRaw:      "not really a boolean",
			propertyInternal: "not really a boolean",
			warningMessage:   "Property P default \"not really a boolean\" is not a boolean",
		},
		{
			name:             "integer value",
			propertyType:     "integer",
			propertyRaw:      int64(42),
			propertyInternal: "42",
		},
		{
			name:             "wrong integer value",
			propertyType:     "integer",
			propertyRaw:      1.2,
			propertyInternal: "%!s(float64=1.2)",
			warningMessage:   "Property P default \"%!s(float64=1.2)\" is not an integer",
		},
		{
			name:             "number value (integer)",
			propertyType:     "number",
			propertyRaw:      int64(-100),
			propertyInternal: "-100",
		},
		{
			name:             "number value (float)",
			propertyType:     "number",
			propertyRaw:      3.141592,
			propertyInternal: "3.141592",
		},
		{
			name:             "string value",
			propertyType:     "string",
			propertyRaw:      "some text",
			propertyInternal: "some text",
		},
		{
			name:             "array value",
			propertyType:     "array",
			propertyRaw:      []int{1, 2, 3},
			propertyInternal: "[%!s(int=1) %!s(int=2) %!s(int=3)]",
			errorMessage:     "type array is not yet implemented",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			property := PropertyModel{}
			diags := property.FromAPI(context.Background(), PropertyAPIModel{
				Title:   "P",
				Type:    tc.propertyType,
				Default: tc.propertyRaw,
			})
			CheckDiagnostics(t, diags, tc.warningMessage, tc.errorMessage)
			CheckEqual(t, property.Default.ValueString(), tc.propertyInternal)
		})
	}
}

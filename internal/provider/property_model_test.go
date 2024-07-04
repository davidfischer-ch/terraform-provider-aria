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
		propertyInternal types.String
		propertyRaw      any
		warningMessage   string
		errorMessage     string
	}{
		{
			name:             "boolean value (false)",
			propertyType:     "boolean",
			propertyInternal: types.StringValue("false"),
			propertyRaw:      false,
		},
		{
			name:             "boolean value (true)",
			propertyType:     "boolean",
			propertyInternal: types.StringValue("true"),
			propertyRaw:      true,
		},
		{
			name:             "boolean value (string)",
			propertyType:     "boolean",
			propertyInternal: types.StringValue("not really a boolean"),
			propertyRaw:      nil,
			errorMessage:     "invalid syntax",
		},
		{
			name:             "boolean value (nil)",
			propertyType:     "boolean",
			propertyInternal: types.StringNull(),
			propertyRaw:      nil,
		},
		{
			name:             "integer value (integer)",
			propertyType:     "integer",
			propertyInternal: types.StringValue("42"),
			propertyRaw:      int64(42),
		},
		{
			name:             "integer value (float)",
			propertyType:     "integer",
			propertyInternal: types.StringValue("1.2"),
			propertyRaw:      nil,
			errorMessage:     "invalid syntax",
		},
		{
			name:             "number value (integer)",
			propertyType:     "number",
			propertyInternal: types.StringValue("-100"),
			propertyRaw:      int64(-100),
		},
		{
			name:             "number value (float)",
			propertyType:     "number",
			propertyInternal: types.StringValue("3.141592"),
			propertyRaw:      3.141592,
		},
		{
			name:             "string value (string)",
			propertyType:     "string",
			propertyInternal: types.StringValue("some text"),
			propertyRaw:      "some text",
		},
		{
			name:             "array value (array)",
			propertyType:     "array",
			propertyInternal: types.StringValue("[1, 2, 3]"),
			errorMessage:     "type array is not yet implemented",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			property := PropertyModel{
				Name:    types.StringValue("p"),
				Title:   types.StringValue("P"),
				Type:    types.StringValue(tc.propertyType),
				Default: tc.propertyInternal,
			}
			name, raw, diags := property.ToAPI(context.Background())
			CheckDiagnostics(t, diags, tc.warningMessage, tc.errorMessage)
			CheckEqual(t, raw.Default, tc.propertyRaw)
			CheckEqual(t, name, "p")
		})
	}
}

func TestPropertyModelDefaultFromAPI(t *testing.T) {
	cases := []struct {
		name             string
		propertyType     string
		propertyRaw      any
		propertyInternal types.String
		warningMessage   string
		errorMessage     string
	}{
		{
			name:             "boolean value (false)",
			propertyType:     "boolean",
			propertyRaw:      false,
			propertyInternal: types.StringValue("false"),
		},
		{
			name:             "boolean value (true)",
			propertyType:     "boolean",
			propertyRaw:      true,
			propertyInternal: types.StringValue("true"),
		},
		{
			name:             "boolean value (string)",
			propertyType:     "boolean",
			propertyRaw:      "not really a boolean",
			propertyInternal: types.StringValue("not really a boolean"),
			warningMessage:   "Property P default \"not really a boolean\" is not a boolean",
		},
		{
			name:             "integer value (integer)",
			propertyType:     "integer",
			propertyRaw:      int64(42),
			propertyInternal: types.StringValue("42"),
		},
		{
			name:             "integer value (float)",
			propertyType:     "integer",
			propertyRaw:      1.2,
			propertyInternal: types.StringValue("%!s(float64=1.2)"),
			warningMessage:   "Property P default \"%!s(float64=1.2)\" is not an integer",
		},
		{
			name:             "number value (integer)",
			propertyType:     "number",
			propertyRaw:      int64(-100),
			propertyInternal: types.StringValue("-100"),
		},
		{
			name:             "number value (float)",
			propertyType:     "number",
			propertyRaw:      3.141592,
			propertyInternal: types.StringValue("3.141592"),
		},
		{
			name:             "number value (nil)",
			propertyType:     "number",
			propertyRaw:      nil,
			propertyInternal: types.StringNull(),
		},
		{
			name:             "string value (string)",
			propertyType:     "string",
			propertyRaw:      "some text",
			propertyInternal: types.StringValue("some text"),
		},
		{
			name:             "array value (array)",
			propertyType:     "array",
			propertyRaw:      []int{1, 2, 3},
			propertyInternal: types.StringValue("[%!s(int=1) %!s(int=2) %!s(int=3)]"),
			errorMessage:     "type array is not yet implemented",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			property := PropertyModel{}
			diags := property.FromAPI(context.Background(), "p", PropertyAPIModel{
				Title:   "P",
				Type:    tc.propertyType,
				Default: tc.propertyRaw,
			})
			CheckDiagnostics(t, diags, tc.warningMessage, tc.errorMessage)
			CheckEqual(t, property.Default, tc.propertyInternal)
			CheckEqual(t, property.Name.ValueString(), "p")
		})
	}
}

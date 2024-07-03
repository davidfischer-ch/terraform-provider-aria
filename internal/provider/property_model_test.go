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
		errorMessage     string
	}{
		{
			name:             "boolean value",
			propertyType:     "boolean",
			propertyInternal: "false",
			propertyRaw:      false,
			errorMessage:     "",
		},
		{
			name:             "boolean value",
			propertyType:     "boolean",
			propertyInternal: "true",
			propertyRaw:      true,
			errorMessage:     "",
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
			errorMessage:     "",
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
			errorMessage:     "",
		},
		{
			name:             "number value (float)",
			propertyType:     "number",
			propertyInternal: "3.141592",
			propertyRaw:      3.141592,
			errorMessage:     "",
		},
		{
			name:             "string value",
			propertyType:     "string",
			propertyInternal: "some text",
			propertyRaw:      "some text",
			errorMessage:     "",
		},
		{
			name:             "array value",
			propertyType:     "array",
			propertyInternal: "[1, 2, 3]",
			propertyRaw:      nil,
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
			CheckDiagnostics(t, diags, tc.errorMessage)
			CheckEqual(t, raw.Default, tc.propertyRaw)
		})
	}
}

func TestPropertyModelDefaultFromAPI(t *testing.T) {
	cases := []struct {
		name             string
		propertyType     string
		propertyInternal string
		propertyRaw      any
		errorMessage     string
	}{
		{
			name:             "boolean value",
			propertyType:     "boolean",
			propertyInternal: "false",
			propertyRaw:      false,
			errorMessage:     "",
		},
		/*{
			name: "boolean value",
			propertyType: "boolean",
			propertyInternal: "true",
			propertyRaw: true,
			errorMessage: "",
		},
		{
			name: "wrong boolean value",
			propertyType: "boolean",
			propertyInternal: "not really a boolean",
			propertyRaw: nil,
			errorMessage: "invalid syntax",
		},
		{
			name: "integer value",
			propertyType: "integer",
			propertyInternal: "42",
			propertyRaw: int64(42),
			errorMessage: "",
		},
		{
			name: "wrong integer value",
			propertyType: "integer",
			propertyInternal: "1.2",
			propertyRaw: nil,
			errorMessage: "invalid syntax",
		},
		{
			name: "number value (integer)",
			propertyType: "number",
			propertyInternal: "-100",
			propertyRaw: int64(-100),
			errorMessage: "",
		},
		{
			name: "number value (float)",
			propertyType: "number",
			propertyInternal: "3.141592",
			propertyRaw: 3.141592,
			errorMessage: "",
		},
		{
			name: "string value",
			propertyType: "string",
			propertyInternal: "some text",
			propertyRaw: "some text",
			errorMessage: "",
		},
		{
			name: "array value",
			propertyType: "array",
			propertyInternal: "[1, 2, 3]",
			propertyRaw: nil,
			errorMessage: "type array is not yet implemented",
		},*/
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			property := PropertyModel{}
			diags := property.FromAPI(context.Background(), PropertyAPIModel{
				Title:   "P",
				Type:    tc.propertyType,
				Default: tc.propertyRaw,
			})
			CheckDiagnostics(t, diags, tc.errorMessage)
			CheckEqual(t, property.Default.ValueString(), tc.propertyInternal)
		})
	}
}

// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestPropertyModel_Default_FromAPI(t *testing.T) {
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
			name:             "integer value (nil)",
			propertyType:     "integer",
			propertyRaw:      nil,
			propertyInternal: types.StringNull(),
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
			t.Parallel()
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

func TestPropertyModel_Default_ToAPI(t *testing.T) {
	cases := []struct {
		name             string
		propertyType     string
		propertyInternal types.String
		propertyRaw      any
		propertyJson     string
		warningMessage   string
		errorMessage     string
	}{
		{
			name:             "boolean value (false)",
			propertyType:     "boolean",
			propertyInternal: types.StringValue("false"),
			propertyRaw:      false,
			propertyJson:     "\"default\":false,",
		},
		{
			name:             "boolean value (true)",
			propertyType:     "boolean",
			propertyInternal: types.StringValue("true"),
			propertyRaw:      true,
			propertyJson:     "\"default\":true,",
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
			propertyJson:     "",
		},
		{
			name:             "integer value (integer)",
			propertyType:     "integer",
			propertyInternal: types.StringValue("42"),
			propertyRaw:      int64(42),
			propertyJson:     "\"default\":42,",
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
			propertyJson:     "\"default\":-100,",
		},
		{
			name:             "number value (float)",
			propertyType:     "number",
			propertyInternal: types.StringValue("3.141592"),
			propertyRaw:      3.141592,
			propertyJson:     "\"default\":3.141592,",
		},
		{
			name:             "string value (string)",
			propertyType:     "string",
			propertyInternal: types.StringValue("some text"),
			propertyRaw:      "some text",
			propertyJson:     "\"default\":\"some text\",",
		},
		{
			name:             "string value (string empty)",
			propertyType:     "string",
			propertyInternal: types.StringNull(),
			propertyRaw:      nil,
			propertyJson:     "",
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
			t.Parallel()
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

			if !diags.HasError() {
				rawJson, err := json.Marshal(raw)
				CheckEqual(t, err, nil)
				CheckEqual(
					t,
					string(rawJson),
					fmt.Sprintf("{"+
						"\"title\":\"P\","+
						"\"description\":\"\","+
						"\"type\":\"%s\","+
						"%s"+
						"\"encrypted\":false,"+
						"\"ReadOnly\":false,"+
						"\"recreateOnUpdate\":false,"+
						"\"pattern\":\"\""+
						"}",
						tc.propertyType,
						tc.propertyJson))
			}
		})
	}
}

func TestPropertyModel_OneOf_FromAPI(t *testing.T) {
	property := PropertyModel{}
	diags := property.FromAPI(context.Background(), "p", PropertyAPIModel{
		Title: "P",
		Type:  "string",
		OneOf: []PropertyOneOfAPIModel{
			{
				Const: "a",
				Title: "A",
			},
			{
				Const:     "b",
				Title:     "B",
				Encrypted: true,
			},
		},
	})
	CheckDiagnostics(t, diags, "", "")
	CheckDeepEqual(t, property.OneOf[0], PropertyOneOfModel{
		Const:     types.StringValue("a"),
		Title:     types.StringValue("A"),
		Encrypted: types.BoolValue(false),
	})
	CheckDeepEqual(t, property.OneOf[1], PropertyOneOfModel{
		Const:     types.StringValue("b"),
		Title:     types.StringValue("B"),
		Encrypted: types.BoolValue(true),
	})
}

func TestPropertyModel_OneOf_FromAPI_nil(t *testing.T) {
	property := PropertyModel{}
	diags := property.FromAPI(context.Background(), "p", PropertyAPIModel{
		Title: "P",
		Type:  "string",
	})
	CheckDiagnostics(t, diags, "", "")
	CheckEqual(t, len(property.OneOf), 0)
}

func TestPropertyModel_OneOf_ToAPI(t *testing.T) {
	property := PropertyModel{
		Name:  types.StringValue("p"),
		Title: types.StringValue("P"),
		Type:  types.StringValue("string"),
		OneOf: []PropertyOneOfModel{
			{
				Const: types.StringValue("a"),
				Title: types.StringValue("A"),
			},
			{
				Const:     types.StringValue("b"),
				Title:     types.StringValue("B"),
				Encrypted: types.BoolValue(true),
			},
		},
	}
	name, raw, diags := property.ToAPI(context.Background())
	CheckDiagnostics(t, diags, "", "")
	CheckDeepEqual(t, raw.OneOf[0], PropertyOneOfAPIModel{Const: "a", Title: "A"})
	CheckDeepEqual(t, raw.OneOf[1], PropertyOneOfAPIModel{
		Const:     "b",
		Title:     "B",
		Encrypted: true,
	})
	CheckEqual(t, name, "p")
}

func TestPropertyModel_OneOf_ToAPI_nil(t *testing.T) {
	property := PropertyModel{
		Name:  types.StringValue("p"),
		Title: types.StringValue("P"),
		Type:  types.StringValue("string"),
	}
	_, raw, diags := property.ToAPI(context.Background())
	CheckDiagnostics(t, diags, "", "")
	CheckEqual(t, len(raw.OneOf), 0)
}

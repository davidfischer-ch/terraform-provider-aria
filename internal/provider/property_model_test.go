// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestPropertyModel_Default_FromAPI(t *testing.T) {
	cases := []struct {
		name             string
		propertyType     string
		propertyRaw      any
		propertyInternal jsontypes.Normalized
		warningMessage   string
		errorMessage     string
	}{
		{
			name:             "boolean value (false)",
			propertyType:     "boolean",
			propertyRaw:      false,
			propertyInternal: jsontypes.NewNormalizedValue("false"),
		},
		{
			name:             "boolean value (true)",
			propertyType:     "boolean",
			propertyRaw:      true,
			propertyInternal: jsontypes.NewNormalizedValue("true"),
		},
		{
			name:             "boolean value (string)",
			propertyType:     "boolean",
			propertyRaw:      "not really a boolean",
			propertyInternal: jsontypes.NewNormalizedValue("\"not really a boolean\""),
		},
		{
			name:             "integer value (integer)",
			propertyType:     "integer",
			propertyRaw:      int64(42),
			propertyInternal: jsontypes.NewNormalizedValue("42"),
		},
		{
			name:             "integer value (float round)",
			propertyType:     "integer",
			propertyRaw:      float64(99),
			propertyInternal: jsontypes.NewNormalizedValue("99"),
		},
		{
			name:             "integer value (float)",
			propertyType:     "integer",
			propertyRaw:      float64(1.2),
			propertyInternal: jsontypes.NewNormalizedValue("1.2"),
		},
		{
			name:             "integer value (nil)",
			propertyType:     "integer",
			propertyRaw:      nil,
			propertyInternal: jsontypes.NewNormalizedNull(),
		},
		{
			name:             "number value (integer)",
			propertyType:     "number",
			propertyRaw:      int64(-100),
			propertyInternal: jsontypes.NewNormalizedValue("-100"),
		},
		{
			name:             "number value (float)",
			propertyType:     "number",
			propertyRaw:      float64(3.141592),
			propertyInternal: jsontypes.NewNormalizedValue("3.141592"),
		},
		{
			name:             "number value (nil)",
			propertyType:     "number",
			propertyRaw:      nil,
			propertyInternal: jsontypes.NewNormalizedNull(),
		},
		{
			name:             "string value (string)",
			propertyType:     "string",
			propertyRaw:      "some text",
			propertyInternal: jsontypes.NewNormalizedValue("\"some text\""),
		},
		{
			name:             "array value (nil)",
			propertyType:     "array",
			propertyRaw:      nil,
			propertyInternal: jsontypes.NewNormalizedNull(),
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

func TestPropertyModel_Default_ToAPI(t *testing.T) {
	cases := []struct {
		name             string
		propertyType     string
		propertyInternal jsontypes.Normalized
		propertyRaw      any
		propertyJson     string
		warningMessage   string
		errorMessage     string
	}{
		{
			name:             "boolean value (false)",
			propertyType:     "boolean",
			propertyInternal: jsontypes.NewNormalizedValue("false"),
			propertyRaw:      false,
			propertyJson:     "\"default\":false,",
		},
		{
			name:             "boolean value (true)",
			propertyType:     "boolean",
			propertyInternal: jsontypes.NewNormalizedValue("true"),
			propertyRaw:      true,
			propertyJson:     "\"default\":true,",
		},
		{
			name:             "boolean value (string)",
			propertyType:     "boolean",
			propertyInternal: jsontypes.NewNormalizedValue("\"not really a boolean\""),
			propertyRaw:      "not really a boolean",
			propertyJson:     "\"default\":\"not really a boolean\",",
		},
		{
			name:             "boolean value (nil)",
			propertyType:     "boolean",
			propertyInternal: jsontypes.NewNormalizedNull(),
			propertyRaw:      nil,
			propertyJson:     "",
		},
		{
			name:             "integer value (integer)",
			propertyType:     "integer",
			propertyInternal: jsontypes.NewNormalizedValue("42"),
			propertyRaw:      float64(42), // Always float, don't mind it
			propertyJson:     "\"default\":42,",
		},
		{
			name:             "integer value (float)",
			propertyType:     "integer",
			propertyInternal: jsontypes.NewNormalizedValue("1.2"),
			propertyRaw:      float64(1.2),
			propertyJson:     "\"default\":1.2,",
		},
		{
			name:             "number value (integer)",
			propertyType:     "number",
			propertyInternal: jsontypes.NewNormalizedValue("-100"),
			propertyRaw:      float64(-100), // Always float, don't mind it
			propertyJson:     "\"default\":-100,",
		},
		{
			name:             "number value (float)",
			propertyType:     "number",
			propertyInternal: jsontypes.NewNormalizedValue("3.141592"),
			propertyRaw:      3.141592,
			propertyJson:     "\"default\":3.141592,",
		},
		{
			name:             "string value (string)",
			propertyType:     "string",
			propertyInternal: jsontypes.NewNormalizedValue("\"some text\""),
			propertyRaw:      "some text",
			propertyJson:     "\"default\":\"some text\",",
		},
		{
			name:             "string value (string empty)",
			propertyType:     "string",
			propertyInternal: jsontypes.NewNormalizedNull(),
			propertyRaw:      nil,
			propertyJson:     "",
		},
		{
			name:             "array value (string empty)",
			propertyType:     "array",
			propertyInternal: jsontypes.NewNormalizedNull(),
			propertyRaw:      nil,
			propertyJson:     "",
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
			CheckDeepEqual(t, raw.Default, tc.propertyRaw)
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
						"\"readOnly\":false,"+
						"\"recreateOnUpdate\":false"+
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

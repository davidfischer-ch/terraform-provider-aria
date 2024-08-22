// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Model interface {
	String() string

	LockKey() string

	CreatePath() string
	ReadPath() string
	UpdatePath() string
	DeletePath() string
}

type APIModel interface{}

func StringOrNullValue(value string) types.String {
	// Replace empty value by nil
	if len(value) == 0 {
		return types.StringNull()
	}
	return types.StringValue(value)
}

// Convert raw value to JSON encoded attribute.
func JSONNormalizedFromAny(name string, value any) (jsontypes.Normalized, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	if value == nil {
		return jsontypes.NewNormalizedNull(), diags
	}

	valueJSON, err := json.Marshal(value)
	if err != nil {
		diags.AddError(
			"Client error",
			fmt.Sprintf(
				"Unable to JSON encode %s default \"%s\", got error: %s",
				name, value, err))
		return jsontypes.NewNormalizedNull(), diags
	}

	return jsontypes.NewNormalizedValue(string(valueJSON)), diags
}

// Convert JSON encoded attribute to raw value.
func JSONNormalizedToAny(attribute jsontypes.Normalized) (any, diag.Diagnostics) {
	if attribute.IsNull() {
		return nil, diag.Diagnostics{}
	}

	var value any
	diags := attribute.Unmarshal(&value)
	return value, diags
}

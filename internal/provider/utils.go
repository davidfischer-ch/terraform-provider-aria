// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TODO How to deduplicate code without introducing more loc?
// https://www.golinuxcloud.com/golang-function-accept-two-types/

func GetDataSourceClient(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) *resty.Client {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return nil
	}

	client, ok := req.ProviderData.(*resty.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *resty.Client, got: %T. Please report this issue to the "+
				"provider developers.", req.ProviderData),
		)
		return nil
	}

	return client
}

func GetResourceClient(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) *resty.Client {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return nil
	}

	client, ok := req.ProviderData.(*resty.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *resty.Client, got: %T. Please report this issue to the "+
				"provider developers.", req.ProviderData),
		)
		return nil
	}

	return client
}

func SkipEmpty(value []string) []string {
	return slices.DeleteFunc(value, func(e string) bool { return e == "" })
}

func CleanString(value string) string {
	return strings.Replace(value, "\r", "", -1)
}

func StringOrNullValue(value string) types.String {
	// Replace empty value by nil
	if len(value) == 0 {
		return types.StringNull()
	}
	return types.StringValue(value)
}

// Decode then encode to JSON a given JSON encoded value.
// Default to given value in case of error. You must check diagnostics.
// Use case: Stabilize Aria API response data to stabilize its content.
//
// To prevent this:
//
// Error: Provider produced inconsistent result after apply
//
// When applying changes to aria_resource_action.test, provider
// "provider[\"registry.terraform.io/hashicorp/aria\"]" produced an unexpected
// new value: .form_definition.form: was
// cty.StringVal("{\"layout\":{\"pages\":[{\"id\":\"page_1\",\"sections\":[],\"title\":\"Page
// Numéro 1\"}]},\"schema\":{}}"), but now
// cty.StringVal("{\"layout\":{\"pages\":[{\"id\":\"page_1\",\"title\":\"Page
// Numéro 1\",\"sections\":[]}]},\"schema\":{}}").
//
// This is a bug in the provider, which should be reported in the provider's own issue tracker.
func JSONRencode(value string, title string) (string, diag.Diagnostics) {

	diags := diag.Diagnostics{}

	// Value JSON Encoded -> data
	var data interface{}
	err := json.Unmarshal([]byte(value), &data)
	if err != nil {
		diags.AddError(
			"Client error",
			fmt.Sprintf("Unable to JSON decode %s, got error: %s", title, err))
		return value, diags
	}

	// data -> JSON Encoded value
	valueBack, err := json.Marshal(data)
	if err != nil {
		diags.AddError(
			"Client error",
			fmt.Sprintf("Unable to JSON encode %s, got error: %s", title, err))
		return value, diags
	}
	return string(valueBack), diags
}

// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
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
			fmt.Sprintf("Expected *resty.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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
			fmt.Sprintf("Expected *resty.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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

// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// TODO How to deduplicate code without introducing more loc?
// https://www.golinuxcloud.com/golang-function-accept-two-types/

func GetDataSourceClient(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) *AriaClient {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return nil
	}

	client, ok := req.ProviderData.(*AriaClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *AriaClient, got: %T. Please report this issue to the "+
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
) *AriaClient {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return nil
	}

	client, ok := req.ProviderData.(*AriaClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *AriaClient, got: %T. Please report this issue to the "+
				"provider developers.", req.ProviderData),
		)
		return nil
	}

	return client
}

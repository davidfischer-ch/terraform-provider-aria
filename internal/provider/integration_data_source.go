// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &IntegrationDataSource{}

func NewIntegrationDataSource() datasource.DataSource {
	return &IntegrationDataSource{}
}

// IntegrationDataSource defines the data source implementation.
type IntegrationDataSource struct {
	client *AriaClient
}

func (self *IntegrationDataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_integration"
}

func (self *IntegrationDataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = IntegrationDataSourceSchema()
}

func (self *IntegrationDataSource) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	self.client = GetDataSourceClient(ctx, req, resp)
}

func (self *IntegrationDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	// Read Terraform configuration data into the model
	var integration IntegrationDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &integration)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var responseFromAPI IntegrationResponseAPIodel
	path := integration.ReadPath()
	response, err := self.client.R(path).
		SetQueryParam("size", "1").
		SetQueryParam("page", "0").
		SetQueryParam("sort", "name,asc").
		SetResult(&responseFromAPI).
		Get(path)
	err = self.client.HandleAPIResponse(response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to get %s, got error: %s", integration.String(), err))
		return
	}

	if len(responseFromAPI.Content) == 0 {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to get %s, no content found.", integration.String()))
		return
	}

	for _, contentsRaw := range responseFromAPI.Content {
		integration.FromAPI(contentsRaw.Integration)
		break // Make sure we don't set it multiple times for nothing
	}

	// Save updated integration into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &integration)...)
}

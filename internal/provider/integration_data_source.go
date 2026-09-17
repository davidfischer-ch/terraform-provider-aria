// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

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

	description := integration.String()

	entries, someDiags := ReadIntegrations(self.client, description)
	resp.Diagnostics.Append(someDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	match, someDiags := SelectIntegration(
		entries,
		integration.IntegrationType(),
		integration.Name.ValueString(),
		description,
	)
	resp.Diagnostics.Append(someDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	integration.FromAPI(match)

	// Save updated integration into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &integration)...)
}

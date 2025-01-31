// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &OrchestratorConfigurationDataSource{}

func NewOrchestratorConfigurationDataSource() datasource.DataSource {
	return &OrchestratorConfigurationDataSource{}
}

// OrchestratorConfigurationDataSource defines the data source implementation.
type OrchestratorConfigurationDataSource struct {
	client *AriaClient
}

func (self *OrchestratorConfigurationDataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_orchestrator_configuration"
}

func (self *OrchestratorConfigurationDataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = OrchestratorConfigurationDataSourceSchema()
}

func (self *OrchestratorConfigurationDataSource) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	self.client = GetDataSourceClient(ctx, req, resp)
}

func (self *OrchestratorConfigurationDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	// Read Terraform configuration data into the model
	var configuration OrchestratorConfigurationDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &configuration)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var configurationRaw OrchestratorConfigurationAPIModel
	response, err := self.client.Client.R().
		// TODO SetQueryParam("apiVersion", ORCHESTRATOR_API_VERSION).
		SetResult(&configurationRaw).
		Get(configuration.ReadPath())

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read %s, got error: %s", configuration.String(), err))
		return
	}

	// Save updated configuration into Terraform state
	resp.Diagnostics.Append(configuration.FromAPI(ctx, configurationRaw, response)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &configuration)...)
}

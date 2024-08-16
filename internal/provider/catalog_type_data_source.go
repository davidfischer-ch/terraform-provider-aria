// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CatalogTypeDataSource{}

func NewCatalogTypeDataSource() datasource.DataSource {
	return &CatalogTypeDataSource{}
}

// CatalogTypeDataSource defines the data source implementation.
type CatalogTypeDataSource struct {
	client *AriaClient
}

func (self *CatalogTypeDataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_catalog_type"
}

func (self *CatalogTypeDataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = CatalogTypeDataSourceSchema()
}

func (self *CatalogTypeDataSource) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	self.client = GetDataSourceClient(ctx, req, resp)
}

func (self *CatalogTypeDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	// Read Terraform configuration data into the model
	var catalogType CatalogTypeModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &catalogType)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var catalogTypeRaw CatalogTypeAPIModel
	catalogTypeId := catalogType.Id.ValueString()
	response, err := self.client.Client.R().
		SetQueryParam("apiVersion", CATALOG_API_VERSION).
		SetResult(&catalogTypeRaw).
		Get("/catalog/api/types/" + catalogTypeId)

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read catalog type %s, got error: %s", catalogTypeId, err))
		return
	}

	// Save updated catalog type into Terraform state
	resp.Diagnostics.Append(catalogType.FromAPI(ctx, catalogTypeRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &catalogType)...)
}

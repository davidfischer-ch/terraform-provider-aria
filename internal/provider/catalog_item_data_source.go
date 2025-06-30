// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CatalogItemDataSource{}

func NewCatalogItemDataSource() datasource.DataSource {
	return &CatalogItemDataSource{}
}

// CatalogItemDataSource defines the data source implementation.
type CatalogItemDataSource struct {
	client *AriaClient
}

func (self *CatalogItemDataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_catalog_item"
}

func (self *CatalogItemDataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = CatalogItemDataSourceSchema()
}

func (self *CatalogItemDataSource) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	self.client = GetDataSourceClient(ctx, req, resp)
}

func (self CatalogItemDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("external_id"),
		),
	}
}

func (self *CatalogItemDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	// Read Terraform configuration data into the model
	var item CatalogItemModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &item)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var itemFromAPI CatalogItemAPIModel

	if len(item.Id.ValueString()) > 0 {
		// Retrieve details from the item's API endpoint
		path := item.ReadPath()
		response, err := self.client.R(path).SetResult(&itemFromAPI).Get(path)
		err = self.client.HandleAPIResponse(response, err, []int{200})
		if err != nil {
			resp.Diagnostics.AddError(
				"Client error",
				fmt.Sprintf("Unable to read %s, got error: %s", item.String(), err))
			return
		}
	} else {
		// Retrieve details from the items API endpoint

		var listFromAPI CatalogItemLstAPIModel
		listPath := item.ListPath()
		query := self.client.R(listPath)

		// Setup search query
		name := item.Name.ValueString()
		if len(name) > 0 {
			query = query.SetQueryParam("search", name)
		}
		sourceId := item.SourceId.ValueString()
		if len(sourceId) > 0 {
			query = query.SetQueryParam("sourceIds", sourceId)
		}
		typeId := item.TypeId.ValueString()
		if len(typeId) > 0 {
			query = query.SetQueryParam("types", typeId)
		}

		// Execute search query
		response, err := query.
			SetQueryParam("size", "1000"). // Don't want to play with pagination
			SetResult(&listFromAPI).Get(listPath)
		err = self.client.HandleAPIResponse(response, err, []int{200})
		if err != nil {
			resp.Diagnostics.AddError(
				"Client error",
				fmt.Sprintf("Unable to list items to get %s, got error: %s", item.String(), err))
			return
		}

		externalId := item.ExternalId.ValueString()
		found := false

		// Lookup for the item matching given external ID
		for _, itemRaw := range listFromAPI.Content {
			// Retrieve details from the item's API endpoint
			path := CatalogItemModel{Id: types.StringValue(itemRaw.Id)}.ReadPath()
			response, err = self.client.R(path).SetResult(&itemFromAPI).Get(path)
			err = self.client.HandleAPIResponse(response, err, []int{200})
			if err != nil {
				resp.Diagnostics.AddError(
					"Client error",
					fmt.Sprintf("Unable to read candidate %s, got error: %s", item.String(), err))
				return
			}
			// Found a match!
			if itemFromAPI.ExternalId == externalId {
				found = true
				break
			}
		}

		// Fail miserably
		if !found {
			resp.Diagnostics.AddError(
				"Client error",
				fmt.Sprintf(
					"Unable to find %s matching external ID & attributes, found %d candidate items",
					item.String(), len(listFromAPI.Content)))
			return
		}
	}

	// Save updated catalog type into Terraform state
	resp.Diagnostics.Append(item.FromAPI(itemFromAPI)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &item)...)
}

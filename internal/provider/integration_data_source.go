// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"slices"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

	entries, someDiags := self.ReadIntegrations(ctx, integration)
	resp.Diagnostics.Append(someDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The listing mixes every kind of integration, hence the filtering on the integration type
	// behind the requested catalog source type.
	integrationType := integration.IntegrationType()
	candidates := []IntegrationAPIModel{}
	for _, entry := range entries {
		if entry.IntegrationType == integrationType {
			candidates = append(candidates, entry.ToIntegrationAPI())
		}
	}

	// Filtering by name is what makes the lookup deterministic when there are several candidates.
	name := integration.Name.ValueString()
	matches := candidates
	if len(name) > 0 {
		matches = slices.DeleteFunc(slices.Clone(candidates), func(c IntegrationAPIModel) bool {
			return c.Name != name
		})
	}

	if len(matches) == 0 {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf(
				"Unable to get %s, no integration of type %s matches, %d found: %s.",
				integration.String(),
				integrationType,
				len(candidates),
				IntegrationCandidates(candidates)))
		return
	}

	if len(matches) > 1 {
		resp.Diagnostics.AddError(
			"Configuration error",
			fmt.Sprintf(
				"Unable to get %s, %d integrations of type %s match, set name to one of: %s.",
				integration.String(),
				len(matches),
				integrationType,
				IntegrationCandidates(matches)))
		return
	}

	integration.FromAPI(matches[0])

	// Save updated integration into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &integration)...)
}

// -------------------------------------------------------------------------------------------------

// Return every integration of the organization, walking the pages of the listing.
//
// The offset is the number of entries received so far rather than a page number, which keeps the
// walk correct whatever page size the API decides to serve.
func (self *IntegrationDataSource) ReadIntegrations(
	ctx context.Context,
	integration IntegrationDataSourceModel,
) ([]IntegrationEntryAPIModel, diag.Diagnostics) {

	diags := diag.Diagnostics{}
	entries := []IntegrationEntryAPIModel{}
	identifiers := map[string]bool{}
	received := 0

	path := integration.ReadPath()
	for {
		var responseFromAPI IntegrationsResponseAPIModel
		response, err := self.client.R(path).
			SetQueryParam("$top", strconv.Itoa(INTEGRATIONS_PAGE_SIZE)).
			SetQueryParam("$skip", strconv.Itoa(received)).
			SetQueryParam("$orderby", "name asc"). // An offset is only meaningful over a stable order
			SetResult(&responseFromAPI).
			Get(path)
		err = self.client.HandleAPIResponse(response, err, []int{200})
		if err != nil {
			diags.AddError(
				"Client error",
				fmt.Sprintf("Unable to get %s, got error: %s", integration.String(), err))
			return entries, diags
		}

		tflog.Debug(
			ctx,
			fmt.Sprintf(
				"Read %d integrations from offset %d of %d in total",
				len(responseFromAPI.Content), received, responseFromAPI.TotalElements))

		received += len(responseFromAPI.Content)

		// An API serving the same page over and over would otherwise loop forever, and counting
		// an integration twice would make it ambiguous with itself.
		fresh := 0
		for _, entry := range responseFromAPI.Content {
			if !identifiers[entry.Id] {
				identifiers[entry.Id] = true
				entries = append(entries, entry)
				fresh++
			}
		}

		if len(entries) >= responseFromAPI.TotalElements || len(responseFromAPI.Content) == 0 {
			return entries, diags
		}

		if fresh == 0 {
			diags.AddError(
				"Client error",
				fmt.Sprintf(
					"Unable to get %s, the API returned %d integrations of %d and then stopped "+
						"returning new ones, is it honouring $top and $skip?",
					integration.String(), len(entries), responseFromAPI.TotalElements))
			return entries, diags
		}
	}
}

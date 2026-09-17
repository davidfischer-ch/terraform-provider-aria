// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// Looking an integration up is shared by the aria_integration data source and by the provider
// resolving the Orchestrator to address from the name of its integration.

// Return every integration of the organization, walking the pages of the listing. Description
// names the lookup in the diagnostics.
//
// The offset is the number of entries received so far rather than a page number, which keeps the
// walk correct whatever page size the API decides to serve.
func ReadIntegrations(
	client *AriaClient,
	description string,
) ([]IntegrationEntryAPIModel, diag.Diagnostics) {

	diags := diag.Diagnostics{}
	entries := []IntegrationEntryAPIModel{}
	identifiers := map[string]bool{}
	received := 0

	path := INTEGRATIONS_PATH
	for {
		var responseFromAPI IntegrationsResponseAPIModel
		response, err := client.R(path).
			SetQueryParam("$top", strconv.Itoa(INTEGRATIONS_PAGE_SIZE)).
			SetQueryParam("$skip", strconv.Itoa(received)).
			SetQueryParam("$orderby", "name asc"). // An offset is only meaningful over a stable order
			SetResult(&responseFromAPI).
			Get(path)
		err = client.HandleAPIResponse(response, err, []int{200})
		if err != nil {
			diags.AddError(
				"Client error",
				fmt.Sprintf("Unable to get %s, got error: %s", description, err))
			return entries, diags
		}

		client.Debug(
			"Read %d integrations from offset %d of %d in total",
			len(responseFromAPI.Content), received, responseFromAPI.TotalElements)

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
					description, len(entries), responseFromAPI.TotalElements))
			return entries, diags
		}
	}
}

// Return the sole integration of integrationType named name. An empty name matches when the
// platform exposes a single integration of that type, naming one is what makes the lookup
// deterministic otherwise. Description names the lookup in the diagnostics.
func SelectIntegration(
	entries []IntegrationEntryAPIModel,
	integrationType string,
	name string,
	description string,
) (IntegrationAPIModel, diag.Diagnostics) {

	diags := diag.Diagnostics{}

	// The listing mixes every kind of integration, hence the filtering on the integration type.
	candidates := []IntegrationAPIModel{}
	for _, entry := range entries {
		if entry.IntegrationType == integrationType {
			candidates = append(candidates, entry.ToIntegrationAPI())
		}
	}

	matches := candidates
	if len(name) > 0 {
		matches = slices.DeleteFunc(slices.Clone(candidates), func(c IntegrationAPIModel) bool {
			return c.Name != name
		})
	}

	if len(matches) == 0 {
		diags.AddError(
			"Client error",
			fmt.Sprintf(
				"Unable to get %s, no integration of type %s matches, %d found: %s.",
				description,
				integrationType,
				len(candidates),
				IntegrationCandidates(candidates)))
		return IntegrationAPIModel{}, diags
	}

	if len(matches) > 1 {
		diags.AddError(
			"Configuration error",
			fmt.Sprintf(
				"Unable to get %s, %d integrations of type %s match, set name to one of: %s.",
				description,
				len(matches),
				integrationType,
				IntegrationCandidates(matches)))
		return IntegrationAPIModel{}, diags
	}

	return matches[0], diags
}

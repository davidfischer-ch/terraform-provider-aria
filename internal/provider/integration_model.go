// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// IntegrationModel describes the resource data model.
type IntegrationModel struct {
	Name                      types.String `tfsdk:"name"`
	EndpointConfigurationLink types.String `tfsdk:"endpoint_configuration_link"`
	EndpointURI               types.String `tfsdk:"endpoint_uri"`
}

// IntegrationDataSourceModel describes the data source data model.
type IntegrationDataSourceModel struct {
	TypeId types.String `tfsdk:"type_id"`
	IntegrationModel
}

// IntegrationAPIModel describes the resource API model.
type IntegrationAPIModel struct {
	Name                      string `json:"name"`
	EndpointConfigurationLink string `json:"endpointConfigurationLink"`
	EndpointURI               string `json:"endpointUri"`
}

// IntegrationsResponseAPIModel describes the response of the integrations API endpoint.
type IntegrationsResponseAPIModel struct {
	Content       []IntegrationEntryAPIModel `json:"content"`
	TotalElements int                        `json:"totalElements"`
}

// IntegrationEntryAPIModel describes an integration as listed by the integrations API endpoint.
type IntegrationEntryAPIModel struct {
	Id               string            `json:"id"`
	Name             string            `json:"name"`
	IntegrationType  string            `json:"integrationType"`
	CustomProperties map[string]string `json:"customProperties"`
}

// Return the integration in the shape the catalog API uses. The endpoint configuration link points
// at the endpoint document, which shares the identifier of the integration wrapping it.
func (self IntegrationEntryAPIModel) ToIntegrationAPI() IntegrationAPIModel {
	return IntegrationAPIModel{
		Name:                      self.Name,
		EndpointConfigurationLink: "/resources/endpoints/" + self.Id,
		EndpointURI:               self.CustomProperties["hostName"],
	}
}

// Return the candidates sorted and joined for use in a diagnostic message. Duplicates are kept,
// how many times a name occurs is part of what the message has to report.
func IntegrationCandidates(candidates []IntegrationAPIModel) string {
	entries := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		entries = append(entries, fmt.Sprintf("%q (%s)", candidate.Name, candidate.EndpointURI))
	}
	slices.Sort(entries)
	return strings.Join(entries, ", ")
}

func (self *IntegrationModel) String() string {
	return fmt.Sprintf(
		"Integration %s (%s)",
		self.Name.ValueString(),
		self.EndpointURI.ValueString())
}

func (self *IntegrationModel) FromAPI(raw IntegrationAPIModel) {
	self.Name = types.StringValue(raw.Name)
	self.EndpointConfigurationLink = types.StringValue(raw.EndpointConfigurationLink)
	self.EndpointURI = types.StringValue(raw.EndpointURI)
}

func (self *IntegrationModel) ToAPI() IntegrationAPIModel {
	return IntegrationAPIModel{
		Name:                      self.Name.ValueString(),
		EndpointConfigurationLink: self.EndpointConfigurationLink.ValueString(),
		EndpointURI:               self.EndpointURI.ValueString(),
	}
}

// Return a description of what is looked up, the name included when it is set.
func (self IntegrationDataSourceModel) String() string {
	typeId := self.TypeId.ValueString()
	if name := self.Name.ValueString(); len(name) > 0 {
		return fmt.Sprintf("Integration %q of type %s", name, typeId)
	}
	return fmt.Sprintf("Integration of type %s", typeId)
}

func (self IntegrationDataSourceModel) ReadPath() string {
	return "iaas/api/integrations"
}

// Return the integration type the catalog source type identifier stands for. The listing mixes
// every kind of integration, this is what tells an Orchestrator from an extensibility endpoint.
func (self IntegrationDataSourceModel) IntegrationType() string {
	typeId := self.TypeId.ValueString()
	if typeId == "com.vmw.vro.workflow" {
		return "vro"
	}
	// Panic is intentional: this is a programming bug, not a runtime error.
	panic(fmt.Sprintf("Internal error: %s as unexpected type: %s.", self.String(), typeId))
}

// Utils -------------------------------------------------------------------------------------------

// Used to convert structure to a types.Object.
func (self IntegrationModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":                        types.StringType,
		"endpoint_configuration_link": types.StringType,
		"endpoint_uri":                types.StringType,
	}
}

// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"

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

// IntegrationResponseAPIodel describes the resource API model.
type IntegrationResponseAPIodel struct {
	Content []IntegrationResponseContentAPIModel `json:"content"`
}

// IntegrationResponseContentAPIModel describes the resource API model.
type IntegrationResponseContentAPIModel struct {
	Integration IntegrationAPIModel `json:"integration"`
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

func (self IntegrationDataSourceModel) ReadPath() string {
	var resource string
	typeId := self.TypeId.ValueString()
	if typeId == "com.vmw.vro.workflow" {
		resource = "workflows"
	} else {
		panic(fmt.Sprintf("Internal error: %s as unexpected type: %s.", self.String(), typeId))
	}
	return fmt.Sprintf("/catalog/api/types/%s/data/%s", typeId, resource)
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

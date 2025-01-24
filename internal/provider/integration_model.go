// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// IntegrationModel describes the resource data model.
type IntegrationModel struct {
	Name                      types.String `tfsdk:"name"`
	EndpointConfigurationLink types.String `tfsdk:"endpoint_configuration_link"`
	EndpointURI               types.String `tfsdk:"endpoint_uri"`
}

// IntegrationAPIModel describes the resource API model.
type IntegrationAPIModel struct {
	Name                      string `json:"name"`
	EndpointConfigurationLink string `json:"endpointConfigurationLink"`
	EndpointURI               string `json:"endpointUri"`
}

func (self *IntegrationModel) String() string {
	return fmt.Sprintf(
		"Integration %s (%s)",
		self.Name.ValueString(),
		self.EndpointURI.ValueString())
}

func (self *IntegrationModel) FromAPI(
	ctx context.Context,
	raw IntegrationAPIModel,
) diag.Diagnostics {
	self.Name = types.StringValue(raw.Name)
	self.EndpointConfigurationLink = types.StringValue(raw.EndpointConfigurationLink)
	self.EndpointURI = types.StringValue(raw.EndpointURI)
	return diag.Diagnostics{}
}

func (self *IntegrationModel) ToAPI(
	ctx context.Context,
) (IntegrationAPIModel, diag.Diagnostics) {
	return IntegrationAPIModel{
		Name:                      self.Name.ValueString(),
		EndpointConfigurationLink: self.EndpointConfigurationLink.ValueString(),
		EndpointURI:               self.EndpointURI.ValueString(),
	}, diag.Diagnostics{}
}

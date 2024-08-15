// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

// CloudTemplateV1ContentAPIModel describes the resource API model.
type CloudTemplateV1ContentAPIModel struct {
	/*FormatVersion string                         `json:"formatVersion"`*/
	Inputs    UnorderedPropertiesAPIModel    `yaml:"inputs"`
	Resources CloudTemplateResourcesAPIModel `yaml:"resources"`
}

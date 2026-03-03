// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import "github.com/hashicorp/terraform-plugin-framework/resource"

func NewOrchestratorEnvironmentRepositoryResource() resource.Resource {
	return &SimpleGenericResource[
		OrchestratorEnvironmentRepositoryModel,
		*OrchestratorEnvironmentRepositoryModel,
		OrchestratorEnvironmentRepositoryAPIModel,
	]{
		config: GenericResourceConfig{
			TypeName:    "_orchestrator_environment_repository",
			SchemaFunc:  OrchestratorEnvironmentRepositorySchema,
			UpdateCodes: []int{202},
			ImportStateSetAttributes: map[string]string{
				"system_credentials": "",
			},
		},
	}
}

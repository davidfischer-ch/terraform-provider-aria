// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import "github.com/hashicorp/terraform-plugin-framework/resource"

func NewOrchestratorTaskResource() resource.Resource {
	return &GenericResource[OrchestratorTaskModel, *OrchestratorTaskModel, OrchestratorTaskAPIModel]{
		config: GenericResourceConfig{
			TypeName:     "_orchestrator_task",
			SchemaFunc:   OrchestratorTaskSchema,
			CreateCodes:  []int{202},
			UpdateMethod: "POST",
		},
	}
}

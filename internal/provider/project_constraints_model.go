// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	//"github.com/hashicorp/terraform-plugin-framework/types"
)

// ProjectConstraintsModel describes the resource data model.
type ProjectConstraintsModel struct {
}

// ProjectConstraintsAPIModel describes the resource API model.
type ProjectConstraintsAPIModel struct {
}

func (self ProjectConstraintsModel) String() string {
	return "Project Constraints"
}

func (self *ProjectConstraintsModel) FromAPI(
	ctx context.Context,
	raw ProjectConstraintsAPIModel,
) diag.Diagnostics {
	return diag.Diagnostics{}
}

func (self ProjectConstraintsModel) ToAPI(
	ctx context.Context,
) (ProjectConstraintsAPIModel, diag.Diagnostics) {
	return ProjectConstraintsAPIModel{}, diag.Diagnostics{}
}

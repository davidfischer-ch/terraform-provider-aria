// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

// ProjectConstraintsModel describes the resource data model.
type ProjectConstraintsModel struct {
}

// ProjectConstraintsAPIModel describes the resource API model.
type ProjectConstraintsAPIModel struct {
}

func (self ProjectConstraintsModel) String() string {
	return "Project Constraints"
}

func (self *ProjectConstraintsModel) FromAPI(raw ProjectConstraintsAPIModel) {
}

func (self ProjectConstraintsModel) ToAPI() ProjectConstraintsAPIModel {
	return ProjectConstraintsAPIModel{}
}

// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// IconModel describes the resource data model.
type IconModel struct {
	Id      types.String `tfsdk:"id"`
	Content types.String `tfsdk:"content"`
}

func (self IconModel) String() string {
	return fmt.Sprintf("Icon %s", self.Id.ValueString())
}

func (self IconModel) CreatePath() string {
	return "icon/api/icons"
}

func (self IconModel) ReadPath() string {
	return "icon/api/icons/" + self.Id.ValueString()
}

func (self IconModel) UpdatePath() string {
	return self.ReadPath() // Even if not possible ...
}

func (self IconModel) DeletePath() string {
	return self.ReadPath()
}

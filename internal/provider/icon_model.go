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

func (self *IconModel) String() string {
	return fmt.Sprintf("Icon %s", self.Id.ValueString())
}

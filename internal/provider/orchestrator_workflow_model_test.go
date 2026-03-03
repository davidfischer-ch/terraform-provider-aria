// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestOrchestratorWorkflowModel_ReadFormPath(t *testing.T) {
	model := OrchestratorWorkflowModel{
		Id: types.StringValue("abc-123"),
	}
	CheckEqual(
		t,
		model.ReadFormPath(),
		"vco/api/forms/?conditions=workflow=abc-123&designerMod=true")
}

func TestOrchestratorWorkflowModel_ReadFormPath_SpecialChars(t *testing.T) {
	model := OrchestratorWorkflowModel{
		Id: types.StringValue("id with spaces&special=chars"),
	}
	CheckEqual(
		t,
		model.ReadFormPath(),
		"vco/api/forms/?conditions=workflow=id+with+spaces%26special%3Dchars&designerMod=true")
}

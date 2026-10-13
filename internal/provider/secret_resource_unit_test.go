// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSecretValueChecksum(t *testing.T) {
	CheckEqual(
		t,
		SecretValueChecksum("a"),
		"ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb")
	CheckEqual(
		t,
		SecretValueChecksum(""),
		"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
}

func TestSecretValueChanged(t *testing.T) {
	checksum := types.StringValue(SecretValueChecksum("a"))
	CheckEqual(t, SecretValueChanged(types.StringValue("a"), checksum), false)
	CheckEqual(t, SecretValueChanged(types.StringValue("b"), checksum), true)
	// Unknown value (computed from something not yet applied): will change
	CheckEqual(t, SecretValueChanged(types.StringUnknown(), checksum), true)
	// No checksum (import, secret modified outside of Terraform): cannot be verified
	CheckEqual(t, SecretValueChanged(types.StringValue("a"), types.StringNull()), true)
	CheckEqual(t, SecretValueChanged(types.StringValue("a"), types.StringUnknown()), true)
}

func TestSecretModifiedOutside(t *testing.T) {
	updatedAt := types.StringValue("2026-10-13T12:56:33.466Z")
	CheckEqual(t, SecretModifiedOutside(updatedAt, "2026-10-13T12:56:33.466Z"), false)
	CheckEqual(t, SecretModifiedOutside(updatedAt, "2026-10-13T12:56:33.403Z"), true)
	// Import: nothing stored to compare with
	CheckEqual(t, SecretModifiedOutside(types.StringNull(), "2026-10-13T12:56:33.403Z"), false)
}

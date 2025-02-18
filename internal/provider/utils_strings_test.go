// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"
)

func TestCleanString(t *testing.T) {
	result := CleanString("\r\nSome text from Windows.\r\nThat should be cleaned.\r\n")
	CheckEqual(t, result, "\nSome text from Windows.\nThat should be cleaned.\n")
}

func TestSkipEmpty(t *testing.T) {
	result := SkipEmpty([]string{"", "a", "", "b", "", "", "some c", " and d"})
	CheckDeepEqual(t, result, []string{"a", "b", "some c", " and d"})
}

// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"reflect"
	"testing"
)

func TestCleanString(t *testing.T) {
	expected := "\nSome text from Windows.\nThat should be cleaned.\n"
	result := CleanString("\r\nSome text from Windows.\r\nThat should be cleaned.\r\n")
	if result != expected {
		t.Errorf("Result was incorrect, got: %s, expected %s", result, expected)
	}
}

func TestSkipEmpty(t *testing.T) {
	expected := []string{"a", "b", "some c", " and d"}
	result := SkipEmpty([]string{"", "a", "", "b", "", "", "some c", " and d"})
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Result was incorrect, got: %s, expected %s", result, expected)
	}
}

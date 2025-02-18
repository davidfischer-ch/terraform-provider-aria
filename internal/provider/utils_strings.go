// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"slices"
	"strings"
)

func SkipEmpty(value []string) []string {
	return slices.DeleteFunc(value, func(e string) bool { return e == "" })
}

func CleanString(value string) string {
	return strings.Replace(value, "\r", "", -1)
}

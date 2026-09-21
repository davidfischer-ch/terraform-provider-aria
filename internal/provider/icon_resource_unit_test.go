// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"
)

func TestIconMIMEType(t *testing.T) {
	// An SVG is the reason this helper exists, content sniffing reports it as text/xml.
	CheckEqual(t, IconMIMEType("../../tests/icon.svg"), "image/svg+xml")
	CheckEqual(t, IconMIMEType("../../tests/icon.png"), "image/png")
	CheckEqual(t, IconMIMEType("ICON.SVG"), "image/svg+xml")
	CheckEqual(t, IconMIMEType("icon.unknown"), "application/octet-stream")
	CheckEqual(t, IconMIMEType("icon"), "application/octet-stream")
}

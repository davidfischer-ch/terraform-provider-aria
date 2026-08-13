// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"
)

func TestUnorderedPropertiesAPIModel(t *testing.T) {
	raw := UnorderedPropertiesAPIModel{}

	// Insert some properties then retrieve them
	raw["some"] = PropertyAPIModel{Title: "Some"}
	raw["other"] = PropertyAPIModel{Title: "Other"}
	raw["another"] = PropertyAPIModel{Title: "Another"}
	raw["latest"] = PropertyAPIModel{Title: "Latest"}
	CheckDeepEqual(t, raw["other"], PropertyAPIModel{Title: "Other"})
	CheckDeepEqual(t, raw, UnorderedPropertiesAPIModel{
		// Order is not guarantee but not relevant too...
		"another": PropertyAPIModel{Title: "Another"},
		"some":    PropertyAPIModel{Title: "Some"},
		"latest":  PropertyAPIModel{Title: "Latest"},
		"other":   PropertyAPIModel{Title: "Other"},
	})

	// Overwrite
	raw["another"] = PropertyAPIModel{Title: "Another v2"}
	CheckDeepEqual(t, raw, UnorderedPropertiesAPIModel{
		// Order is not guarantee but not relevant too...
		"another": PropertyAPIModel{Title: "Another v2"},
		"latest":  PropertyAPIModel{Title: "Latest"},
		"other":   PropertyAPIModel{Title: "Other"},
		"some":    PropertyAPIModel{Title: "Some"},
	})
}

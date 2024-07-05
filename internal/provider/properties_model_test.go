// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"encoding/json"
	"testing"
)

func TestPropertiesAPIModel(t *testing.T) {
	raw := PropertiesAPIModel{}
	raw.Init()
	CheckDeepEqual(t, raw.Get("missing"), PropertyAPIModel{})
	CheckDeepEqual(t, raw.Items(), []PropertiesAPIModelItem{})

	// Insert some properties then retrieve them
	CheckEqual(t, raw.Set("some", PropertyAPIModel{Title: "Some"}), true)
	CheckEqual(t, raw.Set("other", PropertyAPIModel{Title: "Other"}), true)
	CheckEqual(t, raw.Set("another", PropertyAPIModel{Title: "Another"}), true)
	CheckEqual(t, raw.Set("latest", PropertyAPIModel{Title: "Latest"}), true)
	CheckDeepEqual(t, raw.Get("other"), PropertyAPIModel{Title: "Other"})
	CheckDeepEqual(t, raw.Items(), []PropertiesAPIModelItem{
		{Name: "some", Property: PropertyAPIModel{Title: "Some"}},
		{Name: "other", Property: PropertyAPIModel{Title: "Other"}},
		{Name: "another", Property: PropertyAPIModel{Title: "Another"}},
		{Name: "latest", Property: PropertyAPIModel{Title: "Latest"}},
	})

	// Pop
	_, exists := raw.Pop("other")
	CheckEqual(t, exists, true)

	CheckDeepEqual(t, raw.Items(), []PropertiesAPIModelItem{
		{Name: "some", Property: PropertyAPIModel{Title: "Some"}},
		{Name: "another", Property: PropertyAPIModel{Title: "Another"}},
		{Name: "latest", Property: PropertyAPIModel{Title: "Latest"}},
	})

	_, exists = raw.Pop("other")
	CheckEqual(t, exists, false)

	// Overwrite
	CheckEqual(t, raw.Set("another", PropertyAPIModel{Title: "Another v2"}), false)
	CheckDeepEqual(t, raw.Items(), []PropertiesAPIModelItem{
		{Name: "some", Property: PropertyAPIModel{Title: "Some"}},
		{Name: "another", Property: PropertyAPIModel{Title: "Another v2"}},
		{Name: "latest", Property: PropertyAPIModel{Title: "Latest"}},
	})
}

func TestPropertiesAPIModelJSON(t *testing.T) {
	raw := PropertiesAPIModel{}
	raw.Init()
	CheckEqual(t, raw.Set("some", PropertyAPIModel{Title: "Some"}), true)
	CheckEqual(t, raw.Set("other", PropertyAPIModel{Title: "Other"}), true)
	CheckEqual(t, raw.Set("another", PropertyAPIModel{Title: "Another"}), true)
	CheckEqual(t, raw.Set("yet more", PropertyAPIModel{Title: "Yet More"}), true)
	CheckEqual(t, raw.Set("latest", PropertyAPIModel{Title: "Latest"}), true)

	// Marshal and unmarshal loop will provide the same properties with the same ordering
	data, err := json.Marshal(raw)
	if err != nil {
		panic(err)
	}

	rawBis := PropertiesAPIModel{}
	CheckDeepEqual(t, rawBis.Items(), []PropertiesAPIModelItem{})

	err = json.Unmarshal(data, &rawBis)
	if err != nil {
		panic(err)
	}

	CheckDeepEqual(t, rawBis.Items(), []PropertiesAPIModelItem{
		{Name: "some", Property: PropertyAPIModel{Title: "Some"}},
		{Name: "other", Property: PropertyAPIModel{Title: "Other"}},
		{Name: "another", Property: PropertyAPIModel{Title: "Another"}},
		{Name: "yet more", Property: PropertyAPIModel{Title: "Yet More"}},
		{Name: "latest", Property: PropertyAPIModel{Title: "Latest"}},
	})
}

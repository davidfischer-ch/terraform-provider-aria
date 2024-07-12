// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type PropertiesModel []PropertyModel

// PropertiesAPIModel describes the resource API model.
// Refers
//
//	JSON and Go        https://blog.golang.org/json-and-go
//	Go-Ordered-JSON    https://github.com/virtuald/go-ordered-json
//	Python OrderedDict https://github.com/python/cpython/blob/2.7/Lib/collections.py#L38
//	port OrderedDict   https://github.com/cevaris/ordered_map
type PropertiesAPIModel struct {
	Names []string
	Data  map[string]PropertyAPIModel
}

type PropertiesAPIModelItem struct {
	Name     string
	Property PropertyAPIModel
}

func (self *PropertiesModel) FromAPI(
	ctx context.Context,
	raw PropertiesAPIModel,
) diag.Diagnostics {
	diags := diag.Diagnostics{}
	*self = PropertiesModel{}
	for _, propertyItem := range raw.Items() {
		property := PropertyModel{}
		diags.Append(property.FromAPI(ctx, propertyItem.Name, propertyItem.Property)...)
		*self = append(*self, property)
	}
	return diags
}

func (self *PropertiesModel) ToAPI(
	ctx context.Context,
) (PropertiesAPIModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	properties := PropertiesAPIModel{}
	properties.Init()
	for _, property := range *self {
		propertyName, propertyRaw, propertyDiags := property.ToAPI(ctx)
		properties.Set(propertyName, propertyRaw)
		diags.Append(propertyDiags...)
	}
	return properties, diags
}

// Reset the structure (prepare it for collecting new data).
func (self *PropertiesAPIModel) Init() {
	self.Names = []string{}
	self.Data = map[string]PropertyAPIModel{}
}

// Return a property matching given name, nil if not exists.
func (self *PropertiesAPIModel) Get(
	name string,
) PropertyAPIModel {
	return self.Data[name]
}

// Set the property by name, this will remember the order of insertion.
// The initial insertion order is kept even if the property is overwritten.
// Returns a boolean indicating if the value is newly inserted (not overwritten).
func (self *PropertiesAPIModel) Set(
	name string,
	property PropertyAPIModel,
) bool {
	_, exists := self.Data[name]
	if !exists {
		self.Names = append(self.Names, name)
	}
	self.Data[name] = property
	return !exists
}

// Drop the property if present.
// Return the property (if found) and a flag indicating if the property was found.
func (self *PropertiesAPIModel) Pop(
	name string,
) (PropertyAPIModel, bool) {
	property, exists := self.Data[name]
	if exists {
		// Find and remove property from names
		for index, content := range self.Names {
			if content == name {
				self.Names = append(self.Names[:index], self.Names[index+1:]...)
				break
			}
		}
		delete(self.Data, name)
	}
	return property, exists
}

// Return a slice with given the name, property pair in insertion order.
func (self *PropertiesAPIModel) Items() []PropertiesAPIModelItem {
	items := []PropertiesAPIModelItem{}
	for _, name := range self.Names {
		items = append(items, PropertiesAPIModelItem{name, self.Data[name]})
	}
	return items
}

// Implement type json.Marshaler interface. Will be called when marshaling PropertiesAPIModel.
func (self PropertiesAPIModel) MarshalJSON() ([]byte, error) {
	data := []byte{'{'}
	items := self.Items()
	last := len(items) - 1
	for index, item := range items {
		data = append(data, fmt.Sprintf("%q:", item.Name)...)
		propertyJson, err := json.Marshal(item.Property)
		if err != nil {
			return data, err
		}
		data = append(data, propertyJson...)
		if index < last {
			data = append(data, ',')
		}
	}
	data = append(data, '}')
	return data, nil
}

// Implement type json.Unmarshaler interface. Will be called when unmarshaling PropertiesAPIModel.
func (self *PropertiesAPIModel) UnmarshalJSON(data []byte) error {
	self.Init()

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	// Must open with a delim token '{'
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delim, ok := token.(json.Delim)
	if !ok || delim != '{' {
		return fmt.Errorf("expect JSON object open with '{'")
	}

	for decoder.More() {
		// Must be a string for the key (name)
		token, err = decoder.Token()
		if err != nil {
			return err
		}
		name, ok := token.(string)
		if !ok {
			return fmt.Errorf("expecting JSON key should be always a string: %T: %v", token, token)
		}

		var property PropertyAPIModel
		err = decoder.Decode(&property)
		if err != nil {
			return err
		}

		self.Set(name, property)
	}

	// Must end with a delim token '}'
	token, err = decoder.Token()
	if err != nil {
		return err
	}
	delim, ok = token.(json.Delim)
	if !ok || delim != '}' {
		return fmt.Errorf("expect JSON object close with '}'")
	}

	// Must be the end of the document
	token, err = decoder.Token()
	if err != io.EOF {
		return fmt.Errorf(
			"expect end of JSON object but got more token: %T: %v or err: %v",
			token, token, err)
	}

	return nil
}

// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TagModel describes the resource data model.
type TagModel struct {
	Id    types.String `tfsdk:"id"`
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`

	ForceDelete   types.Bool `tfsdk:"force_delete"`
	KeepOnDestroy types.Bool `tfsdk:"keep_on_destroy"`
}

// TagAPIModel describes the resource API model.
type TagAPIModel struct {
	Id    string `json:"id,omitempty"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type TagListAPIModel struct {
	Content          []TagAPIModel `json:"content"`
	TotalElements    int           `json:"totalElements"`
	NumberOfElements int           `json:"numberOfElements"`
}

func (self TagModel) String() string {
	return fmt.Sprintf(
		"Tag %s (%s)",
		self.Id.ValueString(),
		self.Key.ValueString())
}

// Return an appropriate key that can be used for naming mutexes.
// Create: Identifier can be used to prevent concurrent creation of tags.
// Read Update Delete: Identifier can be used to prevent concurrent modifications on the instance.
func (self TagModel) LockKey() string {
	return "tag-" + self.Id.ValueString()
}

func (self TagModel) ListPath() string {
	return "iaas/api/tags"
}

func (self TagModel) CreatePath() string {
	return "iaas/api/tags"
}

func (self TagModel) ReadPath() string {
	return ""
}

func (self TagModel) UpdatePath() string {
	panic(fmt.Sprintf("Cannot update %s, this type of resource is immutable.", self.String()))
}

func (self TagModel) DeletePath() string {
	path := "iaas/api/tags/" + self.Id.ValueString()
	if self.ForceDelete.ValueBool() {
		return path + "?ignoreUsage=true"
	}
	return path
}

func (self *TagModel) FromAPI(raw TagAPIModel) {
	self.Id = types.StringValue(raw.Id)
	self.Key = types.StringValue(raw.Key)
	self.Value = types.StringValue(raw.Value)
}

func (self TagModel) ToAPI() TagAPIModel {
	return TagAPIModel{
		Id:    self.Id.ValueString(),
		Key:   self.Key.ValueString(),
		Value: self.Value.ValueString(),
	}
}

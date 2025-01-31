// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CustomNamingModel describes the resource data model.
type CustomNamingModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	Projects  []CustomNamingProjectFilterModel     `tfsdk:"projects"`
	Templates map[string]CustomNamingTemplateModel `tfsdk:"templates"`
}

// CustomNamingAPIModel describes the resource API model.
type CustomNamingAPIModel struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Projects  []CustomNamingProjectFilterAPIModel `json:"projects"`
	Templates []CustomNamingTemplateAPIModel      `json:"templates"`
}

func (self *CustomNamingModel) String() string {
	return fmt.Sprintf(
		"Custom Naming %s (%s)",
		self.Id.ValueString(),
		self.Name.ValueString())
}

// Return an appropriate key that can be used for naming mutexes.
// Create: Identifier can be used to prevent concurrent creation of custom naming.
// Read Update Delete: Identifier can be used to prevent concurrent modifications on the instance.
func (self CustomNamingModel) LockKey() string {
	return "custom-naming-" + self.Id.ValueString()
}

func (self CustomNamingModel) CreatePath() string {
	return "iaas/api/naming"
}

func (self CustomNamingModel) ReadPath() string {
	return "iaas/api/naming/" + self.Id.ValueString()
}

func (self CustomNamingModel) UpdatePath() string {
	return self.CreatePath() // Its not a mistake...
}

func (self CustomNamingModel) DeletePath() string {
	return self.ReadPath()
}

func (self *CustomNamingModel) FromAPI(raw CustomNamingAPIModel) {

	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Description = types.StringValue(raw.Description)

	self.Projects = []CustomNamingProjectFilterModel{}
	for _, projectRaw := range raw.Projects {
		project := CustomNamingProjectFilterModel{}
		project.FromAPI(projectRaw)
		self.Projects = append(self.Projects, project)
	}

	// Match templates by resource type and static pattern
	self.Templates = map[string]CustomNamingTemplateModel{}
	for _, templateRaw := range raw.Templates {
		template := CustomNamingTemplateModel{}
		template.FromAPI(templateRaw)
		self.Templates[template.Key()] = template
	}
}

func (self *CustomNamingModel) ToAPI(state CustomNamingModel) CustomNamingAPIModel {

	projectsRaw := []CustomNamingProjectFilterAPIModel{}
	for _, project := range self.Projects {
		projectsRaw = append(projectsRaw, project.ToAPI())
	}

	templatesRaw := []CustomNamingTemplateAPIModel{}
	for key, template := range self.Templates {
		templateState := state.Templates[key]
		templatesRaw = append(templatesRaw, template.ToAPI(templateState))
	}

	return CustomNamingAPIModel{
		Id:          self.Id.ValueString(),
		Name:        self.Name.ValueString(),
		Description: CleanString(self.Description.ValueString()),
		Projects:    projectsRaw,
		Templates:   templatesRaw,
	}
}

// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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

func (self *CustomNamingModel) FromAPI(
	ctx context.Context,
	raw CustomNamingAPIModel,
) diag.Diagnostics {

	diags := diag.Diagnostics{}

	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Description = types.StringValue(raw.Description)

	self.Projects = []CustomNamingProjectFilterModel{}
	for _, projectRaw := range raw.Projects {
		project := CustomNamingProjectFilterModel{}
		diags.Append(project.FromAPI(ctx, projectRaw)...)
		self.Projects = append(self.Projects, project)
	}

	// Match templates by resource type and static pattern
	self.Templates = map[string]CustomNamingTemplateModel{}
	for _, templateRaw := range raw.Templates {
		template := CustomNamingTemplateModel{}
		diags.Append(template.FromAPI(ctx, templateRaw)...)
		self.Templates[template.Key()] = template
	}

	return diags
}

func (self *CustomNamingModel) ToAPI(
	ctx context.Context,
	state CustomNamingModel,
) (CustomNamingAPIModel, diag.Diagnostics) {

	diags := diag.Diagnostics{}

	projectsRaw := []CustomNamingProjectFilterAPIModel{}
	for _, project := range self.Projects {
		projectRaw, projectDiags := project.ToAPI(ctx)
		projectsRaw = append(projectsRaw, projectRaw)
		diags.Append(projectDiags...)
	}

	templatesRaw := []CustomNamingTemplateAPIModel{}
	for key, template := range self.Templates {
		templateState := state.Templates[key]
		templateRaw, templateDiags := template.ToAPI(ctx, templateState)
		templatesRaw = append(templatesRaw, templateRaw)
		diags.Append(templateDiags...)
	}

	raw := CustomNamingAPIModel{
		Id:          self.Id.ValueString(),
		Name:        self.Name.ValueString(),
		Description: CleanString(self.Description.ValueString()),
		Projects:    projectsRaw,
		Templates:   templatesRaw,
	}

	return raw, diags
}

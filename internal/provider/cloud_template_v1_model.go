// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gopkg.in/yaml.v2"
)

// CloudTemplateV1Model describes the resource data model.
type CloudTemplateV1Model struct {
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	RequestScopeOrg types.Bool   `tfsdk:"request_scope_org"`
	Status          types.String `tfsdk:"status"`

	Inputs    UnorderedPropertiesModel    `tfsdk:"inputs"`
	Resources CloudTemplateResourcesModel `tfsdk:"resources"`

	Valid              types.Bool `tfsdk:"valid"`
	ValidationMessages types.List `tfsdk:"validation_messages"`
	// Of type CloudTemplateV1ValidationMessageModel

	ProjectId types.String `tfsdk:"project_id"`
	OrgId     types.String `tfsdk:"org_id"`
}

// CloudTemplateV1APIModel describes the resource API model.
type CloudTemplateV1APIModel struct {
	Id              string `json:"id,omitempty"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	RequestScopeOrg bool   `json:"requestScopeOrg"`
	Status          string `json:"status"`
	Content         string `json:"content"`

	Valid              bool                                       `json:"valid"`
	ValidationMessages []CloudTemplateV1ValidationMessageAPIModel `json:"ValidationMessages,omitempty"`

	ProjectId string `json:"projectId"`
	OrgId     string `json:"orgId,omitempty"`
}

func (self CloudTemplateV1Model) String() string {
	return fmt.Sprintf(
		"Cloud Template v1 %s (%s)",
		self.Id.ValueString(),
		self.Name.ValueString())
}

// Return an appropriate key that can be used for naming mutexes.
// Create: Identifier can be used to prevent concurrent creation of cloud templates.
// Read Update Delete: Identifier can be used to prevent concurrent modifications on the instance.
func (self CloudTemplateV1Model) LockKey() string {
	return "cloud-template-" + self.Id.ValueString()
}

func (self CloudTemplateV1Model) CreatePath() string {
	return "blueprint/api/blueprints"
}

func (self CloudTemplateV1Model) ReadPath() string {
	return "blueprint/api/blueprints/" + self.Id.ValueString()
}

func (self CloudTemplateV1Model) UpdatePath() string {
	return self.ReadPath()
}

func (self CloudTemplateV1Model) DeletePath() string {
	return self.ReadPath()
}

func (self *CloudTemplateV1Model) FromAPI(
	ctx context.Context,
	raw CloudTemplateV1APIModel,
) diag.Diagnostics {

	diags := diag.Diagnostics{}

	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Description = types.StringValue(raw.Description)
	self.RequestScopeOrg = types.BoolValue(raw.RequestScopeOrg)
	self.Status = types.StringValue(raw.Status)
	self.Valid = types.BoolValue(raw.Valid)
	self.ProjectId = types.StringValue(raw.ProjectId)
	self.OrgId = types.StringValue(raw.OrgId)

	// Extract inputs and resources from raw content
	var contentRaw CloudTemplateV1ContentAPIModel
	err := yaml.Unmarshal([]byte(raw.Content), &contentRaw)
	if err == nil {
		diags.Append(self.Inputs.FromAPI(ctx, contentRaw.Inputs)...)
		diags.Append(self.Resources.FromAPI(ctx, contentRaw.Resources)...)
	} else {
		diags.AddError(
			"Configuration error",
			fmt.Sprintf("Unable to YAML decode %s content", self.String()))
	}

	// Convert raw validation messages to a ListValue of objects
	messages := []CloudTemplateV1ValidationMessageModel{}
	for _, messageRaw := range raw.ValidationMessages {
		message := CloudTemplateV1ValidationMessageModel{}
		message.FromAPI(messageRaw)
		messages = append(messages, message)
	}

	var messagesDiags diag.Diagnostics
	attrs := types.ObjectType{AttrTypes: CloudTemplateV1ValidationMessageModel{}.AttributeTypes()}
	self.ValidationMessages, messagesDiags = types.ListValueFrom(ctx, attrs, messages)
	diags.Append(messagesDiags...)

	return diags
}

func (self CloudTemplateV1Model) ToAPI(
	ctx context.Context,
) (CloudTemplateV1APIModel, diag.Diagnostics) {

	// Convert inputs and resources to raw content
	inputsRaw, diags := self.Inputs.ToAPI(ctx)
	resourcesRaw, resourcesDiags := self.Resources.ToAPI(ctx)
	diags.Append(resourcesDiags...)

	contentRaw := CloudTemplateV1ContentAPIModel{
		Inputs:    inputsRaw,
		Resources: resourcesRaw,
	}

	contentRawBytes, err := yaml.Marshal(contentRaw)
	if err != nil {
		diags.AddError(
			"Configuration error",
			fmt.Sprintf("Unable to YAML encode %s content", self.String()))
	}

	// Convert validation messages to raw API
	/*messagesRaw := []CloudTemplateV1ValidationMessageAPIModel{}
	for key, obj := range self.ValidationMessages {
		messageRaw, messageDiags := CloudTemplateV1ValidationMessageAPIModelFromObject(ctx, obj)
		messagesRaw = append(messagesRaw, messageRaw)
		diags.Append(messageDiags...)
	}*/

	return CloudTemplateV1APIModel{
		Id:              self.Id.ValueString(),
		Name:            self.Name.ValueString(),
		Description:     CleanString(self.Description.ValueString()),
		Content:         string(contentRawBytes),
		RequestScopeOrg: self.RequestScopeOrg.ValueBool(),
		Status:          self.Status.ValueString(),
		Valid:           self.Valid.ValueBool(),
		/*ValidationMessages: messagesRaw,*/
		ProjectId: self.ProjectId.ValueString(),
		OrgId:     self.OrgId.ValueString(),
	}, diags
}

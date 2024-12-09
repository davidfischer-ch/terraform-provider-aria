// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// OrchestratorWorkflowModel describes the resource data model.
type OrchestratorWorkflowModel struct {
	// Schema

	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	CategoryId  types.String `tfsdk:"category_id"`
	Version     types.String `tfsdk:"version"`

	//Attrib               jsontypes.Normalized `tfsdk:"attrib"`
	AllowedOperations types.String `tfsdk:"allowed_operations"`
	ObjectName        types.String `tfsdk:"object_name"`

	// TODO types.Object ... "Of type PositionModel"
	Position PositionModel `tfsdk:"position"`

	//Presentation         jsontypes.Normalized `tfsdk:"presentation"`
	RestartMode          types.Int32  `tfsdk:"restart_mode"`
	ResumeFromFailedMode types.Int32  `tfsdk:"resume_from_failed_mode"`
	RootName             types.String `tfsdk:"root_name"`
	//WorkflowItem         jsontypes.Normalized `tfsdk:"workflow_item"`

	InputParameters  types.List `tfsdk:"input_parameters"`
	OutputParameters types.List `tfsdk:"output_parameters"`
	// Of type ParameterModel

	ApiVersion    types.String `tfsdk:"api_version"`
	EditorVersion types.String `tfsdk:"editor_version"`

	// Of type OrchestratorWorkflowInputForm
	//InputForms types.List `tfsdk:"input_forms"`

	ForceDelete types.Bool `tfsdk:"force_delete"`
}

// OrchestratorWorkflowCreateAPIModel describes the resource create API model.
type OrchestratorWorkflowCreateAPIModel struct {
	Id         string `json:"id,omitempty"`
	Name       string `json:"name"`
	CategoryId string `json:"category-id"`
}

// OrchestratorWorkflowContentAPIModel describes the version API model.
type OrchestratorWorkflowContentAPIModel struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"display-name"`
	Description string `json:"description"`
	CategoryId  string `json:"category-id"`
	Version     string `json:"version"`

	Attrib               any              `json:"attrib"`
	AllowedOperations    string           `json:"allowed-operations"`
	ObjectName           string           `json:"object-name"`
	Position             PositionAPIModel `json:"position"`
	Presentation         any              `json:"presentation,omitempty"` // e.g. {}
	RestartMode          int32            `json:"restartMode"`
	ResumeFromFailedMode int32            `json:"resumeFromFailedMode"`
	RootName             string           `json:"root-name"`
	WorkflowItem         any              `json:"workflow-item"`

	Input  OrchestratorWorkflowIOAPIModel `json:"input"`
	Output OrchestratorWorkflowIOAPIModel `json:"output"`

	ApiVersion    string `json:"api-version"`
	EditorVersion string `json:"editor-version"`
}

type OrchestratorWorkflowIOAPIModel struct {
	Param []ParameterAPIModel `json:"param"`
}

type OrchestratorWorkflowVersionAPIModel struct {
	InputForms []map[string]any                    `json:"inputForms"`
	ParentId   string                              `json:"parentId,omitempty"`
	Schema     OrchestratorWorkflowContentAPIModel `json:"workflowSchema"`
}

type OrchestratorWorkflowVersionResponseAPIModel struct {
	ObjectId string `json:"objectId"`
}

func (self OrchestratorWorkflowModel) String() string {
	return fmt.Sprintf(
		"Orchestrator Workflow %s (%s)",
		self.Id.ValueString(),
		self.Name.ValueString())
}

// Return an appropriate key that can be used for naming mutexes.
// Create: Identifier can be used to prevent concurrent creation of vRO actions.
// Read Update Delete: Identifier can be used to prevent concurrent modifications on the instance.
func (self OrchestratorWorkflowModel) LockKey() string {
	return "orchestrator-workflow-" + self.Id.ValueString()
}

func (self OrchestratorWorkflowModel) CreatePath() string {
	return "vco/api/workflows"
}

func (self OrchestratorWorkflowModel) ReadPath() string {
	return self.ReadContentPath()
}

func (self OrchestratorWorkflowModel) ReadContentPath() string {
	return fmt.Sprintf("vco/api/workflows/%s/content", self.Id.ValueString())
}

func (self OrchestratorWorkflowModel) ReadVersionsPath() string {
	return fmt.Sprintf("vco/api/workflows/%s/versions", self.Id.ValueString())
}

func (self OrchestratorWorkflowModel) ReadVersionPath(versionId string) string {
	return fmt.Sprintf("vco/api/workflows/%s/versions/%s", self.Id.ValueString(), versionId)
}

func (self OrchestratorWorkflowModel) UpdatePath() string {
	return self.ReadPath()
}

func (self OrchestratorWorkflowModel) UpdateContentPath() string {
	return self.ReadContentPath()
}

func (self OrchestratorWorkflowModel) DeletePath() string {
	if self.ForceDelete.ValueBool() {
		return fmt.Sprintf("vco/api/workflows/%s?force=true", self.Id.ValueString())
	}
	return fmt.Sprintf("vco/api/workflows/%s", self.Id.ValueString())
}

// Save response from create API endpoint (only ID, name and category attributes are available)
func (self *OrchestratorWorkflowModel) FromCreateAPI(
	ctx context.Context,
	raw OrchestratorWorkflowCreateAPIModel,
) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.CategoryId = types.StringValue(raw.CategoryId)
	return diag.Diagnostics{}
}

// Save response from content API endpoint
func (self *OrchestratorWorkflowModel) FromContentAPI(
	ctx context.Context,
	raw OrchestratorWorkflowContentAPIModel,
) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.Description = types.StringValue(raw.Description)
	// FIXME How to retrieve CategoryId ? Yet another API endpoint to call?
	self.Version = types.StringValue(raw.Version)
	self.AllowedOperations = types.StringValue(raw.AllowedOperations)
	self.ObjectName = types.StringValue(raw.ObjectName)
	self.RestartMode = types.Int32Value(raw.RestartMode)
	self.ResumeFromFailedMode = types.Int32Value(raw.ResumeFromFailedMode)
	self.RootName = types.StringValue(raw.RootName)
	self.ApiVersion = types.StringValue(raw.ApiVersion)
	self.EditorVersion = types.StringValue(raw.EditorVersion)

	var parametersDiags diag.Diagnostics
	diags := self.Position.FromAPI(ctx, raw.Position)

	self.InputParameters, parametersDiags = ParameterModelListFromAPI(ctx, raw.Input.Param)
	diags.Append(parametersDiags...)

	self.OutputParameters, parametersDiags = ParameterModelListFromAPI(ctx, raw.Output.Param)
	diags.Append(parametersDiags...)

	return diags
}

// Save response from version API endpoint
func (self *OrchestratorWorkflowModel) FromVersionAPI(
	ctx context.Context,
	raw OrchestratorWorkflowVersionAPIModel,
) diag.Diagnostics {
	return self.FromContentAPI(ctx, raw.Schema)
}

// Create data for calling the create API endpoint (only ID, name and category attributes are set).
func (self OrchestratorWorkflowModel) ToCreateAPI(
	ctx context.Context,
) (OrchestratorWorkflowCreateAPIModel, diag.Diagnostics) {
	return OrchestratorWorkflowCreateAPIModel{
		Id:         self.Id.ValueString(),
		Name:       self.Name.ValueString(),
		CategoryId: self.CategoryId.ValueString(),
	}, diag.Diagnostics{}
}

// Create data for calling the content API endpoint
func (self OrchestratorWorkflowModel) ToContentAPI(
	ctx context.Context,
) (OrchestratorWorkflowContentAPIModel, diag.Diagnostics) {
	positionRaw, diags := self.Position.ToAPI(ctx)

	inputRaw, inputDiags := ParameterModelListToAPI(
		ctx,
		self.InputParameters,
		fmt.Sprintf("%s, %s", self.String(), "input_parameters"),
	)
	diags.Append(inputDiags...)

	outputRaw, outputDiags := ParameterModelListToAPI(
		ctx,
		self.OutputParameters,
		fmt.Sprintf("%s, %s", self.String(), "output_parameters"),
	)
	diags.Append(outputDiags...)

	return OrchestratorWorkflowContentAPIModel{
		Id:                   self.Id.ValueString(),
		Name:                 self.Name.ValueString(),
		CategoryId:           self.CategoryId.ValueString(),
		Description:          self.Description.ValueString(),
		Version:              self.Version.ValueString(),
		AllowedOperations:    self.AllowedOperations.ValueString(),
		ObjectName:           self.ObjectName.ValueString(),
		Position:             positionRaw,
		RestartMode:          self.RestartMode.ValueInt32(),
		ResumeFromFailedMode: self.ResumeFromFailedMode.ValueInt32(),
		RootName:             self.RootName.ValueString(),
		Input: OrchestratorWorkflowIOAPIModel{
			Param: inputRaw,
		},
		Output: OrchestratorWorkflowIOAPIModel{
			Param: outputRaw,
		},
		ApiVersion:    self.ApiVersion.ValueString(),
		EditorVersion: self.EditorVersion.ValueString(),
	}, diags
}

// Create data for calling the version API endpoint
func (self OrchestratorWorkflowModel) ToVersionAPI(
	ctx context.Context,
) (OrchestratorWorkflowVersionAPIModel, diag.Diagnostics) {
	schema, diags := self.ToContentAPI(ctx)
	return OrchestratorWorkflowVersionAPIModel{
		InputForms: []map[string]any{},
		ParentId:   "",
		Schema:     schema,
	}, diags
}

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
	Id      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Version types.String `tfsdk:"version"`

	CategoryId types.String `tfsdk:"category_id"`

	//ApiVersion
	//EditorVersion
	//ObjectName
	//Position
	//Presentation
	RestartMode          types.Int32 `tfsdk:"restart_mode"`
	ResumeFromFailedMode types.Int32 `tfsdk:"resume_from_failed_mode"`
	//RootName

	// Other fields that are not yet enforced by this model
	//Schema jsontypes.Normalized `tfsdk:"schema"`

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

// OrchestratorWorkflowUpdateAPIModel describes the resource read & update API model.
// TODO Use versions API endpoint instead (and rename it OrchestratorWorkflowVersionAPIModel)
type OrchestratorWorkflowUpdateAPIModel struct {
	Id                   string           `json:"id,omitempty"`
	Name                 string           `json:"display-name"`
	Version              string           `json:"version"`              // e.g. "1.0.0"
	ApiVersion           string           `json:"api-version"`          // e.g. "6.0.0"
	EditorVersion        string           `json:"editor-version"`       // e.g. "2.0"
	ObjectName           string           `json:"object-name"`          // e.g. workflow:name=generic
	Position             PositionAPIModel `json:"position"`             // e.g. {"x": 100.0, "y": 50.0}
	Presentation         any              `json:"presentation"`         // e.g. {}
	RestartMode          int32            `json:"restartMode"`          // e.g. 1
	ResumeFromFailedMode int32            `json:"resumeFromFailedMode"` // e.g. 0
	RootName             string           `json:"root-name"`            // e.g. "item0"

	CategoryId string `json:"category-id"`

	// List of ... {
	//   "position": {"y":50.0,"x":240.0},
	//   "name":"item0",
	//   "type":"end",
	//   "end-mode":"0",
	//   "comparator":0
	// }
	WorkflowItem any `json:"workflow-item"`
}

type PositionAPIModel struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
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
	// TODO Use versions API endpoint instead
	return fmt.Sprintf("vco/api/workflows/%s/content", self.Id.ValueString())
}

func (self OrchestratorWorkflowModel) UpdatePath() string {
	return self.ReadPath()
}

func (self OrchestratorWorkflowModel) DeletePath() string {
	if self.ForceDelete.ValueBool() {
		return fmt.Sprintf("vco/api/workflows/%s?force=true", self.Id.ValueString())
	}
	return fmt.Sprintf("vco/api/workflows/%s", self.Id.ValueString())
}

// Save response from create API endpoint (only ID, name and category attributes are available).
func (self *OrchestratorWorkflowModel) FromCreateAPI(
	ctx context.Context,
	raw OrchestratorWorkflowCreateAPIModel,
) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	self.CategoryId = types.StringValue(raw.CategoryId)
	return diag.Diagnostics{}
}

// Save response from read & update API endpoint
// TODO Use versions API endpoint instead (and rename it FromVersionAPI)
func (self *OrchestratorWorkflowModel) FromUpdateAPI(
	ctx context.Context,
	raw OrchestratorWorkflowUpdateAPIModel,
) diag.Diagnostics {
	self.Id = types.StringValue(raw.Id)
	self.Name = types.StringValue(raw.Name)
	// FIXME How to retrieve CategoryId ? Yet another API endpoint to call?
	self.Version = types.StringValue(raw.Version)
	//ApiVersion
	//EditorVersion
	//ObjectName
	//Position
	//Presentation
	self.RestartMode = types.Int32Value(raw.RestartMode)
	self.ResumeFromFailedMode = types.Int32Value(raw.ResumeFromFailedMode)
	//RootName
	return diag.Diagnostics{}
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

// Create data for calling the update API endpoint
// TODO Use versions API endpoint instead (and rename it ToVersionAPI)
func (self OrchestratorWorkflowModel) ToUpdateAPI(
	ctx context.Context,
) (OrchestratorWorkflowUpdateAPIModel, diag.Diagnostics) {
	return OrchestratorWorkflowUpdateAPIModel{
		Id:      self.Id.ValueString(),
		Name:    self.Name.ValueString(),
		Version: self.Version.ValueString(),
		//ApiVersion
		//EditorVersion
		//ObjectName
		//Position
		//Presentation
		RestartMode:          self.RestartMode.ValueInt32(),
		ResumeFromFailedMode: self.ResumeFromFailedMode.ValueInt32(),
		//RootName
		CategoryId: self.CategoryId.ValueString(),
	}, diag.Diagnostics{}
}

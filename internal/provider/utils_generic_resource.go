// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// CRUDModel is a Model with context-based API conversions.
type CRUDModel[A any] interface {
	Model
	ToAPI(ctx context.Context) (A, diag.Diagnostics)
	FromAPI(ctx context.Context, raw A) diag.Diagnostics
}

// SimpleCRUDModel is a Model with simple API conversions (no context/diagnostics).
type SimpleCRUDModel[A any] interface {
	Model
	ToAPI() A
	FromAPI(raw A)
}

// GenericResourceConfig configures a generic CRUD resource.
type GenericResourceConfig struct {
	TypeName     string
	SchemaFunc   func() schema.Schema
	CreateCodes  []int  // Defaults to [201]
	UpdateMethod string // Defaults to "PUT"
	UpdateCodes  []int  // Defaults to [200]
	// Extra attributes to set during ImportState (attribute path -> default value).
	ImportStateSetAttributes map[string]string
}

func (c GenericResourceConfig) getUpdateMethod() string {
	if c.UpdateMethod == "" {
		return "PUT"
	}
	return c.UpdateMethod
}

// --- GenericResource handles standard CRUD for models with ctx-based conversions ---

type GenericResource[M any, PM interface {
	*M
	CRUDModel[A]
}, A any] struct {
	client *AriaClient
	config GenericResourceConfig
}

func (self *GenericResource[M, PM, A]) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + self.config.TypeName
}

func (self *GenericResource[M, PM, A]) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = self.config.SchemaFunc()
}

func (self *GenericResource[M, PM, A]) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *GenericResource[M, PM, A]) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var model M
	pm := PM(&model)
	resp.Diagnostics.Append(req.Plan.Get(ctx, pm)...)
	if resp.Diagnostics.HasError() {
		return
	}

	toAPI, diags := pm.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var raw A
	_, createDiags := self.client.CreateIt(pm, &raw, toAPI, self.config.CreateCodes...)
	resp.Diagnostics.Append(createDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(pm.FromAPI(ctx, raw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, pm)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", pm.String()))
}

func (self *GenericResource[M, PM, A]) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var model M
	pm := PM(&model)
	resp.Diagnostics.Append(req.State.Get(ctx, pm)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var raw A
	found, _, readDiags := self.client.ReadIt(pm, &raw)
	resp.Diagnostics.Append(readDiags...)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(pm.FromAPI(ctx, raw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, pm)...)
}

func (self *GenericResource[M, PM, A]) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var model M
	pm := PM(&model)
	resp.Diagnostics.Append(req.Plan.Get(ctx, pm)...)
	if resp.Diagnostics.HasError() {
		return
	}

	toAPI, diags := pm.ToAPI(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var raw A
	_, updateDiags := self.client.UpdateIt(
		pm, &raw, toAPI, self.config.getUpdateMethod(), self.config.UpdateCodes...)
	resp.Diagnostics.Append(updateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(pm.FromAPI(ctx, raw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, pm)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", pm.String()))
}

func (self *GenericResource[M, PM, A]) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var model M
	pm := PM(&model)
	resp.Diagnostics.Append(req.State.Get(ctx, pm)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(self.client.DeleteIt(pm)...)
	}
}

func (self *GenericResource[M, PM, A]) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	for attr, value := range self.config.ImportStateSetAttributes {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attr), value)...)
	}
}

// --- SimpleGenericResource handles CRUD for models with simple conversions ---

type SimpleGenericResource[M any, PM interface {
	*M
	SimpleCRUDModel[A]
}, A any] struct {
	client *AriaClient
	config GenericResourceConfig
}

func (self *SimpleGenericResource[M, PM, A]) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + self.config.TypeName
}

func (self *SimpleGenericResource[M, PM, A]) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = self.config.SchemaFunc()
}

func (self *SimpleGenericResource[M, PM, A]) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *SimpleGenericResource[M, PM, A]) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var model M
	pm := PM(&model)
	resp.Diagnostics.Append(req.Plan.Get(ctx, pm)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var raw A
	_, createDiags := self.client.CreateIt(pm, &raw, pm.ToAPI(), self.config.CreateCodes...)
	resp.Diagnostics.Append(createDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pm.FromAPI(raw)
	resp.Diagnostics.Append(resp.State.Set(ctx, pm)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", pm.String()))
}

func (self *SimpleGenericResource[M, PM, A]) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var model M
	pm := PM(&model)
	resp.Diagnostics.Append(req.State.Get(ctx, pm)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var raw A
	found, _, readDiags := self.client.ReadIt(pm, &raw)
	resp.Diagnostics.Append(readDiags...)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	pm.FromAPI(raw)
	resp.Diagnostics.Append(resp.State.Set(ctx, pm)...)
}

func (self *SimpleGenericResource[M, PM, A]) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var model M
	pm := PM(&model)
	resp.Diagnostics.Append(req.Plan.Get(ctx, pm)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var raw A
	_, updateDiags := self.client.UpdateIt(
		pm, &raw, pm.ToAPI(), self.config.getUpdateMethod(), self.config.UpdateCodes...)
	resp.Diagnostics.Append(updateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pm.FromAPI(raw)
	resp.Diagnostics.Append(resp.State.Set(ctx, pm)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", pm.String()))
}

func (self *SimpleGenericResource[M, PM, A]) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var model M
	pm := PM(&model)
	resp.Diagnostics.Append(req.State.Get(ctx, pm)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(self.client.DeleteIt(pm)...)
	}
}

func (self *SimpleGenericResource[M, PM, A]) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	for attr, value := range self.config.ImportStateSetAttributes {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attr), value)...)
	}
}

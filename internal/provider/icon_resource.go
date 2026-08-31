// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"crypto/sha256"
	"fmt"
	"mime"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// IconMIMEType returns the MIME type of an icon, derived from its extension.
//
// Resty types a multipart file part by sniffing its content, and Go's sniffer has no SVG rule: an
// SVG opens with `<?xml`, which it reports as text/xml. The VCF 9 icon API rejects any part that is
// not an image type with "file must contain valid MIME type", where Aria Automation 8.x accepted
// it. Deriving the type from the extension keeps both versions happy.
func IconMIMEType(path string) string {
	if mimeType := mime.TypeByExtension(filepath.Ext(path)); len(mimeType) > 0 {
		return mimeType
	}
	return "application/octet-stream"
}

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IconResource{}

func NewIconResource() resource.Resource {
	return &IconResource{}
}

// IconResource defines the resource implementation.
type IconResource struct {
	client *AriaClient
}

func (self *IconResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_icon"
}

func (self *IconResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = IconSchema()
}

func (self *IconResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

func (self *IconResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var icon IconModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &icon)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Creating "icons" with the same content multiple times in parallel
	// will lead to a platform's Internal Error (HTTP 500)...
	// The platform is not handling properly concurrent requests to icon create/delete API
	// So we implement this protection (mutex) at the client side (provider)

	iconPath := icon.Path.ValueString()
	file, err := os.Open(iconPath)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to open %s, got error: %s", icon.String(), err))
		return
	}
	// Read-only file, nothing to report on close.
	defer func() { _ = file.Close() }()

	lockKey := icon.LockKey()
	path := icon.CreatePath()
	self.client.Mutex.Lock(ctx, lockKey)
	defer self.client.Mutex.Unlock(ctx, lockKey)
	response, err := self.client.R(path).
		SetMultipartField("file", filepath.Base(iconPath), IconMIMEType(iconPath), file).
		Post(path)
	err = self.client.HandleAPIResponse(response, err, []int{201})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create %s, got error: %s", icon.String(), err))
		return
	}

	iconId, err := self.client.GetIdFromLocation(response)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to parse Icon ID, got error: %s", err))
		return
	}

	// Save icon into Terraform state
	icon.Id = types.StringValue(iconId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &icon)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", icon.String()))

	// Read the icon to retrieve its content (duplicated code with read)

	path = icon.ReadPath()
	response, err = self.client.R(path).Get(path)
	err = self.client.HandleAPIResponse(response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read %s, got error: %s", icon.String(), err))
		return
	}

	// Save updated icon into Terraform state
	icon.ContentHash = types.StringValue(fmt.Sprintf("%x", sha256.Sum256(response.Body())))
	resp.Diagnostics.Append(resp.State.Set(ctx, &icon)...)
	tflog.Debug(ctx, fmt.Sprintf("Refreshed %s successfully", icon.String()))

}

func (self *IconResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var icon IconModel
	resp.Diagnostics.Append(req.State.Get(ctx, &icon)...)
	if resp.Diagnostics.HasError() {
		return
	}

	path := icon.ReadPath()
	self.client.Mutex.RLock(ctx, icon.LockKey())
	defer self.client.Mutex.RUnlock(ctx, icon.LockKey())
	response, err := self.client.R(path).Get(path)

	// Handle gracefully a resource that has vanished on the platform
	// Beware that some APIs respond with HTTP 404 instead of 403 ...
	if response.StatusCode() == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	err = self.client.HandleAPIResponse(response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read %s, got error: %s", icon.String(), err))
		return
	}

	// Save updated icon into Terraform state
	icon.ContentHash = types.StringValue(fmt.Sprintf("%x", sha256.Sum256(response.Body())))
	resp.Diagnostics.Append(resp.State.Set(ctx, &icon)...)
}

func (self *IconResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var icon IconModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &icon)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated icon into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &icon)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", icon.String()))
}

func (self *IconResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var icon IconModel
	resp.Diagnostics.Append(req.State.Get(ctx, &icon)...)
	if !resp.Diagnostics.HasError() && !icon.KeepOnDestroy.ValueBool() {
		self.client.Mutex.Lock(ctx, icon.LockKey())
		defer self.client.Mutex.Unlock(ctx, icon.LockKey())
		resp.Diagnostics.Append(self.client.DeleteIt(&icon)...)
	}
}

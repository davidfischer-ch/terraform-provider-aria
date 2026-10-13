// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SecretResource{}
var _ resource.ResourceWithImportState = &SecretResource{}
var _ resource.ResourceWithModifyPlan = &SecretResource{}

func NewSecretResource() resource.Resource {
	return &SecretResource{}
}

// SecretResource defines the resource implementation.
type SecretResource struct {
	client *AriaClient
}

func (self *SecretResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_secret"
}

func (self *SecretResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = SecretSchema()
}

func (self *SecretResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	self.client = GetResourceClient(ctx, req, resp)
}

// ModifyPlan proposes an update when the configured value is not the one last written.
//
// The value is write-only, so the plan never differs on it: the plan is made different by hand,
// through the attributes that an update changes.
func (self *SecretResource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	// Nothing to do on create (everything is unknown) and on destroy.
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	var value, checksum types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("value"), &value)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("value_checksum"), &checksum)...)
	if resp.Diagnostics.HasError() || !SecretValueChanged(value, checksum) {
		return
	}

	for _, name := range []string{"value_checksum", "updated_at", "updated_by"} {
		resp.Diagnostics.Append(
			resp.Plan.SetAttribute(ctx, path.Root(name), types.StringUnknown())...)
	}
}

func (self *SecretResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Read Terraform plan data into the model
	var secret SecretModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &secret)...)

	// The value is write-only, it is only available from the configuration
	var value types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("value"), &value)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := secret.ToAPI(ctx, value.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var secretFromAPI SecretAPIModel
	_, createDiags := self.client.CreateIt(&secret, &secretFromAPI, body, 201)
	resp.Diagnostics.Append(createDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save secret into Terraform state right now, so a failure below does not leave it orphan
	resp.Diagnostics.Append(secret.FromAPI(ctx, secretFromAPI)...)
	secret.ValueChecksum = types.StringValue(SecretValueChecksum(value.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &secret)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save the secret as read back into Terraform state
	resp.Diagnostics.Append(self.refresh(ctx, &secret)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &secret)...)
	tflog.Debug(ctx, fmt.Sprintf("Created %s successfully", secret.String()))
}

func (self *SecretResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Read Terraform prior state data into the model
	var secret SecretModel
	resp.Diagnostics.Append(req.State.Get(ctx, &secret)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var secretFromAPI SecretAPIModel
	found, _, readDiags := self.client.ReadIt(&secret, &secretFromAPI)
	resp.Diagnostics.Append(readDiags...)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Somebody else modified the secret: forget which value we wrote, the next plan then proposes
	// to overwrite the secret.
	if SecretModifiedOutside(secret.UpdatedAt, secretFromAPI.UpdatedAt) {
		tflog.Warn(ctx, fmt.Sprintf(
			"%s was modified outside of Terraform (updated_at %s, expected %s)",
			secret.String(), secretFromAPI.UpdatedAt, secret.UpdatedAt.ValueString()))
		secret.ValueChecksum = types.StringNull()
	}

	// Save updated secret into Terraform state
	resp.Diagnostics.Append(secret.FromAPI(ctx, secretFromAPI)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &secret)...)
}

func (self *SecretResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Read Terraform plan data into the model
	var secret SecretModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &secret)...)

	// The value is write-only, it is only available from the configuration
	var value types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("value"), &value)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := secret.ToAPI(ctx, value.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	self.client.Mutex.Lock(ctx, secret.LockKey())
	defer self.client.Mutex.Unlock(ctx, secret.LockKey())

	// Only PATCH is supported, the reply is not trusted: it can be missing fields
	var secretFromAPI SecretAPIModel
	_, updateDiags := self.client.UpdateIt(&secret, &secretFromAPI, body, "PATCH")
	resp.Diagnostics.Append(updateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save secret into Terraform state
	secret.ValueChecksum = types.StringValue(SecretValueChecksum(value.ValueString()))
	resp.Diagnostics.Append(self.refresh(ctx, &secret)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &secret)...)
	tflog.Debug(ctx, fmt.Sprintf("Updated %s successfully", secret.String()))
}

func (self *SecretResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Read Terraform prior state data into the model
	var secret SecretModel
	resp.Diagnostics.Append(req.State.Get(ctx, &secret)...)
	if !resp.Diagnostics.HasError() {
		self.client.Mutex.Lock(ctx, secret.LockKey())
		defer self.client.Mutex.Unlock(ctx, secret.LockKey())
		resp.Diagnostics.Append(self.client.DeleteIt(&secret)...)
	}
}

func (self *SecretResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// refresh reads the secret back into the model, after it was written.
//
// The reply of a write is not what a read returns: the creation date-time is rounded to the
// millisecond and `orgScoped` goes missing after an update. Keeping it would make the next
// refresh believe that somebody modified the secret.
func (self *SecretResource) refresh(ctx context.Context, secret *SecretModel) diag.Diagnostics {
	var secretFromAPI SecretAPIModel
	found, _, diags := self.client.ReadIt(secret, &secretFromAPI)
	if diags.HasError() {
		return diags
	}
	if !found {
		diags.AddError(
			"Client error",
			fmt.Sprintf("Unable to read %s after writing it, not found", secret.String()))
		return diags
	}
	diags.Append(secret.FromAPI(ctx, secretFromAPI)...)
	return diags
}

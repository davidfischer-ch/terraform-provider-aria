// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ABXSecretResource{}

func NewABXSecretResource() resource.Resource {
	return &ABXSecretResource{}
}

// ABXSecretResource defines the resource implementation.
type ABXSecretResource struct {
	client *resty.Client
}

// ABXSecretResourceModel describes the resource data model.
type ABXSecretResourceModel struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Value     types.String `tfsdk:"value"`
	Encrypted types.Bool   `tfsdk:"encrypted"`
	OrgId     types.String `tfsdk:"org_id"`
}

// ABXSecretResourceAPIModel describes the resource API model.
type ABXSecretResourceAPIModel struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Value         string `json:"value"`
	Encrypted     bool   `json:"encrypted"`
	OrgId         string `json:"orgId"`
	CreatedMillis int64  `json:"createdMillis"`
}

func (r *ABXSecretResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_abx_secret"
}

func (r *ABXSecretResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "ABX secret resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Secret identifier",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Secret name",
				Required:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "Secret value (cannot be enforced since API don't return it)",
				Required:            true,
				Sensitive:           true,
			},
			"encrypted": schema.BoolAttribute{
				MarkdownDescription: "Secret should be always encrypted!",
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"org_id": schema.StringAttribute{
				MarkdownDescription: "Secret organisation ID",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (r *ABXSecretResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = GetResourceClient(ctx, req, resp)
}

func (r *ABXSecretResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var secret ABXSecretResourceModel
	var secretRaw ABXSecretResourceAPIModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &secret)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.R().
		SetBody(ABXSecretResourceAPIModel{
			Name:      secret.Name.ValueString(),
			Value:     secret.Value.ValueString(),
			Encrypted: secret.Encrypted.ValueBool(),
		}).
		SetResult(&secretRaw).
		Post("abx/api/resources/action-secrets")

	err = handleAPIResponse(ctx, response, err, 200)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to create secret, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Secret %s created", secretRaw.Id))

	secret.Id = types.StringValue(secretRaw.Id)
	secret.Encrypted = types.BoolValue(secretRaw.Encrypted)
	secret.OrgId = types.StringValue(secretRaw.OrgId)

	// Save secret into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &secret)...)
}

func (r *ABXSecretResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var secret ABXSecretResourceModel
	var secretRaw ABXSecretResourceAPIModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &secret)...)
	if resp.Diagnostics.HasError() {
		return
	}

	secretId := secret.Id.ValueString()
	response, err := r.client.R().
		SetResult(&secretRaw).
		Get("abx/api/resources/action-secrets/" + secretId)

	// Handle gracefully a resource that has vanished on the platform
	// Beware that some APIs respond with HTTP 404 instead of 403 ...
	if response.StatusCode() == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	err = handleAPIResponse(ctx, response, err, 200)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read secret %s, got error: %s", secretId, err))
		return
	}

	secret.Name = types.StringValue(secretRaw.Name)
	// secret.Value = types.StringValue(secretRaw.Value)
	secret.Encrypted = types.BoolValue(secretRaw.Encrypted)
	secret.OrgId = types.StringValue(secretRaw.OrgId)

	// Save updated secret into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &secret)...)
}

func (r *ABXSecretResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var secret ABXSecretResourceModel
	var secretRaw ABXSecretResourceAPIModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &secret)...)
	if resp.Diagnostics.HasError() {
		return
	}

	secretId := secret.Id.ValueString()
	response, err := r.client.R().
		SetBody(ABXSecretResourceAPIModel{
			Name:      secret.Name.ValueString(),
			Value:     secret.Value.ValueString(),
			Encrypted: secret.Encrypted.ValueBool(),
		}).
		SetResult(&secretRaw).
		Put("abx/api/resources/action-secrets/" + secretId)

	err = handleAPIResponse(ctx, response, err, 200)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to update secret %s, got error: %s", secretId, err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Secret %s update", secretId))

	secret.Name = types.StringValue(secretRaw.Name)
	// value is returned with the following '*****' awesome :)
	secret.Encrypted = types.BoolValue(secretRaw.Encrypted)
	secret.OrgId = types.StringValue(secretRaw.OrgId)

	// Save secret into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &secret)...)
}

func (r *ABXSecretResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var secret ABXSecretResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &secret)...)
	if resp.Diagnostics.HasError() {
		return
	}

	secretId := secret.Id.ValueString()
	if len(secretId) == 0 {
		return
	}

	response, err := r.client.R().Delete("abx/api/resources/action-secrets/" + secretId)
	err = handleAPIResponse(ctx, response, err, 200)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to delete secret %s, got error: %s", secretId, err))
	}
}

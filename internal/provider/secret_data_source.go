// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SecretDataSource{}

func NewSecretDataSource() datasource.DataSource {
	return &SecretDataSource{}
}

// SecretDataSource defines the data source implementation.
type SecretDataSource struct {
	client *resty.Client
}

func (self *SecretDataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_secret"
}

func (self *SecretDataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Secret data source",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Secret identifier",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Secret name",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Describe the secret in few sentences",
				Computed:            true,
			},
			"org_id": ComputedOrganizationIdSchema(),
			"org_scoped": schema.BoolAttribute{
				MarkdownDescription: "Scoped to the organization?",
				Computed:            true,
			},
			"project_ids": schema.SetAttribute{
				MarkdownDescription: "Restrict to given projects (an empty list means all)",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Creation date-time",
				Computed:            true,
			},
			"created_by": schema.StringAttribute{
				MarkdownDescription: "Ask VMware",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "Changed date-time",
				Computed:            true,
			},
			"updated_by": schema.StringAttribute{
				MarkdownDescription: "Ask VMware",
				Computed:            true,
			},
		},
	}
}

func (self *SecretDataSource) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	self.client = GetDataSourceClient(ctx, req, resp)
}

func (self *SecretDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	// Read Terraform configuration data into the model
	var secret SecretModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &secret)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var secretRaw SecretAPIModel
	secretId := secret.Id.ValueString()
	response, err := self.client.R().
		// TODO SetQueryParam("apiVersion", PLATFORM_API_VERSION).
		SetResult(&secretRaw).
		Get("/platform/api/secrets/" + secretId)

	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read secret %s, got error: %s", secretId, err))
		return
	}

	// Save updated secret into Terraform state
	resp.Diagnostics.Append(secret.FromAPI(ctx, secretRaw)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &secret)...)
}

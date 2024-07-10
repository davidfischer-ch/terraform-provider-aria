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
var _ datasource.DataSource = &IconDataSource{}

func NewIconDataSource() datasource.DataSource {
	return &IconDataSource{}
}

// IconDataSource defines the data source implementation.
type IconDataSource struct {
	client *resty.Client
}

func (self *IconDataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_icon"
}

func (self *IconDataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Icon data source",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Icon identifier",
				Required:            true,
			},
			"content": schema.StringAttribute{
				MarkdownDescription: "Icon content",
				Computed:            true,
			},
		},
	}
}

func (self *IconDataSource) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	self.client = GetDataSourceClient(ctx, req, resp)
}

func (self *IconDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	// Read Terraform configuration data into the model
	var icon IconModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &icon)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := self.client.R().
		// TODO SetQueryParam("apiVersion", ICON_API_VERSION).
		Get("icon/api/icons/" + icon.Id.ValueString())
	err = handleAPIResponse(ctx, response, err, []int{200})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read icon %s, got error: %s", icon.Id.ValueString(), err))
		return
	}

	// Save updated icon into Terraform state
	icon.Content = types.StringValue(response.String())
	resp.Diagnostics.Append(resp.State.Set(ctx, &icon)...)
}

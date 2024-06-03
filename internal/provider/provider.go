// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure AriaProvider satisfies various provider interfaces.
var _ provider.Provider = &AriaProvider{}
var _ provider.ProviderWithFunctions = &AriaProvider{}

// AriaProvider defines the provider implementation.
type AriaProvider struct {
	// version is set to the provider version on release, "dev" when the provider is built and ran
	// locally, and "test" when running acceptance testing.
	version string
}

// AriaProviderModel describes the provider data model.
type AriaProviderModel struct {
	Host         types.String `tfsdk:"host"`
	RefreshToken types.String `tfsdk:"refresh_token"`
	Insecure     types.Bool   `tfsdk:"insecure"`
}

func (p *AriaProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "aria"
	resp.Version = p.version
}

func (p *AriaProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The URI to Aria. May also be provided via ARIA_HOST environment variable.",
			},
			"refresh_token": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The refresh token to use for making API requests. May also be provided via ARIA_REFRESH_TOKEN environment variable.",
			},
			"insecure": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether server should be accessed without verifying the TLS certificate. May also be provided via ARIA_INSECURE environment variable.",
			},
		},
	}
}

func (p *AriaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	// Retrieve provider data from configuration
	var config AriaProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve default values from environment variables if set

	host := os.Getenv("ARIA_HOST")
	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}
	if len(host) == 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Aria API Host",
			"Set the host in the provider configuration "+
				"or use the ARIA_HOST and ensure its not empty.",
		)
	}

	refresh_token := os.Getenv("ARIA_REFRESH_TOKEN")
	if !config.RefreshToken.IsNull() {
		refresh_token = config.RefreshToken.ValueString()
	}
	if len(refresh_token) == 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("refresh_token"),
			"Missing Aria API Refresh Token",
			"Set the refresh token in the provider configuration "+
				"or use the ARIA_REFRESH_TOKEN and ensure its not empty.",
		)
	}

	var insecure bool
	var err error
	if !config.Insecure.IsNull() {
		insecure = config.Insecure.ValueBool()
	} else {
		insecure, err = strconv.ParseBool(os.Getenv("ARIA_INSECURE"))
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("insecure"),
				"Invalid Aria Insecure Flag",
				"Environment variable ARIA_INSECURE is not a valid boolean.",
			)
		}
	}

	ctx = tflog.SetField(ctx, "aria_host", host)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "aria_refresh_token", refresh_token)
	ctx = tflog.SetField(ctx, "aria_insecure", insecure)

	tflog.Debug(ctx, "Creating Aria client")

	// Create a new Aria client using the configuration values
	cfg := AriaClientConfig{
		Host:         host,
		RefreshToken: refresh_token,
		Insecure:     insecure,
	}

	err = cfg.Check()
	if err != nil {
		resp.Diagnostics.AddError(
			"Aria client configuration is invalid",
			err.Error(),
		)
		return
	}

	err = cfg.GetAccessToken()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to retrieve a valid access token",
			err.Error(),
		)
		return
	}

	client := cfg.Client()

	// Make the Aria client available for DataSource and Resource type Configure methods
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Aria client", map[string]any{"success": true})
}

func (p *AriaProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
	}
}

func (p *AriaProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
		// NewIconDataSource,
	}
}

func (p *AriaProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewExampleFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AriaProvider{
			version: version,
		}
	}
}

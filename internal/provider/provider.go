// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	Host               types.String `tfsdk:"host"`
	Insecure           types.Bool   `tfsdk:"insecure"`
	RefreshToken       types.String `tfsdk:"refresh_token"`
	AccessToken        types.String `tfsdk:"access_token"`
	OKAPICallsLogLevel types.String `tfsdk:"ok_api_calls_log_level"`
	KOAPICallsLogLevel types.String `tfsdk:"ko_api_calls_log_level"`
}

func (self *AriaProvider) Metadata(
	ctx context.Context,
	req provider.MetadataRequest,
	resp *provider.MetadataResponse,
) {
	resp.TypeName = "aria"
	resp.Version = self.version
}

func (self *AriaProvider) Schema(
	ctx context.Context,
	req provider.SchemaRequest,
	resp *provider.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "The URI to Aria. " +
					"May also be provided via ARIA_HOST environment variable.",
				Optional:            true,
			},
			"refresh_token": schema.StringAttribute{
				MarkdownDescription: "The refresh token to use for making API requests. " +
					"May also be provided via ARIA_REFRESH_TOKEN environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"access_token": schema.StringAttribute{
				MarkdownDescription: "The access token to use for making API requests. " +
					"May also be provided via ARIA_ACCESS_TOKEN environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"insecure": schema.BoolAttribute{
				MarkdownDescription: "Whether server should be accessed without verifying the " +
					"TLS certificate. May also be provided via ARIA_INSECURE environment variable.",
				Optional:            true,
			},
			"ok_api_calls_log_level": schema.StringAttribute{
				Optional:             true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"INFO", "DEBUG", "TRACE"}...),
				},
			},
			"ko_api_calls_log_level": schema.StringAttribute{
				Optional:             true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"INFO", "DEBUG", "TRACE"}...),
				},
			},
		},
	}
}

func (self *AriaProvider) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse,
) {

	// Retrieve provider data from configuration
	var config AriaProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prevent an unexpectedly misconfigured client, if Terraform configuration values are only
	// known after another resource is applied.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Aria API Host",
			"Either set the host in the provider configuration to a static value, "+
				"apply the source of the value first, or use ARIA_HOST.",
		)
	}

	if config.RefreshToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("refresh_token"),
			"Unknown Aria API Refresh Token",
			"Either set the refresh token in the provider configuration to a static value, "+
				"apply the source of the value first, or use ARIA_REFRESH_TOKEN.",
		)
	}

	if config.AccessToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("access_token"),
			"Unknown Aria API Access Token",
			"Either set the access token in the provider configuration to a static value, "+
				"apply the source of the value first, or use ARIA_ACCESS_TOKEN.",
		)
	}

	if config.Insecure.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("insecure"),
			"Unknown Aria API Insecure Flag",
			"Either set the insecure flag in the provider configuration to a static value, "+
				"apply the source of the value first, or use ARIA_INSECURE.",
		)
	}

	if config.OKAPICallsLogLevel.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("ok_api_calls_log_level"),
			"Unknown OK API Calls Log Level",
			"Either set the level in the provider configuration to a static value or "+
				"apply the source of the value first.",
		)
	}

	if config.KOAPICallsLogLevel.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("ko_api_calls_log_level"),
			"Unknown KO API Calls Log Level",
			"Either set the level in the provider configuration to a static value "+
				"or apply the source of the value first.",
		)
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
				"or use ARIA_HOST and ensure its not empty.",
		)
	}

	var insecure bool
	if !config.Insecure.IsNull() {
		insecure = config.Insecure.ValueBool()
	} else {
		var err error
		insecure, err = strconv.ParseBool(os.Getenv("ARIA_INSECURE"))
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("insecure"),
				"Invalid Aria Insecure Flag",
				"Environment variable ARIA_INSECURE is not a valid boolean.",
			)
		}
	}

	refresh_token := os.Getenv("ARIA_REFRESH_TOKEN")
	if !config.RefreshToken.IsNull() {
		refresh_token = config.RefreshToken.ValueString()
	}

	access_token := os.Getenv("ARIA_ACCESS_TOKEN")
	if !config.AccessToken.IsNull() {
		access_token = config.AccessToken.ValueString()
	}

	if len(refresh_token) == 0 && len(access_token) == 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("refresh_token"),
			"Missing Aria API Token",
			"Set either the refresh or access token in the provider configuration "+
				"or use one of ARIA_{ACCESS,REFRESH}_TOKEN and ensure its not empty.",
		)
	}

	okLogLevel := "TRACE"
	if !config.OKAPICallsLogLevel.IsNull() {
		okLogLevel = config.OKAPICallsLogLevel.ValueString()
	}

	koLogLevel := "DEBUG"
	if !config.KOAPICallsLogLevel.IsNull() {
		koLogLevel = config.KOAPICallsLogLevel.ValueString()
	}

	ctx = tflog.SetField(ctx, "aria_host", host)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "aria_refresh_token", refresh_token)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "aria_access_token", access_token)
	ctx = tflog.SetField(ctx, "aria_insecure", insecure)

	tflog.Debug(ctx, "Creating Aria client")

	// Create a new Aria client using the configuration values
	client := AriaClient{
		Host:               host,
		RefreshToken:       refresh_token,
		AccessToken:        access_token,
		Insecure:           insecure,
		Context:            ctx,
		OKAPICallsLogLevel: okLogLevel,
		KOAPICallsLogLevel: koLogLevel,
	}

	clientDiags := client.Init()
	resp.Diagnostics.Append(clientDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Make the Aria client available for DataSource and Resource type Configure methods
	resp.DataSourceData = &client
	resp.ResourceData = &client

	tflog.Info(ctx, "Configured Aria client", map[string]any{"success": true})
}

func (self *AriaProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewABXActionResource,
		NewABXConstantResource,
		NewABXSensitiveConstantResource,
		NewCatalogItemIconResource,
		NewCatalogSourceResource,
		NewCloudTemplateV1Resource,
		NewCustomFormResource,
		NewCustomNamingResource,
		NewCustomResourceResource,
		NewIconResource,
		NewOrchestratorActionResource,
		NewOrchestratorCategoryResource,
		NewOrchestratorConfigurationResource,
		NewOrchestratorEnvironmentResource,
		NewOrchestratorEnvironmentRepositoryResource,
		NewOrchestratorTaskResource,
		NewOrchestratorWorkflowResource,
		NewPolicyResource,
		NewProjectResource,
		NewPropertyGroupResource,
		NewResourceActionResource,
		NewSubscriptionResource,
		NewTagResource,
	}
}

func (self *AriaProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCatalogTypeDataSource,
		NewIconDataSource,
		NewIntegrationDataSource,
		NewOrchestratorConfigurationDataSource,
		NewSecretDataSource,
	}
}

func (self *AriaProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AriaProvider{
			version: version,
		}
	}
}

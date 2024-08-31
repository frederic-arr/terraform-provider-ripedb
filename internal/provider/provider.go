// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure RipeProvider satisfies various provider interfaces.
var _ provider.Provider = &RipeProvider{}
var _ provider.ProviderWithFunctions = &RipeProvider{}

// TODO: Add MD5 auth
// TODO: Add CERT auth

// RipeProvider defines the provider implementation.
type RipeProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// RipeProviderModel describes the provider data model.
type RipeProviderModel struct{}

type RipeProviderData struct{}

func (p *RipeProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ripe"
	resp.Version = p.version
}

func (p *RipeProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "RIPE DB.",
	}
}

func (p *RipeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data RipeProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	providerData := RipeProviderData{}
	resp.DataSourceData = &providerData
	resp.ResourceData = &providerData
}

func (p *RipeProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewObjectResource,
	}
}

func (p *RipeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewObjectDataSource,
	}
}

func (p *RipeProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewGetFirstFunction,
		NewGetAllFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &RipeProvider{
			version: version,
		}
	}
}

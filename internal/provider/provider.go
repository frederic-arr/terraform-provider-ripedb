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

// Ensure RipeDbProvider satisfies various provider interfaces.
var _ provider.Provider = &RipeDbProvider{}
var _ provider.ProviderWithFunctions = &RipeDbProvider{}

// TODO: Add MD5 auth
// TODO: Add CERT auth

// RipeDbProvider defines the provider implementation.
type RipeDbProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// RipeDbProviderModel describes the provider data model.
type RipeDbProviderModel struct{}

type RipeDbProviderData struct{}

func (p *RipeDbProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ripedb"
	resp.Version = p.version
}

func (p *RipeDbProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "RIPE DB.",
	}
}

func (p *RipeDbProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data RipeDbProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	providerData := RipeDbProviderData{}
	resp.DataSourceData = &providerData
	resp.ResourceData = &providerData
}

func (p *RipeDbProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewObjectResource,
	}
}

func (p *RipeDbProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewObjectDataSource,
	}
}

func (p *RipeDbProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewGetFirstFunction,
		NewGetAllFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &RipeDbProvider{
			version: version,
		}
	}
}

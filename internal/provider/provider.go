// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/frederic-arr/ripedb-go/ripedb"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure RipeDbProvider satisfies various provider interfaces.
var _ provider.Provider = &RipeDbProvider{}
var _ provider.ProviderWithFunctions = &RipeDbProvider{}

// TODO: Add MD5 auth
// TODO: Add CERT auth

const (
	DefaultEndpoint = "https://rest.db.ripe.net"
	DefaultSource   = "RIPE"
)

// RipeDbProvider defines the provider implementation.
type RipeDbProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// RipeDbProviderModel describes the provider data model.
type RipeDbProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Source   types.String `tfsdk:"database"`

	User     types.String `tfsdk:"user"`
	Password types.String `tfsdk:"password"`

	Certificate types.String `tfsdk:"certificate"`
	Key         types.String `tfsdk:"key"`
}

type RipeDbProviderData struct {
	Client *ripedb.RipeDbClient
}

func (p *RipeDbProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ripedb"
	resp.Version = p.version
}

func (p *RipeDbProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "RIPE DB.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "the endpoint of the RIPE DB",
				Optional:            true,
			},
			"database": schema.StringAttribute{
				MarkdownDescription: "the source of the RIPE DB",
				Optional:            true,
			},
			"user": schema.StringAttribute{
				MarkdownDescription: "the name of the MNTNER object",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "the password of the MNTNER object",
				Optional:            true,
				Sensitive:           true,
			},

			"certificate": schema.StringAttribute{
				MarkdownDescription: "the certificate of the MNTNER object",
				Optional:            true,
				Sensitive:           true,
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "the private key of the MNTNER object",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *RipeDbProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data RipeDbProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var endpoint string
	if data.Endpoint.ValueString() != "" {
		endpoint = data.Endpoint.ValueString()
	} else {
		endpoint = DefaultEndpoint
	}

	var source string
	if data.Source.ValueString() != "" {
		source = data.Source.ValueString()
	} else {
		source = DefaultSource
	}

	var client ripedb.RipeDbClient
	if !data.Password.IsNull() {
		client = ripedb.NewRipePasswordClient(data.User.ValueStringPointer(), data.Password.ValueString(), nil)
	} else if !data.Certificate.IsNull() || !data.Key.IsNull() {
		if data.Certificate.IsNull() || data.Key.IsNull() {
			resp.Diagnostics.AddError("both key and cert must be provided", "")
		}

		client = ripedb.NewRipeX509Client([]byte(data.Certificate.ValueString()), []byte(data.Key.ValueString()), nil)
	} else {
		client = ripedb.NewRipeAnonymousClient(nil)
	}

	client.SetEndpoint(endpoint)
	client.SetSource(source)

	providerData := RipeDbProviderData{
		Client: &client,
	}
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

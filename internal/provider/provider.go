// Copyright (c) The RIPE DB Provider for Terraform Authors
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
	Client *ripedb.RipeClient
}

func (p *RipeDbProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ripedb"
	resp.Version = p.version
}

func (p *RipeDbProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The RIPE DB provider is used to interact with the objects in the RIPE database. The provider needs to be configured with the proper credentials before objects can be modified.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "The endpoint of the RIPE Database RESTful API.",
				Optional:            true,
			},
			"database": schema.StringAttribute{
				MarkdownDescription: "The database where the queries should be made. This is equivalent to the `source` field of the objects.",
				Optional:            true,
			},

			"user": schema.StringAttribute{
				MarkdownDescription: "Username for the basic authentication protocol. You cannot use X.509 Client Authentication and Basic authentication at the same time.",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password for the basic authentication protocol. If no `user` is provided, the authentication will be made through the `password` query parameter instead of the `Authorizatio` header. You cannot use X.509 Client Authentication and Basic authentication at the same time.",
				Optional:            true,
				Sensitive:           true,
			},

			"certificate": schema.StringAttribute{
				MarkdownDescription: "PEM-encoded client certificate for TLS authentication. Both `certificate` and `key` must be provided. The `endpoint` field must be set appropriately if you are not using the default production API. You cannot use X.509 Client Authentication and Basic authentication at the same time.",
				Optional:            true,
				Sensitive:           true,
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "PEM-encoded client certificate key for TLS authentication. Both `certificate` and `key` must be provided. The `endpoint` field must be set appropriately if you are not using the default production API. You cannot use X.509 Client Authentication and Basic authentication at the same time.",
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

	opts := ripedb.RipeClientOptions{
		Endpoint: data.Endpoint.ValueStringPointer(),
		Source:   data.Source.ValueStringPointer(),
		User:     data.User.ValueStringPointer(),
		Password: data.Password.ValueStringPointer(),
	}

	if !data.Certificate.IsNull() {
		cert := []byte(data.Certificate.ValueString())
		opts.Certificate = &cert
	}
	if !data.Key.IsNull() {
		key := []byte(data.Key.ValueString())
		opts.Key = &key
	}

	client, err := ripedb.NewRipeClient(&opts)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create RIPE HTTP client", err.Error())
		return
	}

	providerData := RipeDbProviderData{
		Client: client,
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

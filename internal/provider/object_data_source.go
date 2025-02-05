// Copyright (c) The RIPE DB Provider for Terraform Authors
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/frederic-arr/ripedb-go/ripedb"
	"github.com/frederic-arr/ripedb-go/ripedb/models"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ObjectDataSource{}

func NewObjectDataSource() datasource.DataSource {
	return &ObjectDataSource{}
}

type ObjectDataSource struct {
	client *ripedb.RipeClient
}

func (d *ObjectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_object"
}

func (d *ObjectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This data source provides information a generic object in the RIPE Database.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "the ID of the object",
				Computed:            true,
			},
			"class": schema.StringAttribute{
				MarkdownDescription: "the class of the object",
				Required:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "the key of the object",
				Required:            true,
			},
			"attributes": schema.ListNestedAttribute{
				MarkdownDescription: "the attributes of the object",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "the name of the attribute",
							Computed:            true,
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "the value of the attribute",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *ObjectDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*RipeDbProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *RipeDbProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = data.Client
}

func (d *ObjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ObjectModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resource := data.Class.ValueString()
	key := data.Value.ValueString()
	obj, err := (*d.client).GetObject(resource, key)
	if err != nil {
		resp.Diagnostics.AddError("failed to query RIPE database", err.Error())
		return
	}

	m := models.ObjectToModelUnchecked(resource, *obj)
	objectToModel(obj, &data)
	data.Id = types.StringValue(fmt.Sprintf("%s:%s", resource, m.Key()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

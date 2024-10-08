// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ObjectDataSource{}

func NewObjectDataSource() datasource.DataSource {
	return &ObjectDataSource{}
}

type ObjectDataSource struct{}

func (d *ObjectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_object"
}

func (d *ObjectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This data source provides information about an Autonomous System Number (ASN) object in the RIPE Database.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "the ID of the object",
				Computed:            true,
			},
			"class": schema.StringAttribute{
				MarkdownDescription: "the class of the object",
				Required:            true,
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "the key of the object",
				Required:            true,
			},
			"attributes": schema.MapAttribute{
				MarkdownDescription: "a map of attributes with their values",
				Computed:            true,
				ElementType: types.ListType{
					ElemType: types.StringType,
				},
			},
			"raw_attributes": schema.ListNestedAttribute{
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
}

func (d *ObjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ObjectModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := queryData(data.Class.ValueString(), data.Key.ValueString(), &data)
	if err != nil {
		resp.Diagnostics.AddError("failed to query RIPE Database", err.Error())
		return
	}

	data.Id = types.StringValue(fmt.Sprintf("%s:%s", data.Class.ValueString(), data.Key.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

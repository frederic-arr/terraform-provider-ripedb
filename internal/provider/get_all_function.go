// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ function.Function = GetAllFunction{}
)

func NewGetAllFunction() function.Function {
	return GetAllFunction{}
}

type GetAllFunction struct{}

func (r GetAllFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "get_all"
}

func (r GetAllFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Get all values of an attribute",
		MarkdownDescription: "This function returns all the values of an attribute. \n" +
			"If the attribute does not exist, returns `[]` (*an empty list*).",
		Parameters: []function.Parameter{
			function.ListParameter{
				Name:                "attributes",
				MarkdownDescription: "The data to get the values from",
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name":  types.StringType,
						"value": types.StringType,
					},
				},
			},
			function.StringParameter{
				Name:                "key",
				MarkdownDescription: "The key of the attribute to get the values of",
			},
		},
		Return: function.ListReturn{
			ElementType: types.StringType,
		},
	}
}

func (r GetAllFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var attributes []ObjectModelAttribute
	var key types.String

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &attributes, &key))
	if resp.Error != nil {
		return
	}

	var values []types.String = make([]types.String, 0)
	for _, attr := range attributes {
		if attr.Name == key {
			if attr.Value.ValueString() == "" {
				continue
			}

			values = append(values, attr.Value)
		}
	}

	resp.Error = resp.Result.Set(ctx, values)
}

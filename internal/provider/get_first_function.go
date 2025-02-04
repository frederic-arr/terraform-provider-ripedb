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
	_ function.Function = GetFirstFunction{}
)

func NewGetFirstFunction() function.Function {
	return GetFirstFunction{}
}

type GetFirstFunction struct{}

func (r GetFirstFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "get_first"
}

func (r GetFirstFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Get the first value of an attribute",
		MarkdownDescription: "This function returns the first value of an attribute. \n" +
			"If the attribute does not exist, returns `null`.",
		Parameters: []function.Parameter{
			function.ListParameter{
				Name:                "attributes",
				MarkdownDescription: "The data to get the first value from",
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name":  types.StringType,
						"value": types.StringType,
					},
				},
			},
			function.StringParameter{
				Name:                "key",
				MarkdownDescription: "The key to get the first value of",
			},
		},
		Return: function.StringReturn{},
	}
}

func (r GetFirstFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var attributes []ObjectModelAttribute
	var key types.String

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &attributes, &key))
	if resp.Error != nil {
		return
	}

	for _, attr := range attributes {
		if attr.Name == key {
			resp.Error = resp.Result.Set(ctx, attr.Value)
			return
		}
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, types.StringNull()))
}

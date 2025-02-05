// Copyright (c) The RIPE DB Provider for Terraform Authors
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"slices"

	"github.com/frederic-arr/rpsl-go"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var OMIT_KEYS = []string{"source", "created", "last-modified"}

type ObjectModel struct {
	Id         types.String           `tfsdk:"id"`
	Class      types.String           `tfsdk:"class"`
	Value      types.String           `tfsdk:"value"`
	Attributes []ObjectModelAttribute `tfsdk:"attributes"`
}

type ObjectModelAttribute struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

func objectToModel(obj *rpsl.Object, data *ObjectModel) {
	data.Attributes = []ObjectModelAttribute{}
	for _, a := range obj.Attributes {
		data.Attributes = append(data.Attributes, ObjectModelAttribute{
			Name:  types.StringValue(a.Name),
			Value: types.StringValue(a.Value),
		})
	}
}

func modelToObject(data *ObjectModel) *rpsl.Object {
	obj := rpsl.Object{
		Attributes: make([]rpsl.Attribute, 0),
	}

	for _, attr := range data.Attributes {
		obj.Attributes = append(obj.Attributes, rpsl.Attribute{
			Name:  attr.Name.ValueString(),
			Value: attr.Value.ValueString(),
		})
	}

	return &obj
}

func filterObject(obj *rpsl.Object, data *ObjectModel) {
	data.Attributes = []ObjectModelAttribute{}
	for i, a := range obj.Attributes {
		if i == 0 || slices.Contains(OMIT_KEYS, a.Name) {
			continue
		}

		data.Attributes = append(data.Attributes, ObjectModelAttribute{
			Name:  types.StringValue(a.Name),
			Value: types.StringValue(a.Value),
		})
	}
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/frederic-arr/ripedb-go/ripedb"
	"github.com/frederic-arr/rpsl-go"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ObjectDataSourceModel describes the data source data model.
type ObjectModel struct {
	Id            types.String              `tfsdk:"id"`
	Class         types.String              `tfsdk:"class"`
	Key           types.String              `tfsdk:"key"`
	Attributes    map[string][]types.String `tfsdk:"attributes"`
	RawAttributes []ObjectModelAttribute    `tfsdk:"raw_attributes"`
}

type ObjectModelAttribute struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

func query(kind string, key string) (*rpsl.Object, error) {
	client := ripedb.NewRipeAnonymousClient()
	model, err := client.Get("ripe", kind, key)
	if err != nil {
		return nil, err
	}

	obj, err := model.FindOne()
	if err != nil {
		return nil, err
	}

	mapped := rpsl.Object{
		Attributes: []rpsl.Attribute{},
	}

	for _, attr := range obj.Attributes.Attribute {
		mapped.Attributes = append(mapped.Attributes, rpsl.Attribute{
			Name:  attr.Name,
			Value: attr.Value.(string),
		})
	}

	return &mapped, nil
}

func queryData(kind string, key string, data *ObjectModel) error {
	obj, err := query(kind, key)
	if err != nil {
		return err
	}

	data.Attributes = map[string][]types.String{}
	for _, attr := range obj.Attributes {
		if data.Attributes[attr.Name] == nil {
			data.Attributes[attr.Name] = []types.String{}
		}

		data.Attributes[attr.Name] = append(data.Attributes[attr.Name], types.StringValue(attr.Value))
		data.RawAttributes = append(data.RawAttributes, ObjectModelAttribute{
			Name:  types.StringValue(attr.Name),
			Value: types.StringValue(attr.Value),
		})
	}

	return nil
}

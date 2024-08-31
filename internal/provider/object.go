// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"net"

	"github.com/frederic-arr/rpsl-go"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ObjectDataSourceModel describes the data source data model.
type ObjectModel struct {
	Id         types.String           `tfsdk:"id"`
	Class      types.String           `tfsdk:"class"`
	Key        types.String           `tfsdk:"key"`
	Attributes []ObjectModelAttribute `tfsdk:"attributes"`
}

type ObjectModelAttribute struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

func query(kind string, key string) (*rpsl.Object, error) {
	server := "whois-test.ripe.net:43"
	query := fmt.Sprintf("-rBGT%s %s\r\n", kind, key)

	conn, err := net.Dial("tcp", server)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(query))
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 1024)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			break
		}
	}

	raw := string(buf)
	obj, err := rpsl.Parse(raw)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func queryData(kind string, key string, data *ObjectModel) error {
	obj, err := query(kind, key)
	if err != nil {
		return err
	}

	for _, attr := range obj.Attributes {
		data.Attributes = append(data.Attributes, ObjectModelAttribute{
			Name:  types.StringValue(attr.Name),
			Value: types.StringValue(attr.Value),
		})
	}

	return nil
}

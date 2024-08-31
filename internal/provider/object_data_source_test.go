// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccObjectDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "ripedb_object" "test" {
					class = "aut-num"
					key = "AS3333"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ripedb_object.test", "id", "aut-num:AS3333"),
					resource.TestCheckResourceAttr("data.ripedb_object.test", "attributes.0.name", "aut-num"),
					resource.TestCheckResourceAttr("data.ripedb_object.test", "attributes.0.value", "AS3333"),
				),
			},
		},
	})
}

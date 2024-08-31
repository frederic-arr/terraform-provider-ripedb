// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestGetAllFunction_Known(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				output "mnt_by" {
					value = provider::ripe::get_all([
						{ name = "mnt-by", value = "RIPE-MNT" },
						{ name = "mnt-by", value = "ARIN-MNT" },
					], "mnt-by")
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("mnt_by", knownvalue.ListExact(
						[]knownvalue.Check{
							knownvalue.StringExact("RIPE-MNT"),
							knownvalue.StringExact("ARIN-MNT"),
						},
					)),
				},
			},
		},
	})
}

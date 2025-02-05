// Copyright (c) The RIPE DB Provider for Terraform Authors
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
					value = provider::ripedb::get_all([
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

func TestGetAllFunction_Unknown(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				output "admin_c" {
					value = provider::ripedb::get_all([
						{ name = "mnt-by", value = "RIPE-MNT" },
						{ name = "mnt-by", value = "ARIN-MNT" },
					], "admin-c")
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("admin_c", knownvalue.ListExact(
						[]knownvalue.Check{},
					)),
				},
			},
		},
	})
}

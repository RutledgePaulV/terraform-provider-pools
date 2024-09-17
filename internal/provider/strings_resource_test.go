// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExampleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccExampleResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pools_strings.this", "resource_borrowers.one", "alice"),
					resource.TestCheckResourceAttr("pools_strings.this", "resource_borrowers.two", "bob"),
					resource.TestCheckNoResourceAttr("pools_strings.this", "resource_borrowers.three"),
					resource.TestCheckResourceAttr("pools_strings.this", "borrower_resources.alice", "one"),
					resource.TestCheckResourceAttr("pools_strings.this", "borrower_resources.bob", "two"),
				),
			},
			// Update and Read testing
			{
				Config: testAccExampleResourceConfigAgain(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pools_strings.this", "resource_borrowers.one", "alice"),
					resource.TestCheckResourceAttr("pools_strings.this", "resource_borrowers.three", "bob"),
					resource.TestCheckNoResourceAttr("pools_strings.this", "resource_borrowers.two"),
					resource.TestCheckResourceAttr("pools_strings.this", "borrower_resources.alice", "one"),
					resource.TestCheckResourceAttr("pools_strings.this", "borrower_resources.bob", "three"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccExampleResourceConfig() string {
	return `
resource "pools_strings" "this" {
  resources = ["one", "two", "three"]
  borrowers = ["alice", "bob"]
}
`
}

func testAccExampleResourceConfigAgain() string {
	return `
resource "pools_strings" "this" {
  resources = ["one", "three"]
  borrowers = ["alice", "bob"]
}
`
}

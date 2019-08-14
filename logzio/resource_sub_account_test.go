package logzio

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccLogzioSubAccount_CreateSubAccount(t *testing.T) {
	resourceName := "logzio_subaccount.create_subaccount"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: ReadFixtureFromFile("create_subaccount.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_name", "create_subaccount_name"),
				),
			},
		},
	})
}
package logzio

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccDataSourceLogShippingToken(t *testing.T) {
	resourceName := "data.logzio_log_shipping_token.my_log_shipping_token_datasource"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: true,
				Config:             ReadFixtureFromFile("create_log_shipping_token_datasource.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "my_token"),
				),
			},
		},
	})
}

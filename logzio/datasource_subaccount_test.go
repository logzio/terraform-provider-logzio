package logzio

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccDataSourceSubaccount(t *testing.T) {
	resourceName := "data.logzio_subaccount.subaccount_datasource_by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: true,
				Config: ReadFixtureFromFile("create_subaccount_datasource.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_name", "test"),
					resource.TestCheckResourceAttr(resourceName, "retention_days", "2"),
					resource.TestCheckResourceAttr(resourceName, "max_daily_gb", "1"),
				),
			},
		},
	})
}

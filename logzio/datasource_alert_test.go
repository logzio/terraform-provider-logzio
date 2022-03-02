package logzio

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"testing"
)

func TestAccDataSourceLogzIoAlert(t *testing.T) {
	resourceName := "data.logzio_alert.alert_datasource_by_id"
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: true,
				Config:             utils.ReadFixtureFromFile("create_alert_datasource.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "hello"),
					resource.TestCheckResourceAttr(resourceName, "query_string", "loglevel:ERROR"),
					resource.TestCheckResourceAttr(resourceName, "operation", "GREATER_THAN"),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "false"),
				),
			},
		},
	})
}

package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccDataSourceLogzIoAlert(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: true,
				Config:             testAccDataSourceLogzioAlertConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.logzio_alert.by_title", "title", "hello"),
					resource.TestCheckResourceAttr("data.logzio_alert.by_title", "query_string", "loglevel:ERROR"),
					resource.TestCheckResourceAttr("data.logzio_alert.by_title", "operation", "GREATER_THAN"),
				),
			},
		},
	})
}

func testAccDataSourceLogzioAlertBase() string {
	return fmt.Sprintf(`
resource "logzio_alert" "by_title" {
  title = "hello"
  query_string = "loglevel:ERROR"
  operation = "GREATER_THAN"
  notification_emails = ["testx@test.com"]
  search_timeframe_minutes = 5
  value_aggregation_type = "NONE"
  alert_notification_endpoints = []
  suppress_notifications_minutes = 5
  severity_threshold_tiers = [
    {
      "severity" = "HIGH",
      "threshold" = 10
    }
  ]
}
`)
}

func testAccDataSourceLogzioAlertConfig() string {
	return testAccDataSourceLogzioAlertBase() + `

data "logzio_alert" "by_title" {
  title = "hello"
  depends_on = ["logzio_alert.by_title"]
}

`
}

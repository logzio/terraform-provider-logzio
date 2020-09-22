package logzio

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccDataSourceLogzIoAlert(t *testing.T) {
	resourceName := "data.logzio_alert.alert_datasource_by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: true,
				Config: ReadFixtureFromFile("create_alert_datasource.tf"),
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

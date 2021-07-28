package logzio

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccDataSourceLogzIoAlertV2(t *testing.T) {
	resourceName := "data.logzio_alert_v2.alert_v2_datasource_by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: true,
				Config:             ReadFixtureFromFile("create_alert_v2_datasource.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "hello"),
					resource.TestCheckResourceAttr(resourceName, "sub_components.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sub_components.0.severity_threshold_tiers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sub_components.0.severity_threshold_tiers.0.severity", "HIGH"),
					resource.TestCheckResourceAttr(resourceName, "sub_components.0.severity_threshold_tiers.0.threshold", "10"),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "sub_components.0.filter_must"),
				),
			},
		},
	})
}

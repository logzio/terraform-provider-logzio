package logzio

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"testing"
)

func TestAccDataSourceLogzIoAlertV2(t *testing.T) {
	resourceName := "data.logzio_alert_v2.alert_v2_datasource_by_id"
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckApiToken(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:  getConfigResourceAlertV2("10"),
				Destroy: false,
			},
			{
				Config: getConfigResourceAlertV2("10") + getConfigDatasourceAlertV2(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "hello"),
					resource.TestCheckResourceAttr(resourceName, "sub_components.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sub_components.0.severity_threshold_tiers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sub_components.0.severity_threshold_tiers.0.severity", "HIGH"),
					resource.TestCheckResourceAttr(resourceName, "sub_components.0.severity_threshold_tiers.0.threshold", "10"),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "sub_components.0.filter_must"),
					resource.TestCheckResourceAttr(resourceName, alertV2ScheduleCronExpression, "0 0/5 9-17 ? * * *"),
					resource.TestCheckResourceAttr(resourceName, alertV2ScheduleTimezone, "IET"),
				),
			},
			{
				Config: getConfigResourceAlertV2("10.5") + getConfigDatasourceAlertV2(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "hello"),
					resource.TestCheckResourceAttr(resourceName, "sub_components.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sub_components.0.severity_threshold_tiers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sub_components.0.severity_threshold_tiers.0.severity", "HIGH"),
					resource.TestCheckResourceAttr(resourceName, "sub_components.0.severity_threshold_tiers.0.threshold", "10.5"),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "sub_components.0.filter_must"),
					resource.TestCheckResourceAttr(resourceName, alertV2ScheduleCronExpression, "0 0/5 9-17 ? * * *"),
					resource.TestCheckResourceAttr(resourceName, alertV2ScheduleTimezone, "IET"),
				),
			},
		},
	})
}

func getConfigResourceAlertV2(threshold string) string {
	return `resource "logzio_alert_v2" "alert_v2_datasource" {
  title = "hello"
  description = "this is a description"
  tags = [
    "some",
    "test"]
  search_timeframe_minutes = 5
  is_enabled = false
  notification_emails = [
    "testx@test.com"]
  suppress_notifications_minutes = 5
  output_type = "JSON"
  sub_components {
    query_string = "loglevel:ERROR"
    should_query_on_all_accounts = true
    operation = "EQUALS"
    value_aggregation_type = "COUNT"
    severity_threshold_tiers {
      severity = "HIGH"
      threshold =` + threshold + `
    }
    filter_must = jsonencode([
      {
        match_phrase: {
          "some_field": {
            "query": "some_query"
          }
        }
      },
      {
        another_match: {
          "some_field2": {
            "query": "hello world"
          }
        }
      }
    ])
  }
  schedule_cron_expression = "0 0/5 9-17 ? * * *"
  schedule_timezone = "IET"
}
`
}

func getConfigDatasourceAlertV2() string {
	return `data "logzio_alert_v2" "alert_v2_datasource_by_id" {
  id = "${logzio_alert_v2.alert_v2_datasource.id}"
  depends_on = ["logzio_alert_v2.alert_v2_datasource"]
}
`
}

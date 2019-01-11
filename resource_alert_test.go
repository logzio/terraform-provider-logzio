package main

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jonboydell/logzio_client/alerts"
	"os"
	"strconv"
	"testing"
)

func TestAccLogzioAlert_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccLogzioAlertDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLogzioAlertConfig("name"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogzioAlertExists("logzio_alert.name"),
					resource.TestCheckResourceAttr(
						"logzio_alert.name", "title", "my_other_title"),
				),
			},
		},
	})
}

func TestAccLogzioAlert_SomeTest(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccLogzioAlertDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLogzioAlertConfig("name"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogzioAlertExists("logzio_alert.name"),
					resource.TestCheckResourceAttr(
						"logzio_alert.name", "title", "my_other_title"),
				),
			},
			resource.TestStep{
				Config: testAccUpdateLogzioAlertConfig("name"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLogzioAlertExists("logzio_alert.name"),
					resource.TestCheckResourceAttr(
						"logzio_alert.name", "title", "this_is_my_title"),
				),
			},
		},
	})
}

func testAccCheckLogzioAlertExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No alert ID is set")
		}

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)

		var client *alerts.Alerts
		client, _ = alerts.New(os.Getenv("LOGZIO_API_TOKEN"))

		_, err = client.GetAlert(int64(id))

		if err != nil {
			return fmt.Errorf("Alert doesn't exist")
		}

		return nil
	}
}

func testAccLogzioAlertDestroy(s *terraform.State) error {

	for _, r := range s.RootModule().Resources {
		id, err := strconv.ParseInt(r.Primary.ID, 10, 64)
		if err != nil {
			return err
		}

		var client *alerts.Alerts
		client, _ = alerts.New(os.Getenv("LOGZIO_API_TOKEN"))

		_, err = client.GetAlert(int64(id))

		if err == nil {
			return fmt.Errorf("Alert still exists")
		}
	}
	return nil
}

func TestAccDataSourceLogzIoAlert(t *testing.T) {
	rName := "some_name"
	rTitle := "some_title"
	resourceName := "logzio_alert.some_name"
	datasourceName := "data.logzio_alert.by_title"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLogzioAlertConfig(rName, rTitle),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceLogzIoAlertCheck(datasourceName, resourceName),
				),
			},
		},
	})
}

func testAccDataSourceLogzIoAlertCheck(datasourceName, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[datasourceName]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", datasourceName)
		}

		alertRs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", resourceName)
		}

		attrNames := []string{
			"title",
			"query_string",
			"operation",
			"notification_emails",
			"search_timeframe_minutes",
			"value_aggregation_type",
			"alert_notification_endpoonts",
			"suppress_notification_minutes",
			"severity_threshold_tiers",
		}

		for _, attrName := range attrNames {
			if ds.Primary.Attributes[attrName] != alertRs.Primary.Attributes[attrName] {
				return fmt.Errorf(
					"%s is %s; want %s",
					attrName,
					ds.Primary.Attributes[attrName],
					alertRs.Primary.Attributes[attrName],
				)
			}
		}

		return nil
	}
}

func testAccCheckLogzioAlertConfig(rName string) string {
	return fmt.Sprintf(`
resource "logzio_alert" "%s" {
  title = "my_other_title"
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
`, rName)
}

func testAccUpdateLogzioAlertConfig(rName string) string {
	return fmt.Sprintf(`
resource "logzio_alert" "%s" {
  title = "this_is_my_title"
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
`, rName)
}

func testAccDataSourceLogzioAlertConfig(rName string, rTitle string) string {
	return fmt.Sprintf(`
resource "logzio_alert" "%s" {
  title = "%s"
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

data "logzio_alert" "by_title" {
  title = "%s"
  depends_on = ["logzio_alert.%s"]
}
`, rName, rTitle, rTitle, rName)
}
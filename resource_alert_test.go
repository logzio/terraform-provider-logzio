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

		var client *alerts.AlertsClient
		client, _ = alerts.New(os.Getenv(envLogzioApiToken))

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

		var client *alerts.AlertsClient
		client, _ = alerts.New(os.Getenv(envLogzioApiToken))

		_, err = client.GetAlert(int64(id))

		if err == nil {
			return fmt.Errorf("alert still exists")
		}
	}
	return nil
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

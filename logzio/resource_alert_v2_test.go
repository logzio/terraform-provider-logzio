package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/logzio/logzio_terraform_client/alerts_v2"
	"io/ioutil"
	"log"
	"testing"
)

const (
	alertsV2ResourceCreateAlert string = "create_alert_v2"
	alertsV2ResourceUpdateAlert string = "update_alert_v2"
)

func TestAccLogzioAlertV2_CreateAlert(t *testing.T) {
	alertName := "test_create_alert_v2"
	resourceName := "logzio_alert_v2." + alertName

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceTestAlertV2(alertName, alertsV2ResourceCreateAlert),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "hello"),
					resource.TestCheckResourceAttr(resourceName, "sub_components.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sub_components.0.severity_threshold_tiers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sub_components.0.severity_threshold_tiers.0.severity", "HIGH"),
					resource.TestCheckResourceAttr(resourceName, "sub_components.0.severity_threshold_tiers.0.threshold", "10"),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "false"),
				),
			},
			{
				Config: resourceTestAlertV2(alertName, alertsV2ResourceCreateAlert),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogzioAlertV2_UpdateAlert(t *testing.T) {
	alertName := "test_update_alert_v2"
	resourceName := "logzio_alert_v2." + alertName

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceTestAlertV2(alertName, alertsV2ResourceCreateAlert),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "hello"),
				),
			},
			{
				Config: resourceTestAlertV2("test_update_alert_v2", alertsV2ResourceUpdateAlert),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "updated_alert"),
					resource.TestCheckResourceAttr(resourceName, "sub_components.0.severity_threshold_tiers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sub_components.0.value_aggregation_type", alerts_v2.AggregationTypeSum),
					resource.TestCheckResourceAttr(resourceName, "sub_components.0.value_aggregation_field", "some_field"),
				),
			},
		},
	})
}

func resourceTestAlertV2(name string, path string) string {
	content, err := ioutil.ReadFile(fmt.Sprintf("testdata/fixtures/%s.tf", path))
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf(fmt.Sprintf("%s", content), name)
}

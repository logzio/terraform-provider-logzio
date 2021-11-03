package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
					resource.TestCheckResourceAttr(resourceName, "sub_components.0.severity_threshold_tiers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "sub_components.0.severity_threshold_tiers.2401098826.severity", "LOW"),
					resource.TestCheckResourceAttr(resourceName, "sub_components.0.severity_threshold_tiers.869563786.severity", "HIGH"),
					resource.TestCheckResourceAttr(resourceName, "sub_components.0.severity_threshold_tiers.2401098826.threshold", "2"),
					resource.TestCheckResourceAttr(resourceName, "sub_components.0.severity_threshold_tiers.869563786.threshold", "10"),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "sub_components.0.filter_must"),
				),
			},
			{
				Config:            resourceTestAlertV2(alertName, alertsV2ResourceCreateAlert),
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
				Config: resourceTestAlertV2(alertName, alertsV2ResourceUpdateAlert),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "updated_alert"),
					resource.TestCheckResourceAttr(resourceName, "sub_components.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sub_components.0.severity_threshold_tiers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "sub_components.0.severity_threshold_tiers.2401098826.severity", "LOW"),
					resource.TestCheckResourceAttr(resourceName, "sub_components.0.severity_threshold_tiers.869563786.severity", "HIGH"),
					resource.TestCheckResourceAttr(resourceName, "sub_components.0.severity_threshold_tiers.2401098826.threshold", "2"),
					resource.TestCheckResourceAttr(resourceName, "sub_components.0.severity_threshold_tiers.869563786.threshold", "10"),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "sub_components.0.filter_must"),
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

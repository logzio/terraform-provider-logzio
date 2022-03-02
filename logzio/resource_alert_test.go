package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/logzio/logzio_terraform_provider/logzio/utils"
	"io/ioutil"
	"log"
	"testing"
)

func TestAccLogzioAlert_CreateAlert(t *testing.T) {
	alertName := "test_create_alert"
	resourceName := "logzio_alert." + alertName
	defer utils.SleepAfterTest()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceCreateAlert(alertName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "hello"),
					resource.TestCheckResourceAttr(resourceName, "severity_threshold_tiers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "severity_threshold_tiers.0.severity", "HIGH"),
					resource.TestCheckResourceAttr(resourceName, "severity_threshold_tiers.0.threshold", "10"),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "false"),
				),
			},
			{
				Config:            resourceCreateAlert(alertName),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogzioAlert_UpdateAlert(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckApiToken(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceCreateAlert("test_update_alert"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"logzio_alert.test_update_alert", "title", "hello"),
				),
			},
			{
				Config: resourceUpdateAlert("test_update_alert"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"logzio_alert.test_update_alert", "title", "updated_alert"),
				),
			},
		},
	})
}

func resourceCreateAlert(name string) string {
	content, err := ioutil.ReadFile("testdata/fixtures/create_alert.tf")
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf(fmt.Sprintf("%s", content), name)
}

func resourceUpdateAlert(name string) string {
	content, err := ioutil.ReadFile("testdata/fixtures/update_alert.tf")
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf(fmt.Sprintf("%s", content), name)
}

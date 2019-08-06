package logzio

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"io/ioutil"
	"log"
	"testing"
)

func TestAccLogzioAlert_CreateAlert(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceLogzioAlertBase("name"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("logzio_alert.name", "title", "hello"),
					resource.TestCheckResourceAttr("logzio_alert.name", "severity_threshold_tiers.#", "1"),
					resource.TestCheckResourceAttr("logzio_alert.name", "severity_threshold_tiers.0.severity", "HIGH"),
					resource.TestCheckResourceAttr("logzio_alert.name", "severity_threshold_tiers.0.threshold", "10"),
				),
			},
		},
	})
}

func TestAccLogzioAlert_UpdateAlert(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: resourceCreateAlert("name"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"logzio_alert.name", "title", "hello"),
				),
			},
			resource.TestStep{
				Config: resourceUpdateAlert("name"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"logzio_alert.name", "title", "updated_alert"),
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
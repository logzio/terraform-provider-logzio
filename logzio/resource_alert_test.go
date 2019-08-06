package logzio

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jonboydell/logzio_client/alerts"
)

func TestAccLogzioAlert_CreateAlert(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkAlertDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceLogzioAlertBase("name"),
				Check: resource.ComposeTestCheckFunc(
					checkAlertExists("logzio_alert.name"),
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
		CheckDestroy: checkAlertDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: resourceCreateAlert("name"),
				Check: resource.ComposeTestCheckFunc(
					checkAlertExists("logzio_alert.name"),
					resource.TestCheckResourceAttr(
						"logzio_alert.name", "title", "hello"),
				),
			},
			resource.TestStep{
				Config: resourceUpdateAlert("name"),
				Check: resource.ComposeTestCheckFunc(
					checkAlertExists("logzio_alert.name"),
					resource.TestCheckResourceAttr(
						"logzio_alert.name", "title", "updated_alert"),
				),
			},
		},
	})
}

func checkAlertExists(n string) resource.TestCheckFunc {
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
		baseURL := defaultBaseUrl
		if len(os.Getenv(envLogzioBaseURL)) > 0 {
			baseURL = os.Getenv(envLogzioBaseURL)
		}
		client, _ = alerts.New(os.Getenv(envLogzioApiToken), baseURL)

		time.Sleep(2 * time.Second)
		_, err = client.GetAlert(int64(id))

		if err != nil {
			return fmt.Errorf("Alert doesn't exist")
		}

		return nil
	}
}

func checkAlertDestroyed(s *terraform.State) error {

	for _, r := range s.RootModule().Resources {
		id, err := strconv.ParseInt(r.Primary.ID, 10, 64)
		if err != nil {
			return err
		}

		var client *alerts.AlertsClient
		baseURL := defaultBaseUrl
		if len(os.Getenv(envLogzioBaseURL)) > 0 {
			baseURL = os.Getenv(envLogzioBaseURL)
		}
		client, _ = alerts.New(os.Getenv(envLogzioApiToken), baseURL)

		time.Sleep(2 * time.Second)
		_, err = client.GetAlert(int64(id))
		if err == nil {
			return fmt.Errorf("alert still exists")
		}
	}
	return nil
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
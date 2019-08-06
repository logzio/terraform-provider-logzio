package logzio

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jonboydell/logzio_client/endpoints"
)

func TestAccLogzioEndpoint_CreateSlackEndpoint(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkEndpointDestroyed,
		Steps: []resource.TestStep{
			{
				Config: readFixtureFromFile("valid_slack_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
					checkEndpointExists("logzio_endpoint.valid_slack_endpoint"),
					resource.TestCheckResourceAttr(
						"logzio_endpoint.valid_slack_endpoint", "title", "valid_slack_endpoint"),
					testAccCheckOutputExists("logzio_endpoint.valid_slack_endpoint", "test_id"),
					resource.TestMatchOutput("test_id", regexp.MustCompile("\\d")),
				),
			},
		},
	})
}

func TestAccLogzioEndpoint_CreateInvalidSlackEndpoint(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkEndpointDestroyed,
		Steps: []resource.TestStep{
			{
				Config:      readFixtureFromFile("invalid_slack_endpoint.tf"),
				ExpectError: regexp.MustCompile("Bad URL provided. no protocol"),
			},
		},
	})
}

func readFixtureFromFile(name string) string {
	content, err := ioutil.ReadFile("testdata/fixtures/"+name)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s", content)
}


func TestAccLogzioEndpoint_UpdateSlackEndpoint(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkEndpointDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLogzioEndpointConfig("slackHappyPath"),
				Check: resource.ComposeTestCheckFunc(
					checkEndpointExists("logzio_endpoint.slack"),
					resource.TestCheckResourceAttr(
						"logzio_endpoint.slack", "title", "my_slack_title"),
				),
			},
			{
				Config: testAccCheckLogzioEndpointConfig("slackUpdateHappyPath"),
				Check: resource.ComposeTestCheckFunc(
					checkEndpointExists("logzio_endpoint.slack"),
					resource.TestCheckResourceAttr(
						"logzio_endpoint.slack", "title", "my_updated_slack_title"),
				),
			},
		},
	})
}

func TestAccLogzioEndpoint_CreateCustomEndpoint(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkEndpointDestroyed,
		Steps: []resource.TestStep{
			{
				Config: createCustomEndpoint("custom"),
				Check: resource.ComposeTestCheckFunc(
					checkEndpointExists("logzio_endpoint.custom"),
					resource.TestCheckResourceAttr(
						"logzio_endpoint.custom", "title", "my_custom_title"),
				),
			},
		},
	})
}

func TestAccLogzioEndpoint_PagerDuty_HappyPath(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkEndpointDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLogzioEndpointConfig("pagerDutyHappyPath"),
				Check: resource.ComposeTestCheckFunc(
					checkEndpointExists("logzio_endpoint.pagerduty"),
					resource.TestCheckResourceAttr(
						"logzio_endpoint.pagerduty", "title", "my_pagerduty_title"),
				),
			},
		},
	})
}

func TestAccLogzioEndpoint_BigPanda_HappyPath(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkEndpointDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLogzioEndpointConfig("bigPandaHappyPath"),
				Check: resource.ComposeTestCheckFunc(
					checkEndpointExists("logzio_endpoint.bigpanda"),
					resource.TestCheckResourceAttr(
						"logzio_endpoint.bigpanda", "title", "my_bigpanda_title"),
					resource.TestCheckResourceAttr(
						"logzio_endpoint.bigpanda", "big_panda.#", "1"),
					resource.TestCheckResourceAttr(
						"logzio_endpoint.bigpanda", "big_panda.1922960384.api_token", "my_api_token"),
					resource.TestCheckResourceAttr(
						"logzio_endpoint.bigpanda", "big_panda.1922960384.app_key", "my_app_key"),
				),
			},
		},
	})
}

func testAccCheckOutputExists(n string, o string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		id := rs.Primary.ID
		os, ok := s.RootModule().Outputs[o]

		if rs.Primary.ID == "" {
			return errors.New("no endpoint ID is set")
		}

		if os.Value != id {
			return fmt.Errorf("can't find resource that matches output ID")
		}

		return nil
	}
}

func checkEndpointExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no endpoint ID is set")
		}

		id, err := strconv.ParseInt(rs.Primary.ID, BASE_10, BITSIZE_64)

		var client *endpoints.EndpointsClient
		baseURL := defaultBaseUrl
		if len(os.Getenv(envLogzioBaseURL)) > 0 {
			baseURL = os.Getenv(envLogzioBaseURL)
		}
		client, _ = endpoints.New(os.Getenv(envLogzioApiToken), baseURL)

		_, err = client.GetEndpoint(int64(id))

		if err != nil {
			return fmt.Errorf("endpoint doesn't exist")
		}

		return nil
	}
}

func checkEndpointDestroyed(s *terraform.State) error {
	for _, r := range s.RootModule().Resources {
		id, err := strconv.ParseInt(r.Primary.ID, BASE_10, BITSIZE_64)
		if err != nil {
			return err
		}

		var client *endpoints.EndpointsClient
		baseURL := defaultBaseUrl
		if len(os.Getenv(envLogzioBaseURL)) > 0 {
			baseURL = os.Getenv(envLogzioBaseURL)
		}
		client, _ = endpoints.New(os.Getenv(envLogzioApiToken), baseURL)

		time.Sleep(3 * time.Second)
		_, err = client.GetEndpoint(int64(id))
		if err == nil {
			return fmt.Errorf("endpoint still exists")
		}
	}
	return nil
}

func testAccCheckLogzioEndpointConfig(key string) string {
	templates := map[string]string{
		"slackHappyPath": `
resource "logzio_endpoint" "slack" {
  title = "my_slack_title"
  endpoint_type = "Slack"
  description = "this_is_my_description"
  slack {
	url = "https://www.test.com"
  }
}

output "test_id" {
	value = "${logzio_endpoint.slack.id}"
}
`,
		"slackBadUrl": `
resource "logzio_endpoint" "slack" {
  title = "my_slack_title"
  endpoint_type = "Slack"
  description = "this_is_my_description"
  slack {
	url = "https://not_a_url"
  }
}
`,
		"slackUpdateHappyPath": `
resource "logzio_endpoint" "slack" {
  title = "my_updated_slack_title"
  endpoint_type = "Slack"
  description = "this_is_my_description"
  slack {
	url = "https://www.test.com"
  }
}
`,
		"customHappyPath": `
resource "logzio_endpoint" "custom" {
  title = "my_custom_title"
  endpoint_type = "Custom"
  description = "this_is_my_description"
  custom {
	url = "https://www.test.com"
	method = "POST"
	headers = {
		this = "is"
		a = "header"
	}
	body_template = {
		this = "is"
		my = "template"
	}
  }
}
`,
		"pagerDutyHappyPath": `
	resource "logzio_endpoint" "pagerduty" {
		title = "my_pagerduty_title"
		endpoint_type = "PagerDuty"
		description = "this is my description"
		pager_duty {
			service_key = "my_service_key"
		}
	}
`,
		"bigPandaHappyPath": `
	resource "logzio_endpoint" "bigpanda" {
		title = "my_bigpanda_title"
		endpoint_type = "BigPanda"
		description = "this is my description"
		big_panda {
			api_token = "my_api_token"
			app_key = "my_app_key"
		}
	}
`,
	}
	return templates[key]
}


func createCustomEndpoint(name string) string {
	content, err := ioutil.ReadFile("testdata/fixtures/create_custom_endpoint.tf")
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf(fmt.Sprintf("%s", content), name)
}

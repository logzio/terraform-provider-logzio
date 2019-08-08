package logzio

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"io/ioutil"
	"log"
	"regexp"
	"testing"
)

func TestAccLogzioEndpoint_CreateSlackEndpoint(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ReadFixtureFromFile("valid_slack_endpoint.tf"),
				Check: resource.ComposeTestCheckFunc(
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
		Steps: []resource.TestStep{
			{
				Config:      ReadFixtureFromFile("invalid_slack_endpoint.tf"),
				ExpectError: regexp.MustCompile("Bad URL provided. no protocol"),
			},
		},
	})
}

func TestAccLogzioEndpoint_UpdateSlackEndpoint(t *testing.T) {
	endpointName := "test_create_slack_endpoint"
	resourceName := "logzio_endpoint." + endpointName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ReadResourceFromFile("create_slack_endpoint.tf", endpointName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "title", "slack_endpoint"),
				),
			},
			{
				Config: ReadResourceFromFile("update_slack_endpoint.tf", endpointName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName, "title", "updated_slack_endpoint"),
				),
			},
		},
	})
}

func TestAccLogzioEndpoint_CreateCustomEndpoint(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: createCustomEndpoint("custom"),
				Check: resource.ComposeTestCheckFunc(
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
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLogzioEndpointConfig("pagerDutyHappyPath"),
				Check: resource.ComposeTestCheckFunc(
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
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLogzioEndpointConfig("bigPandaHappyPath"),
				Check: resource.ComposeTestCheckFunc(
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
